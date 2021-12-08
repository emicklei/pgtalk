package pgtalk

import (
	"fmt"
)

type FieldAccess[T any] struct {
	ColumnInfo
	fieldWriter   func(dest interface{}, u *T)
	valueToInsert T
}

func NewFieldAccess[T any](
	info ColumnInfo,
	writer func(dest interface{}, u *T)) FieldAccess[T] {
	return FieldAccess[T]{
		ColumnInfo:  info,
		fieldWriter: writer}
}

func (a FieldAccess[T]) Column() ColumnInfo { return a.ColumnInfo }

// Collect is part of SQLExpression
func (a FieldAccess[T]) Collect(list []ColumnAccessor) []ColumnAccessor {
	return append(list, a)
}

// Set returns a new FieldAccess[T] with a value to set on a T.
func (a FieldAccess[T]) Set(v T) FieldAccess[T] {
	a.valueToInsert = v
	return a
}

func (a FieldAccess[T]) SetFieldValue(entity interface{}, fieldValue interface{}) error {
	if fieldValue == nil {
		return nil
	}
	// try both value of T and pointer to T
	v, ok := fieldValue.(T)
	if ok {
		a.fieldWriter(entity, &v)
		return nil
	}
	pv, ok := fieldValue.(*T)
	if !ok {
		var t *T
		return NewValueConversionError(fieldValue, fmt.Sprintf("%T", t))
	}
	a.fieldWriter(entity, pv)
	return nil
}

func (a FieldAccess[T]) ValueToInsert() interface{} {
	return a.valueToInsert
}

// Equals return a
func (a FieldAccess[T]) Equals(operand interface{}) SQLExpression {
	if fat, ok := operand.(FieldAccess[T]); ok {
		return MakeBinaryOperator(a, "=", fat)
	}
	if t, ok := operand.(T); ok {
		return MakeBinaryOperator(a, "=", valuePrinter{t})
	}
	var t T
	panic("expected a " + fmt.Sprintf("%T", t))
}
