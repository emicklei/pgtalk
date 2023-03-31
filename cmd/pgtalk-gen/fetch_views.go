package main

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/emicklei/tre"
	pgx "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const pgLoadViewDef = `
select viewname from pg_catalog.pg_views where schemaname = $1;`

func LoadViews(ctx context.Context, conn *pgx.Conn, schema string) (list []PgTable, err error) {
	views, err := loadViews(ctx, conn, schema)
	if err != nil {
		return list, tre.New(err, "LoadViews", "schema", schema)
	}
	if *oVerbose {
		log.Println("found views:", len(views))
	}
	for _, each := range views {
		columns, err := loadViewColumns(ctx, conn, schema, each.Name)
		if err != nil {
			return list, tre.New(err, "LoadViews", "schema", schema)
		}
		if *oVerbose {
			log.Println("found columns in view:", each.Name, len(columns))
		}
		table := PgTable{
			Schema:   schema,
			Name:     each.Name,
			DataType: each.DataType,
			Columns:  columns,
		}
		list = append(list, table)
	}
	return
}

func loadViews(ctx context.Context, conn *pgx.Conn, schema string) (list []PgClass, err error) {
	rows, err := conn.Query(ctx, pgLoadViewDef, schema)
	if err != nil {
		return list, tre.New(err, "loadViews", "schema", schema)
	}
	defer rows.Close()

	for rows.Next() {
		cls := PgClass{}
		if err := rows.Scan(&cls.Name); err != nil {
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

func loadViewColumns(ctx context.Context, conn *pgx.Conn, schema, view string) (list []PgColumn, err error) {
	rows, err := conn.Query(ctx, pgLoadViewColumnDef, view)
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

const pgLoadViewColumnDef = `
SELECT
	a.attnum AS field_ordinal,
	a.attname AS column_name,
	format_type(a.atttypid, a.atttypmod) AS data_type,
	a.attnotnull AS not_null
FROM pg_class c
INNER JOIN pg_attribute a ON a.attrelid = c.oid
INNER JOIN pg_type t ON t.oid = a.atttypid
INNER JOIN pg_namespace n ON n.oid = c.relnamespace
WHERE c.relkind = 'v'
    AND c.relname = $1;
`
