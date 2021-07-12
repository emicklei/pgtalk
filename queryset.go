package pgtalk

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/jackc/pgx/v4"
)

type QuerySet struct {
	tableInfo TableInfo
	selectors []ColumnAccessor
	distinct  bool
	condition SQLWriter
	limit     int
	factory   NewEntityFunc
	groupBy   []ColumnAccessor
	having    SQLWriter
	orderBy   []ColumnAccessor
	sortOrder string
}

func MakeQuerySet(tableInfo TableInfo, selectors []ColumnAccessor, factory NewEntityFunc) QuerySet {
	return QuerySet{
		tableInfo: tableInfo,
		selectors: selectors,
		condition: EmptyCondition,
		factory:   factory}
}

func (q QuerySet) fromSectionOn(w io.Writer) {
	fmt.Fprintf(w, "%s %s", q.tableInfo.Name, q.tableInfo.Alias)
}

func (q QuerySet) SQLOn(w io.Writer) {
	fmt.Fprint(w, "SELECT ")
	if q.distinct {
		fmt.Fprint(w, "DISTINCT ")
	}
	writeAccessOn(q.selectors, w)
	fmt.Fprint(w, " FROM ")
	q.fromSectionOn(w)
	if _, ok := q.condition.(NoCondition); !ok {
		fmt.Fprint(w, " WHERE ")
		q.condition.SQLOn(w)
	}
	if len(q.groupBy) > 0 {
		fmt.Fprint(w, " GROUP BY ")
		writeAccessOn(q.groupBy, w)
	}
	if len(q.orderBy) > 0 {
		fmt.Fprint(w, " ORDER BY ")
		writeAccessOn(q.orderBy, w)
	}
	if q.limit > 0 {
		fmt.Fprintf(w, " LIMIT %d", q.limit)
	}
}

func (q QuerySet) Distinct() QuerySet                          { q.distinct = true; return q }
func (q QuerySet) Ascending() QuerySet                         { q.sortOrder = "ASC"; return q }
func (q QuerySet) Descending() QuerySet                        { q.sortOrder = "DESC"; return q }
func (q QuerySet) Where(condition SQLWriter) QuerySet          { q.condition = condition; return q }
func (q QuerySet) WhereExists(otherQuerySet QuerySet) QuerySet { return q } // TODO
func (q QuerySet) Limit(limit int) QuerySet                    { q.limit = limit; return q }
func (q QuerySet) GroupBy(cas ...ColumnAccessor) QuerySet      { q.groupBy = cas; return q }
func (q QuerySet) Having(condition SQLWriter) QuerySet         { q.having = condition; return q }
func (q QuerySet) OrderBy(cas ...ColumnAccessor) QuerySet      { q.orderBy = cas; return q }

func (d QuerySet) Exec(conn Connection) *ResultIterator {
	rows, err := conn.Query(context.Background(), SQL(d))
	return &ResultIterator{queryError: err, rows: rows}
}

type ResultIterator struct {
	queryError error
	rows       pgx.Rows
}

func (i *ResultIterator) Err() error {
	return i.queryError
}

func (i *ResultIterator) HasNext() bool {
	if i.queryError != nil {
		return false
	}
	if i.rows.Next() {
		return true
	}
	i.rows.Close()
	return false
}

func (i *ResultIterator) Next(entity interface{}) error {
	list := i.rows.FieldDescriptions()
	for _, each := range list {
		log.Println(each.Name)
		log.Println(each.TableAttributeNumber)
	}
	return nil
}

func (d QuerySet) ExecWithAppender(conn Connection, appender func(each interface{})) (err error) {
	rows, err := conn.Query(context.Background(), SQL(d))
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		entity := d.factory()
		sw := []interface{}{}
		for _, each := range d.selectors {
			rw := scanToWrite{
				access: each,
				entity: entity,
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

func (q QuerySet) Count(cas ...ColumnAccessor) QuerySet {
	for _, each := range cas {
		q.selectors = append(q.selectors, Count{accessor: each})
	}
	return q
}

func (d QuerySet) Join(otherQuerySet Unwrappable) Join {
	return Join{
		LeftSet:  d,
		RightSet: otherQuerySet.Unwrap(),
		Type:     InnerJoinType,
	}
}

func (d QuerySet) LeftOuterJoin(otherQuerySet Unwrappable) Join {
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

type Count struct {
	accessor ColumnAccessor
}

func (c Count) Name() string { return c.accessor.Name() }
func (c Count) SQLOn(w io.Writer) {
	fmt.Fprint(w, "COUNT(")
	c.accessor.SQLOn(w)
	fmt.Fprint(w, ")")
}
func (c Count) ValueAsSQLOn(w io.Writer)   {}
func (c Count) WriteInto(e, v interface{}) {}
