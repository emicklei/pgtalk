package pgtalk

import (
	"context"
	"fmt"
	"io"

	"github.com/jackc/pgx/v4"
)

type QuerySet[T any] struct {
	preparedName string
	tableAccess  TableAccessor
	selectors    []ColumnAccessor
	distinct     bool
	condition    SQLExpression
	limit        int
	groupBy      []ColumnAccessor
	having       SQLExpression
	orderBy      []ColumnAccessor
	sortOrder    string
}

func MakeQuerySet[T any](tableAccess TableAccessor, selectors []ColumnAccessor) QuerySet[T] {
	if assertEnabled {
		assertEachAccessorHasTableInfo(selectors, tableAccess.TableInfo)
	}
	return QuerySet[T]{
		tableAccess: tableAccess,
		selectors:   selectors,
		condition:   EmptyCondition}
}

// querySet
func (q QuerySet[T]) selectAccessors() []ColumnAccessor { return q.selectors }
func (q QuerySet[T]) whereCondition() SQLExpression     { return q.condition }
func (q QuerySet[T]) fromSectionOn(w io.Writer) {
	fmt.Fprintf(w, "%s.%s %s", q.tableAccess.TableInfo.Schema, q.tableAccess.TableInfo.Name, q.tableAccess.TableInfo.Alias)
}

func (q QuerySet[T]) SQLOn(w io.Writer) {
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
	if q.having != nil {
		fmt.Fprint(w, " HAVING ")
		q.having.SQLOn(w)
	}
	if len(q.orderBy) > 0 {
		fmt.Fprint(w, " ORDER BY ")
		writeAccessOn(q.orderBy, w)
	}
	if q.limit > 0 {
		fmt.Fprintf(w, " LIMIT %d", q.limit)
	}
}

func (q QuerySet[T]) Named(preparedName string) QuerySet[T]     { q.preparedName = preparedName; return q }
func (q QuerySet[T]) Distinct() QuerySet[T]                     { q.distinct = true; return q }
func (q QuerySet[T]) Ascending() QuerySet[T]                    { q.sortOrder = "ASC"; return q }
func (q QuerySet[T]) Descending() QuerySet[T]                   { q.sortOrder = "DESC"; return q }
func (q QuerySet[T]) Where(condition SQLExpression) QuerySet[T] { q.condition = condition; return q }
func (q QuerySet[T]) Limit(limit int) QuerySet[T]               { q.limit = limit; return q }
func (q QuerySet[T]) GroupBy(cas ...ColumnAccessor) QuerySet[T] {
	q.groupBy = cas
	if assertEnabled {
		assertEachAccessorIn(cas, q.selectors)
	}
	return q
}
func (q QuerySet[T]) Having(condition SQLExpression) QuerySet[T] { q.having = condition; return q }
func (q QuerySet[T]) OrderBy(cas ...ColumnAccessor) QuerySet[T] {
	q.orderBy = cas
	if assertEnabled {
		assertEachAccessorHasTableInfo(cas, q.tableAccess.TableInfo)
	}
	return q
}
func (q QuerySet[T]) Exists() UnaryOperator {
	return UnaryOperator{Operator: "EXISTS", Operand: q}
}

// Collect is part of SQLExpression
func (d QuerySet[T]) Collect(list []ColumnAccessor) []ColumnAccessor {
	return list // TODO
}

func (d QuerySet[T]) Iterate(ctx context.Context, conn *pgx.Conn) (*ResultIterator[T], error) {
	rows, err := conn.Query(ctx, SQL(d))
	return &ResultIterator[T]{
		queryError: err,
		rows:       rows,
		selectors:  d.selectors,
	}, err
}

func (d QuerySet[T]) Exec(ctx context.Context, conn *pgx.Conn) (list []*T, err error) {
	rows, err := conn.Query(ctx, SQL(d))
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		entity := new(T)
		sw := []interface{}{}
		for _, each := range d.selectors {
			rw := scanToWrite{
				access: each,
				entity: entity,
			}
			sw = append(sw, rw)
		}
		if err := rows.Scan(sw...); err != nil {
			return list, err
		}
		list = append(list, entity)
	}
	return
}

func (d QuerySet[T]) Join(otherQuerySet querySet) Join {
	return Join{
		leftSet:  d,
		rightSet: otherQuerySet,
		joinType: InnerJoinType,
	}
}

func (d QuerySet[T]) LeftOuterJoin(otherQuerySet querySet) Join {
	return Join{
		leftSet:  d,
		rightSet: otherQuerySet,
		joinType: LeftOuterJoinType,
	}
}

func (d QuerySet[T]) RightJoin(otherQuerySet querySet) Join {
	return Join{
		leftSet:  d,
		rightSet: otherQuerySet,
		joinType: RightOuterJoinType,
	}
}

func (d QuerySet[T]) FullJoin(otherQuerySet querySet) Join {
	return Join{
		leftSet:  d,
		rightSet: otherQuerySet,
		joinType: FullOuterJoinType,
	}
}
