package pgtalk

import (
	"context"
	"fmt"
	"io"

	"github.com/jackc/pgx/v4"
)

type QuerySet struct {
	preparedName string
	tableInfo    TableInfo
	selectors    []ColumnAccessor
	distinct     bool
	condition    SQLWriter
	limit        int
	factory      NewEntityFunc
	groupBy      []ColumnAccessor
	having       SQLWriter
	orderBy      []ColumnAccessor
	sortOrder    string
	fieldSetter  FieldSetter
}

func MakeQuerySet(tableInfo TableInfo, selectors []ColumnAccessor, factory NewEntityFunc, fieldSetter FieldSetter) QuerySet {
	return QuerySet{
		tableInfo:   tableInfo,
		selectors:   selectors,
		condition:   EmptyCondition,
		factory:     factory,
		fieldSetter: fieldSetter}
}

func (q QuerySet) fromSectionOn(w io.Writer) {
	fmt.Fprintf(w, "%s.%s %s", q.tableInfo.Schema, q.tableInfo.Name, q.tableInfo.Alias)
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

func (q QuerySet) Named(preparedName string) QuerySet             { q.preparedName = preparedName; return q }
func (q QuerySet) Distinct() QuerySet                             { q.distinct = true; return q }
func (q QuerySet) Ascending() QuerySet                            { q.sortOrder = "ASC"; return q }
func (q QuerySet) Descending() QuerySet                           { q.sortOrder = "DESC"; return q }
func (q QuerySet) Where(condition SQLWriter) QuerySet             { q.condition = condition; return q }
func (q QuerySet) WhereExists(otherQuerySet Unwrappable) QuerySet { return q } // TODO
func (q QuerySet) Limit(limit int) QuerySet                       { q.limit = limit; return q }
func (q QuerySet) GroupBy(cas ...ColumnAccessor) QuerySet         { q.groupBy = cas; return q }
func (q QuerySet) Having(condition SQLWriter) QuerySet            { q.having = condition; return q }
func (q QuerySet) OrderBy(cas ...ColumnAccessor) QuerySet         { q.orderBy = cas; return q }
func (q QuerySet) Exists() UnaryOperator {
	return UnaryOperator{Operator: "EXISTS", Operand: q}
}

func (d QuerySet) Exec(ctx context.Context, conn *pgx.Conn) *ResultIterator {
	sql := SQL(d)
	var rows pgx.Rows
	var err error
	if d.preparedName != "" {
		_, err := conn.Prepare(ctx, d.preparedName, sql)
		if err != nil {
			return &ResultIterator{queryError: err}
		}
		rows, err = conn.Query(ctx, d.preparedName)
	} else {
		rows, err = conn.Query(ctx, sql)
	}
	return &ResultIterator{queryError: err, rows: rows, fieldSetter: d.fieldSetter}
}

type ResultIterator struct {
	queryError  error
	rows        pgx.Rows
	fieldSetter FieldSetter
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
	vals, err := i.rows.Values()
	if err != nil {
		return fmt.Errorf("unable to get values:%v", err)
	}
	for f, each := range list {
		if err := i.fieldSetter(entity, each.TableAttributeNumber, vals[f]); err != nil {
			return err
		}
	}
	return nil
}

func (d QuerySet) ExecWithAppender(ctx context.Context, conn *pgx.Conn, appender func(each interface{})) (err error) {
	rows, err := conn.Query(ctx, SQL(d))
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
		leftSet:  d,
		rightSet: otherQuerySet.Unwrap(),
		joinType: InnerJoinType,
	}
}

func (d QuerySet) LeftOuterJoin(otherQuerySet Unwrappable) Join {
	return Join{
		leftSet:  d,
		rightSet: otherQuerySet.Unwrap(),
		joinType: LeftOuterJoinType,
	}
}

func (d QuerySet) RightJoin(otherQuerySet Unwrappable) Join {
	return Join{
		leftSet:  d,
		rightSet: otherQuerySet.Unwrap(),
		joinType: RightOuterJoinType,
	}
}

func (d QuerySet) FullJoin(otherQuerySet Unwrappable) Join {
	return Join{
		leftSet:  d,
		rightSet: otherQuerySet.Unwrap(),
		joinType: FullOuterJoinType,
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
func (c Count) InsertValue() interface{}   { return nil }
