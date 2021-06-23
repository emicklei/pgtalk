package xs

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/jackc/pgx/v4"
)

type QuerySet struct {
	tableName string
	selectors []ReadWrite
	condition SQLWriter
	limit     int
}

func MakeQuerySet(tableName string, selectors []ReadWrite) QuerySet {
	return QuerySet{tableName: tableName, selectors: selectors, condition: EmptyCondition}
}

// String returns the full SQL query
func (q QuerySet) SQL() string {
	// TODO use GORM here
	buf := new(bytes.Buffer)
	for i, each := range q.selectors {
		if i > 0 {
			io.WriteString(buf, ",")
		}
		io.WriteString(buf, each.Name())
	}
	// TEMP
	where := q.condition.SQL()
	if len(where) > 0 {
		where = fmt.Sprintf(" WHERE %s", where)
	}
	limit := ""
	if q.limit > 0 {
		limit = fmt.Sprintf(" LIMIT %d", q.limit)
	}
	return fmt.Sprintf("SELECT %s FROM %s%s%s", buf, q.tableName, where, limit)
}

func (q QuerySet) Where(condition SQLWriter) QuerySet {
	return QuerySet{tableName: q.tableName, selectors: q.selectors, condition: condition}
}

func (q QuerySet) Limit(limit int) QuerySet {
	return QuerySet{tableName: q.tableName, selectors: q.selectors, condition: q.condition, limit: limit}
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
