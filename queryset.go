package pgtalk

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
)

type QuerySet[T any] struct {
	preparedName       string
	tableInfo          TableInfo
	tableAliasOverride string
	selectors          []ColumnAccessor
	distinct           bool
	condition          SQLExpression
	limit              int
	offset             int
	groupBy            []ColumnAccessor
	having             SQLExpression
	orderBy            []ColumnAccessor
	sortOption         string
}

func MakeQuerySet[T any](tableInfo TableInfo, selectors []ColumnAccessor) QuerySet[T] {
	if assertEnabled {
		assertEachAccessorHasTableInfo(selectors, tableInfo)
	}
	return QuerySet[T]{
		tableInfo: tableInfo,
		selectors: selectors,
		condition: EmptyCondition}
}

// querySet
func (q QuerySet[T]) selectAccessors() []ColumnAccessor { return q.selectors }
func (q QuerySet[T]) whereCondition() SQLExpression     { return q.condition }
func (q QuerySet[T]) fromSectionOn(w WriteContext) {
	fmt.Fprintf(w, "%s.%s %s", q.tableInfo.Schema, q.tableInfo.Name, w.TableAlias(q.tableInfo.Name, q.tableInfo.Alias))
}

func (q QuerySet[T]) augmentedContext(w WriteContext) WriteContext {
	if q.tableAliasOverride != "" {
		return w.WithAlias(q.tableInfo.Name, q.tableAliasOverride)
	}
	return w
}

func (q QuerySet[T]) SQLOn(w WriteContext) {
	w = q.augmentedContext(w)
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
	if q.sortOption != "" {
		fmt.Fprint(w, " ", q.sortOption)
	}
	if q.limit > 0 {
		fmt.Fprintf(w, " LIMIT %d", q.limit)
	}
	if q.offset > 0 {
		fmt.Fprintf(w, " OFFSET %d", q.offset)
	}
}

// TableAlias will override the default table or view alias
func (q QuerySet[T]) TableAlias(alias string) QuerySet[T] { q.tableAliasOverride = alias; return q }

// Named sets the name for preparing the statement
func (q QuerySet[T]) Named(preparedName string) QuerySet[T] { q.preparedName = preparedName; return q }

// Distinct is a SQL instruction
func (q QuerySet[T]) Distinct() QuerySet[T] { q.distinct = true; return q }

// Ascending is a SQL instruction for ASC sort option
func (q QuerySet[T]) Ascending() QuerySet[T] { q.sortOption = "ASC"; return q }

// Descending is a SQL instruction for DESC sort option
func (q QuerySet[T]) Descending() QuerySet[T] { q.sortOption = "DESC"; return q }

// Where is a SQL instruction
func (q QuerySet[T]) Where(condition SQLExpression) QuerySet[T] { q.condition = condition; return q }

// Limit is a SQL instruction
func (q QuerySet[T]) Limit(limit int) QuerySet[T] { q.limit = limit; return q }

// Offset is a SQL instruction
func (q QuerySet[T]) Offset(offset int) QuerySet[T] { q.offset = offset; return q }

// GroupBy is a SQL instruction
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
	}
	return q
}
func (q QuerySet[T]) Exists() unaryExpression {
	return unaryExpression{Operator: "EXISTS", Operand: q}
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

func (d QuerySet[T]) Exec(ctx context.Context, conn Querier) (list []*T, err error) {
	rows, err := conn.Query(ctx, SQL(d))
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		entity := new(T)
		sw := []interface{}{}
		for _, each := range d.selectors {
			sw = append(sw, each.FieldToScan(entity))
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
