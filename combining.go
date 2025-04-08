package pgtalk

import (
	"context"
	"fmt"
)

type QueryCombineable interface {
	SQLWriter

	// Every SELECT statement within UNION must have the same number of columns
	// The columns must also have similar data types
	// The columns in every SELECT statement must also be in the same order
	Union(o QueryCombineable, all ...bool) QueryCombineable

	// There are some mandatory rules for INTERSECT operations such as the
	// number of columns, data types, and other columns must be the same
	// in both SELECT statements for the INTERSECT operator to work correctly.
	Intersect(o QueryCombineable, all ...bool) QueryCombineable

	// There are some mandatory rules for EXCEPT operations such as the
	// number of columns, data types, and other columns must be the same
	// in both EXCEPT statements for the EXCEPT operator to work correctly.
	Except(o QueryCombineable, all ...bool) QueryCombineable

	// ExecIntoMaps executes the query and returns each rows as a map string->any .
	// To include table information for each, append a custom sqlfunction to each select statement.
	// For example, adding `pgtalk.SQLAs("'products'", "table")` will add an entry to the map.
	ExecIntoMaps(ctx context.Context, conn querier, parameters ...*QueryParameter) (list []map[string]any, err error)
}

type queryCombination struct {
	left     QueryCombineable
	operator string // UNION,INTERSECT,EXCEPT
	right    QueryCombineable
}

func combineOperator(combiner string, all ...bool) string {
	if len(all) > 0 && all[0] {
		return combiner + " ALL"
	}
	return combiner
}

func (q queryCombination) SQLOn(w WriteContext) {
	fmt.Fprint(w, "((")
	q.left.SQLOn(w)
	fmt.Fprintf(w, ") %s (", q.operator)
	q.right.SQLOn(w)
	fmt.Fprint(w, "))")
}

func (q queryCombination) Union(o QueryCombineable, all ...bool) QueryCombineable {
	return queryCombination{
		left:     q,
		operator: "UNION",
		right:    o,
	}
}
func (q queryCombination) Intersect(o QueryCombineable, all ...bool) QueryCombineable {
	return queryCombination{
		left:     q,
		operator: "INTERSECT",
		right:    o,
	}
}
func (q queryCombination) Except(o QueryCombineable, all ...bool) QueryCombineable {
	return queryCombination{
		left:     q,
		operator: "EXCEPT",
		right:    o,
	}
}

// Combine implements QueryCombineable
func (q queryCombination) ExecIntoMaps(ctx context.Context, conn querier, parameters ...*QueryParameter) (list []map[string]any, err error) {
	return execIntoMaps(ctx, conn, SQL(q), q.findSet().selectAccessors(), parameters...)
}

func (q queryCombination) findSet() querySet {
	if qc, ok := q.left.(queryCombination); ok {
		return qc.findSet()
	}
	return q.left.(querySet)
}
