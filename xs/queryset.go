package xs

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/jackc/pgx/v4"
)

type QuerySet struct {
	tableName  string
	tableAlias string
	selectors  []ReadWrite
	condition  SQLWriter
	limit      int
}

func MakeQuerySet(tableName string, selectors []ReadWrite) QuerySet {
	// TODO get next free alias
	return QuerySet{tableName: tableName, tableAlias: "t1", selectors: selectors, condition: EmptyCondition}
}

// String returns the full SQL query
func (q QuerySet) SQL() string {
	// TODO use GORM here
	buf := new(bytes.Buffer)
	for i, each := range q.selectors {
		if i > 0 {
			io.WriteString(buf, ",")
		}
		if q.tableAlias != "" {
			io.WriteString(buf, q.tableAlias)
			io.WriteString(buf, ".")
		}
		io.WriteString(buf, each.Name())
	}
	// TEMP
	ctx := SQLContext{}
	where := q.condition.SQL(ctx)
	if len(where) > 0 {
		where = fmt.Sprintf(" WHERE %s", where)
	}
	limit := ""
	if q.limit > 0 {
		limit = fmt.Sprintf(" LIMIT %d", q.limit)
	}
	table := q.tableName
	if q.tableAlias != "" {
		table = fmt.Sprintf("%s %s", table, q.tableAlias)
	}
	return fmt.Sprintf("SELECT %s FROM %s%s%s", buf, table, where, limit)
}

func (q QuerySet) Where(condition SQLWriter) QuerySet {
	return QuerySet{tableName: q.tableName, tableAlias: q.tableAlias, selectors: q.selectors, condition: condition}
}

func (q QuerySet) Limit(limit int) QuerySet {
	return QuerySet{tableName: q.tableName, tableAlias: q.tableAlias, selectors: q.selectors, condition: q.condition, limit: limit}
}

func (d QuerySet) Exec(conn *pgx.Conn, factory func() interface{}, appender func(each interface{})) (err error) {
	rows, err := conn.Query(context.Background(), d.SQL())
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		entity := factory()
		sw := []interface{}{}
		for _, each := range d.selectors {
			rw := ScanToWrite{
				RW:     each,
				Entity: entity,
			}
			sw = append(sw, rw)
		}
		if err := rows.Scan(sw...); err != nil {
			return err
		}
		appender(entity)
	}
	return
}

type Unwrappable interface {
	Unwrap() QuerySet
}

func (d QuerySet) Join(otherQuerySet Unwrappable, onLeft, onRight ReadWrite) InnerJoin {
	return InnerJoin{
		LeftSet:  d,
		RightSet: otherQuerySet.Unwrap(),
		OnLeft:   onLeft,
		onRight:  onRight,
	}
}

type InnerJoin struct {
	LeftSet  QuerySet
	RightSet QuerySet
	OnLeft   ReadWrite
	onRight  ReadWrite
}

func (i InnerJoin) SQL() string {
	return "a"
}
