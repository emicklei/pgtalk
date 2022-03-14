package pgtalk

import (
	"fmt"
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
	var t T
	panic("expected a " + fmt.Sprintf("%T", t) + " got a " + fmt.Sprintf("%T", operand))
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
