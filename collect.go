package pgtalk

import (
	"context"
	"fmt"
	"io"
)

// Work in Progress

type UntypedQuerySet struct {
	expressions []ColumnAccessor
	tableInfo   TableInfo
}

func NewUntypedQuerySet(tableInfo TableInfo, expr []ColumnAccessor) UntypedQuerySet {
	return UntypedQuerySet{expressions: expr, tableInfo: tableInfo}
}

func (c UntypedQuerySet) SQLOn(w WriteContext) {
	fmt.Fprintf(w, "SELECT ")
	for i, each := range c.expressions {
		if i > 0 {
			fmt.Fprint(w, ",")
		}
		each.SQLOn(w)
	}
	fmt.Fprintf(w, " FROM ")
	c.tableInfo.SQLOn(w)
}

func (c UntypedQuerySet) ExecIntoMaps(ctx context.Context, conn Querier) (list []map[string]any, err error) {
	return execIntoMaps(ctx, conn, SQL(c), c.expressions)
}

func (a FieldAccess[T]) Concat(resultName string, ex SQLExpression) ColumnAccessor {
	return &computedField{
		ResultName: resultName,
		Expression: binaryExpression{
			Left:     a,
			Operator: "||",
			Right:    ex,
		}}
}

// FieldSQL returns a ColumnAccessor with a customer SQL expressions.
// The named result will be available in the expressionResults map of the record type.
func FieldSQL(sql, name string) *computedField {
	return &computedField{
		ResultName: name,
		Expression: expressionSource{SQL: sql},
	}
}

type expressionSource struct {
	SQL string
}

// SQLOn is part of SQLWriter
func (e expressionSource) SQLOn(w WriteContext) {
	io.WriteString(w, e.SQL)
}

// computedField is a ColumnAccessor for read.
type computedField struct {
	ResultName string
	Expression SQLExpression
	Value      any
}

func (c *computedField) SQLOn(w WriteContext) {
	c.Expression.SQLOn(w)
	fmt.Fprintf(w, " AS %s", c.ResultName)
}
func (c *computedField) Name() string       { return c.ResultName }
func (c *computedField) ValueToInsert() any { return nil }
func (c *computedField) Column() ColumnInfo { return ColumnInfo{columnName: c.ResultName} }

// FieldValueToScan returns the address of the value of the field in the entity
func (c *computedField) FieldValueToScan(entity any) any {
	addr := &c.Value
	if h, ok := entity.(expressionValueHolder); ok {
		// side effect to update the entity custom expressions
		h.AddExpressionResult(c.ResultName, addr)
	}
	return addr
}

// AppendScannable collects values for scanning by a result Row
// Cannot use ValueToInsert because that looses type information such that the Scanner will use default mapping
func (c *computedField) AppendScannable(list []any) []any {
	return append(list, &c.Value)
}

// Get accesses the value from a map.
func (c *computedField) Get(values map[string]any) any {
	v, ok := values[c.ResultName]
	if !ok {
		return nil
	}
	return v
}
