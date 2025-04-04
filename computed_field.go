package pgtalk

import (
	"fmt"
)

//	SQLAs returns a ColumnAccessor with a customer SQL expression.
//
// The named result will be available using the GetExpressionResult method of the record type.
func SQLAs(sql, name string) *computedField {
	return &computedField{
		ResultName: name,
		Expression: expressionSource{SQL: sql},
	}
}

type expressionSource struct {
	unimplementedBooleanExpression
	SQL string
}

// SQLOn is part of SQLWriter
func (e expressionSource) SQLOn(w WriteContext) {
	fmt.Fprintf(w, "(%s)", e.SQL)
}

// computedField is a ColumnAccessor for read.
type computedField struct {
	ResultName string
	Expression SQLExpression
}

func (c *computedField) SQLOn(w WriteContext) {
	c.Expression.SQLOn(w)
	fmt.Fprintf(w, " AS %s", c.ResultName)
}
func (c *computedField) Name() string { return c.ResultName }
func (c *computedField) SetSource(parameterIndex int) string {
	return fmt.Sprintf("$%d", parameterIndex)
}
func (c *computedField) ValueToInsert() any { return nil }
func (c *computedField) Column() ColumnInfo { return ColumnInfo{columnName: c.ResultName} }

// FieldValueToScan returns the address of the value of the field in the entity
func (c *computedField) FieldValueToScan(entity any) any {
	var value any
	if h, ok := entity.(expressionValueHolder); ok {
		// side effect to update the entity custom expressions
		h.AddExpressionResult(c.ResultName, &value)
	}
	return &value
}

// AppendScannable collects values for scanning by a result Row
// Cannot use ValueToInsert because that looses type information such that the Scanner will use default mapping
func (c *computedField) AppendScannable(list []any) []any {
	var value any
	return append(list, &value)
}

// Get accesses the value from a map.
func (c *computedField) Get(values map[string]any) any {
	v, ok := values[c.ResultName]
	if !ok {
		return nil
	}
	return v
}
