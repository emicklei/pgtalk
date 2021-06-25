package xs

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/jackc/pgx/v4"
)

type NewEntityFunc func() interface{}

type QuerySet struct {
	tableInfo TableInfo
	selectors []ReadWrite
	condition SQLWriter
	limit     int
	factory   NewEntityFunc
}

func MakeQuerySet(tableInfo TableInfo, selectors []ReadWrite, factory NewEntityFunc) QuerySet {
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

func (d QuerySet) Exec(conn *pgx.Conn, appender func(each interface{})) (err error) {
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

type Unwrappable interface {
	Unwrap() QuerySet
}

func (d QuerySet) Join(otherQuerySet Unwrappable) InnerJoin {
	return InnerJoin{
		LeftSet:  d,
		RightSet: otherQuerySet.Unwrap(),
	}
}

type InnerJoin struct {
	LeftSet  QuerySet
	RightSet QuerySet
	OnLeft   ReadWrite
	OnRight  ReadWrite
}

func (i InnerJoin) SQL() string {
	//temp
	return "SELECT " + i.LeftSet.SelectSection() + "," + i.RightSet.SelectSection() +
		" FROM " + i.LeftSet.FromSection() +
		" INNER JOIN " + i.RightSet.FromSection() +
		" ON (" + i.LeftSet.tableInfo.Alias + "." + i.OnLeft.Name() +
		" = " + i.RightSet.tableInfo.Alias + "." + i.OnRight.Name() +
		") WHERE " + i.LeftSet.WhereSection()
}

func (i InnerJoin) On(onLeft, onRight ReadWrite) InnerJoin {
	return InnerJoin{
		LeftSet:  i.LeftSet,
		RightSet: i.RightSet,
		OnLeft:   onLeft,
		OnRight:  onRight,
	}
}

func (i InnerJoin) Exec(conn *pgx.Conn) (it InnerJoinIterator, err error) {
	rows, err := conn.Query(context.Background(), i.SQL())
	if err != nil {
		return
	}
	return InnerJoinIterator{leftSet: i.LeftSet, rightSet: i.RightSet, rows: rows}, nil
}

type InnerJoinIterator struct {
	leftSet  QuerySet
	rightSet QuerySet
	rows     pgx.Rows
}

func (i *InnerJoinIterator) HasNext() bool {
	if i.rows.Next() {
		return true
	} else {
		i.rows.Close()
	}
	return false
}

func (i *InnerJoinIterator) Next(left interface{}, right interface{}) error {
	sw := []interface{}{}
	// left
	for _, each := range i.leftSet.selectors {
		rw := ScanToWrite{
			RW:     each,
			Entity: left,
		}
		sw = append(sw, rw)
	}
	// right
	for _, each := range i.rightSet.selectors {
		rw := ScanToWrite{
			RW:     each,
			Entity: right,
		}
		sw = append(sw, rw)
	}
	return i.rows.Scan(sw...)
}
