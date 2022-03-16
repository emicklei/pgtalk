package pgtalk

import (
	"fmt"
	"strings"
)

type FieldAccess[T any] struct {
	ColumnInfo
	valueFieldWriter FieldAccessFunc
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

// Collect is part of SQLExpression
func (a FieldAccess[T]) Collect(list []ColumnAccessor) []ColumnAccessor {
	return append(list, a)
}

func (a FieldAccess[T]) FieldToScan(entity any) any {
	return a.valueFieldWriter(entity)
}

// Set returns a new FieldAccess[T] with a value to set on a T.
func (a FieldAccess[T]) Set(v T) FieldAccess[T] {
	a.valueToInsert = v
	return a
}

func (a FieldAccess[T]) ValueToInsert() interface{} {
	return a.valueToInsert
}

// Equals returns a SQLExpression
func (a FieldAccess[T]) Equals(operand interface{}) binaryExpression {
	if fat, ok := operand.(FieldAccess[T]); ok {
		return MakeBinaryOperator(a, "=", fat)
	}
	if t, ok := operand.(T); ok {
		return MakeBinaryOperator(a, "=", valuePrinter{t})
	}
	return MakeBinaryOperator(a, "=", valuePrinter{operand})
}

// Less returns a SQLExpression
func (a FieldAccess[T]) LessThan(operand interface{}) binaryExpression {
	if fat, ok := operand.(FieldAccess[T]); ok {
		return MakeBinaryOperator(a, "<", fat)
	}
	if t, ok := operand.(T); ok {
		return MakeBinaryOperator(a, "<", valuePrinter{t})
	}
	var t T
	panic("expected a " + fmt.Sprintf("%T", t) + " got a " + fmt.Sprintf("%T", operand))
}

func (a FieldAccess[T]) In(values ...any) binaryExpression {
	vs := make([]interface{}, len(values))
	for i := 0; i < len(values); i++ {
		vs[i] = values[i]
	}
	return MakeBinaryOperator(a, "IN", valuesPrinter{vs})
}

func (a FieldAccess[T]) Compare(operator string, operand any) binaryExpression {
	if !strings.Contains(validComparisonOperators, operator) {
		panic("invalid comparison operator:" + operator)
	}
	return MakeBinaryOperator(a, operator, valuePrinter{operand})
}

func (a FieldAccess[T]) TableAlias(alias string) FieldAccess[T] {
	a.ColumnInfo = a.ColumnInfo.TableAlias(alias)
	return a
}
