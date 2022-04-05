package pgtalk

import (
	"context"
	"fmt"
	"reflect"
)

type QuerySet[T any] struct {
	unimplementedBooleanExpression
	preparedName       string
	tableInfo          TableInfo
	tableAliasOverride string
	selectors          []ColumnAccessor
	distinct           bool
	condition          SQLWriter
	limit              int
	offset             int
	groupBy            []ColumnAccessor
	having             SQLExpression
	orderBy            []ColumnAccessor
	sortOption         string
}

func MakeQuerySet[T any](tableInfo TableInfo, selectors []ColumnAccessor) QuerySet[T] {
	return QuerySet[T]{
		tableInfo: tableInfo,
		selectors: selectors,
		condition: EmptyCondition}
}

// querySet
func (q QuerySet[T]) selectAccessors() []ColumnAccessor { return q.selectors }
func (q QuerySet[T]) whereCondition() SQLWriter         { return q.condition }
func (q QuerySet[T]) fromSectionOn(w writeContext) {
	fmt.Fprintf(w, "%s.%s %s", q.tableInfo.Schema, q.tableInfo.Name, w.TableAlias(q.tableInfo.Name, q.tableInfo.Alias))
}

func (q QuerySet[T]) augmentedContext(w writeContext) writeContext {
	if q.tableAliasOverride != "" {
		return w.WithAlias(q.tableInfo.Name, q.tableAliasOverride)
	}
	return w
}

func (q QuerySet[T]) SQLOn(w writeContext) {
	w = q.augmentedContext(w)
	fmt.Fprint(w, "SELECT")
	if q.distinct {
		fmt.Fprint(w, " DISTINCT\n")
	} else {
		fmt.Fprint(w, "\n")
	}
	writeAccessOn(q.selectors, w)
	fmt.Fprint(w, "\nFROM ")
	q.fromSectionOn(w)
	if _, ok := q.condition.(noCondition); !ok {
		fmt.Fprint(w, "\nWHERE ")
		q.condition.SQLOn(w)
	}
	if len(q.groupBy) > 0 {
		fmt.Fprint(w, "\nGROUP BY\n")
		writeAccessOn(q.groupBy, w)
	}
	if q.having != nil {
		fmt.Fprint(w, "\nHAVING ")
		q.having.SQLOn(w)
	}
	if len(q.orderBy) > 0 {
		fmt.Fprint(w, "\nORDER BY\n")
		writeAccessOn(q.orderBy, w)
	}
	if q.sortOption != "" {
		fmt.Fprint(w, " ", q.sortOption)
	}
	if q.limit > 0 {
		fmt.Fprintf(w, "\nLIMIT %d", q.limit)
	}
	if q.offset > 0 {
		fmt.Fprintf(w, "\nOFFSET %d", q.offset)
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
	return q
}
func (q QuerySet[T]) Having(condition SQLExpression) QuerySet[T] { q.having = condition; return q }
func (q QuerySet[T]) OrderBy(cas ...ColumnAccessor) QuerySet[T] {
	q.orderBy = cas
	return q
}

func (q QuerySet[T]) Exists() unaryExpression {
	return unaryExpression{Operator: "EXISTS", Operand: q}
}

func (d QuerySet[T]) Iterate(ctx context.Context, conn querier, parameters ...*QueryParameter) (*resultIterator[T], error) {
	params := argumentValues(parameters)
	rows, err := conn.Query(ctx, SQL(d), params...)
	return &resultIterator[T]{
		queryError: err,
		rows:       rows,
		selectors:  d.selectors,
	}, err
}

func (d QuerySet[T]) Exec(ctx context.Context, conn querier, parameters ...*QueryParameter) (list []*T, err error) {
	params := argumentValues(parameters)
	rows, err := conn.Query(ctx, SQL(d), params...)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		entity := new(T)
		sw := []any{}
		for _, each := range d.selectors {
			sw = append(sw, each.FieldValueToScan(entity))
		}
		if err := rows.Scan(sw...); err != nil {
			return list, err
		}
		list = append(list, entity)
	}
	return
}

func (d QuerySet[T]) ExecIntoMaps(ctx context.Context, conn querier, parameters ...*QueryParameter) (list []map[string]any, err error) {
	return execIntoMaps(ctx, conn, SQL(d), d.selectors, parameters...)
}

func execIntoMaps(ctx context.Context, conn querier, query string, selectors []ColumnAccessor, parameters ...*QueryParameter) (list []map[string]any, err error) {
	params := argumentValues(parameters)
	rows, err := conn.Query(ctx, query, params...)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		sw := []any{} // sw holds addresses to the valueToInsert
		for _, each := range selectors {
			sw = each.AppendScannable(sw)
		}
		if err := rows.Scan(sw...); err != nil {
			return list, err
		}
		row := map[string]any{}
		for i, each := range selectors {
			// sw[i] is the address of the valueToInsert of each (ColumnAccessor)
			// use reflect version of dereferencing
			rv := reflect.ValueOf(sw[i])
			row[each.Column().columnName] = rv.Elem().Interface()
		}
		list = append(list, row)
	}
	return
}

func (d QuerySet[T]) Join(otherQuerySet querySet) join {
	return join{
		leftSet:  d,
		rightSet: otherQuerySet,
		joinType: innerJoinType,
	}
}

func (d QuerySet[T]) LeftOuterJoin(otherQuerySet querySet) join {
	return join{
		leftSet:  d,
		rightSet: otherQuerySet,
		joinType: leftOuterJoinType,
	}
}

func (d QuerySet[T]) RightJoin(otherQuerySet querySet) join {
	return join{
		leftSet:  d,
		rightSet: otherQuerySet,
		joinType: rightOuterJoinType,
	}
}

func (d QuerySet[T]) FullJoin(otherQuerySet querySet) join {
	return join{
		leftSet:  d,
		rightSet: otherQuerySet,
		joinType: fullOuterJoinType,
	}
}
