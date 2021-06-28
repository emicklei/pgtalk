package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/emicklei/tre"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

const pgLoadTableDef = `
SELECT
c.relkind AS type,
c.relname AS table_name
FROM pg_class c
JOIN ONLY pg_namespace n ON n.oid = c.relnamespace
WHERE n.nspname = $1
AND c.relkind = 'r'
ORDER BY c.relname
`

// PgTable postgres table
type PgTable struct {
	Schema               string
	Name                 string
	DataType             int32
	HasAutogeneratingKey bool
	PrimaryKeys          []PgColumn
	Columns              []PgColumn
}

// PgClass postgres pg_class
type PgClass struct {
	Name     string
	DataType int32
	Columns  []PgColumn
}

// SQL queries taken from https://github.com/achiku/dgw/blob/master/dgw.go

const pgLoadColumnDef = `
SELECT
    a.attnum AS field_ordinal,
    a.attname AS column_name,
    format_type(a.atttypid, a.atttypmod) AS data_type,
    a.attnotnull AS not_null,
    COALESCE(pg_get_expr(ad.adbin, ad.adrelid), '') AS default_value,
    COALESCE(ct.contype = 'p', false) AS  is_primary_key,
    CASE
        WHEN a.atttypid = ANY ('{int,int8,int2}'::regtype[])
          AND EXISTS (
             SELECT 1 FROM pg_attrdef ad
             WHERE  ad.adrelid = a.attrelid
             AND    ad.adnum   = a.attnum
             AND    pg_get_expr(ad.adbin, ad.adrelid) = 'nextval('''
                || (pg_get_serial_sequence (a.attrelid::regclass::text
                                          , a.attname))::regclass
                || '''::regclass)'
             )
            THEN CASE a.atttypid
                    WHEN 'int'::regtype  THEN 'serial'
                    WHEN 'int8'::regtype THEN 'bigserial'
                    WHEN 'int2'::regtype THEN 'smallserial'
                 END
        WHEN a.atttypid = ANY ('{uuid}'::regtype[]) AND COALESCE(pg_get_expr(ad.adbin, ad.adrelid), '') != ''
            THEN 'autogenuuid'
        ELSE format_type(a.atttypid, a.atttypmod)
    END AS data_type
FROM pg_attribute a
JOIN ONLY pg_class c ON c.oid = a.attrelid
JOIN ONLY pg_namespace n ON n.oid = c.relnamespace
LEFT JOIN pg_constraint ct ON ct.conrelid = c.oid
AND a.attnum = ANY(ct.conkey) AND ct.contype = 'p'
LEFT JOIN pg_attrdef ad ON ad.adrelid = c.oid AND ad.adnum = a.attnum
WHERE a.attisdropped = false
AND n.nspname = $1
AND c.relname = $2
AND a.attnum > 0
ORDER BY a.attnum
`

// PgColumn postgres columns
type PgColumn struct {
	FieldOrdinal int
	Name         string
	DataType     string
	DDLType      string
	NotNull      bool
	DefaultValue sql.NullString
	IsPrimaryKey bool
}

func LoadTables(ctx context.Context, conn *pgx.Conn, schema string) (list []PgTable, err error) {
	classes, err := loadTables(ctx, conn, schema)
	if err != nil {
		return list, tre.New(err, "LoadTables", "schema", schema)
	}
	for _, each := range classes {
		columns, err := loadTableColumns(ctx, conn, schema, each.Name)
		if err != nil {
			return list, tre.New(err, "LoadTables", "schema", schema)
		}
		primColums, isAutoGen := selectPrimaryKeys(columns)
		table := PgTable{
			Schema:               schema,
			Name:                 each.Name,
			DataType:             each.DataType,
			Columns:              columns,
			PrimaryKeys:          primColums,
			HasAutogeneratingKey: isAutoGen,
		}
		list = append(list, table)
	}
	return
}

func loadTables(ctx context.Context, conn *pgx.Conn, schema string) (list []PgClass, err error) {
	rows, err := conn.Query(ctx, pgLoadTableDef, schema)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		cls := PgClass{}
		if err := rows.Scan(&cls.DataType, &cls.Name); err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) {
				fmt.Println(pgErr.Message) // => syntax error at end of input
				fmt.Println(pgErr.Code)    // => 42601
			}
			return list, err
		}
		list = append(list, cls)
	}
	return
}

func loadTableColumns(ctx context.Context, conn *pgx.Conn, schema, table string) (list []PgColumn, err error) {
	rows, err := conn.Query(ctx, pgLoadColumnDef, schema, table)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		each := PgColumn{}
		if err := rows.Scan(
			&each.FieldOrdinal,
			&each.Name,
			&each.DataType,
			&each.NotNull,
			&each.DefaultValue,
			&each.IsPrimaryKey,
			&each.DDLType,
		); err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) {
				fmt.Println(pgErr.Message) // => syntax error at end of input
				fmt.Println(pgErr.Code)    // => 42601
			}
			return list, err
		}
		list = append(list, each)
	}
	return
}

var autoGenKeyTypes = []string{"smallserial", "serial", "bigserial", "autogenuuid"}

func selectPrimaryKeys(cols []PgColumn) (list []PgColumn, isAutogenerating bool) {
	isAutogenerating = false
	for _, each := range cols {
		if each.IsPrimaryKey {
			list = append(list, each)
			for _, other := range autoGenKeyTypes {
				if each.DDLType == other {
					isAutogenerating = true
				}
			}
		}
	}
	return
}
