package pgtalk

import (
	"bytes"
	"context"
	"fmt"
	"io"
)

type QuerySet struct {
	tableInfo TableInfo
	selectors []ColumnAccessor
	condition SQLWriter
	limit     int
	factory   NewEntityFunc
}

func MakeQuerySet(tableInfo TableInfo, selectors []ColumnAccessor, factory NewEntityFunc) QuerySet {
	return QuerySet{
		tableInfo: tableInfo,
		selectors: selectors,
		condition: EmptyCondition,
		factory:   factory}
}

func (q QuerySet) FromSection() string {
	return fmt.Sprintf("%s %s", q.tableInfo.Name, q.tableInfo.Alias)
}

func (q QuerySet) SelectSection() string {
	buf := new(bytes.Buffer)
	for i, each := range q.selectors {
		if i > 0 {
			io.WriteString(buf, ",")
		}
		io.WriteString(buf, q.tableInfo.Alias)
		io.WriteString(buf, ".")
		io.WriteString(buf, each.Name())
	}
	return buf.String()
}

func (q QuerySet) WhereSection() string {
	return q.condition.SQL()
}

// SQL returns the full SQL query
func (q QuerySet) SQL() string {
	// TEMP
	where := q.WhereSection()
	if len(where) > 0 {
		where = fmt.Sprintf(" WHERE %s", where)
	}
	limit := ""
	if q.limit > 0 {
		limit = fmt.Sprintf(" LIMIT %d", q.limit)
	}
	return fmt.Sprintf("SELECT %s FROM %s%s%s", q.SelectSection(), q.FromSection(), where, limit)
}

func (q QuerySet) Where(condition SQLWriter) QuerySet {
	return QuerySet{tableInfo: q.tableInfo, selectors: q.selectors, condition: condition, factory: q.factory}
}

func (q QuerySet) Limit(limit int) QuerySet {
	return QuerySet{tableInfo: q.tableInfo, selectors: q.selectors, condition: q.condition, limit: limit, factory: q.factory}
}

func (d QuerySet) Exec(conn Connection, appender func(each interface{})) (err error) {
	rows, err := conn.Query(context.Background(), d.SQL())
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		entity := d.factory()
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

func (d QuerySet) Join(otherQuerySet Unwrappable) Join {
	return Join{
		LeftSet:  d,
		RightSet: otherQuerySet.Unwrap(),
		Type:     InnerJoinType,
	}
}

func (d QuerySet) LeftJoin(otherQuerySet Unwrappable) Join {
	return Join{
		LeftSet:  d,
		RightSet: otherQuerySet.Unwrap(),
		Type:     LeftOuterJoinType,
	}
}

func (d QuerySet) RightJoin(otherQuerySet Unwrappable) Join {
	return Join{
		LeftSet:  d,
		RightSet: otherQuerySet.Unwrap(),
		Type:     RightOuterJoinType,
	}
}

func (d QuerySet) FullJoin(otherQuerySet Unwrappable) Join {
	return Join{
		LeftSet:  d,
		RightSet: otherQuerySet.Unwrap(),
		Type:     FullOuterJoinType,
	}
}
