package pgtalk

import (
	"context"
	"fmt"
)

// Work in Progress

type UntypedQuerySet struct {
	expressions []ColumnAccessor
	set         SQLWriter
}

func (c UntypedQuerySet) SQLOn(w WriteContext) {
	fmt.Fprintf(w, "SELECT ")
	subW := w.WithAlias("things", "bag")
	for i, each := range c.expressions {
		if i > 0 {
			fmt.Fprint(w, ",")
		}
		each.SQLOn(subW)
	}
	fmt.Fprintf(w, " FROM (")
	c.set.SQLOn(w)
	fmt.Fprintf(w, ") AS bag")
}

func (c UntypedQuerySet) ExecIntoMaps(ctx context.Context, conn Querier) (list []map[string]any, err error) {
	return execIntoMaps(ctx, conn, SQL(c), c.expressions)
}

func (d QuerySet[T]) Collect(expressions ...ColumnAccessor) UntypedQuerySet {
	return UntypedQuerySet{
		set:         d,
		expressions: expressions,
	}
}

func (a FieldAccess[T]) Concat(resultName string, ex SQLExpression) ColumnAccessor {
	return ComputedColumn{
		ResultName: resultName,
		Expression: binaryExpression{
			Left:     a,
			Operator: "||",
			Right:    ex,
		}}
}

// ComputedColumn is a ColumnAccessor for read.
type ComputedColumn struct {
	ResultName string
	Expression SQLExpression
	Value      any
}

func (c ComputedColumn) SQLOn(w WriteContext) {
	c.Expression.SQLOn(w)
	fmt.Fprintf(w, " AS %s", c.ResultName)
}
func (c ComputedColumn) Name() string       { return c.ResultName }
func (c ComputedColumn) ValueToInsert() any { return nil }
func (c ComputedColumn) Column() ColumnInfo { return ColumnInfo{columnName: c.ResultName} }

// FieldValueToScan returns the address of the value of the field in the entity
func (c ComputedColumn) FieldValueToScan(entity any) any {
	return &c.Value
}

// AppendScannable collects values for scanning by a result Row
// Cannot use ValueToInsert because that looses type information such that the Scanner will use default mapping
func (c ComputedColumn) AppendScannable(list []any) []any {
	return append(list, &c.Value)
}

// Get accesses the value from a map.
func (c ComputedColumn) Get(values map[string]any) any {
	v, ok := values[c.ResultName]
	if !ok {
		return nil
	}
	return v
}
