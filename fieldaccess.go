package pgtalk

import (
	"fmt"
	"strings"
)

type FieldAccess[T any] struct {
	unimplementedBooleanExpression
	ColumnInfo
	valueFieldWriter fieldAccessFunc
	valueToInsert    T
}

func NewFieldAccess[T any](
	info ColumnInfo,
	writer func(dest any) any) FieldAccess[T] {
	return FieldAccess[T]{
		ColumnInfo:       info,
		valueFieldWriter: writer}
}

func (a FieldAccess[T]) Column() ColumnInfo { return a.ColumnInfo }

func (a FieldAccess[T]) FieldValueToScan(entity any) any {
	return a.valueFieldWriter(entity)
}

// Set returns a new FieldAccess[T] with a value to set on a T.
func (a FieldAccess[T]) Set(v T) FieldAccess[T] {
	a.valueToInsert = v
	return a
}

// Get returns the value for its columnName from a map (row).
func (a FieldAccess[T]) Get(values map[string]any) any {
	v, ok := values[a.columnName]
	if !ok {
		var none T
		return none
	}
	return v
}

func (a FieldAccess[T]) ValueToInsert() any {
	return a.valueToInsert
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

// Equals returns a SQLExpression
func (a FieldAccess[T]) Equals(operand any) binaryExpression {
	if fat, ok := operand.(FieldAccess[T]); ok {
		return makeBinaryOperator(a, "=", fat)
	}
	if t, ok := operand.(T); ok {
		return makeBinaryOperator(a, "=", valuePrinter{v: t})
	}
	if anyp, ok := operand.(*QueryParameter); ok {
		if _, ok := anyp.value.(T); ok {
			return makeBinaryOperator(a, "=", anyp) // use parameter, not its value
		}
	}
	return makeBinaryOperator(a, "=", valuePrinter{v: operand})
}

// Less returns a SQLExpression
func (a FieldAccess[T]) LessThan(operand any) binaryExpression {
	if fat, ok := operand.(FieldAccess[T]); ok {
		return makeBinaryOperator(a, "<", fat)
	}
	if t, ok := operand.(T); ok {
		return makeBinaryOperator(a, "<", valuePrinter{v: t})
	}
	var t T
	panic("expected a " + fmt.Sprintf("%T", t) + " got a " + fmt.Sprintf("%T", operand))
}

// In returns a binary expression to check that the value of the fieldAccess is in the values collection.
func (a FieldAccess[T]) In(values ...T) binaryExpression {
	vs := make([]any, len(values))
	for i := 0; i < len(values); i++ {
		vs[i] = values[i]
	}
	return makeBinaryOperator(a, "IN", valuesPrinter{vs: vs})
}

func (a FieldAccess[T]) Compare(operator string, operand any) binaryExpression {
	if !strings.Contains(validComparisonOperators, operator) {
		panic("invalid comparison operator:" + operator)
	}
	return makeBinaryOperator(a, operator, valuePrinter{v: operand})
}

// TableAlias changes the table alias for this column accessor.
func (a FieldAccess[T]) TableAlias(alias string) FieldAccess[T] {
	a.ColumnInfo = a.ColumnInfo.TableAlias(alias)
	return a
}

// IsNull creates a SQL Expresion with IS NULL.
func (a FieldAccess[T]) IsNull() binaryExpression {
	return makeBinaryOperator(a, "IS", valuePrinter{v: SQLLiteral{Literal: "NULL"}})
}

// IsNotNull creates a SQL Expresion with IS NOT NULL.
func (a FieldAccess[T]) IsNotNull() binaryExpression {
	return makeBinaryOperator(a, "IS NOT", valuePrinter{v: SQLLiteral{Literal: "NULL"}})
}

// AppendScannable is part of ColumnAccessor
func (a FieldAccess[T]) AppendScannable(list []any) []any {
	return append(list, &a.valueToInsert)
}
