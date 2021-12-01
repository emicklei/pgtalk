package pgtalk

import ("io";"fmt")

type FieldAccess[T any]struct {
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

func (a FieldAccess[T]) Collect(list []ColumnAccessor) []ColumnAccessor {
	return append(list, a)
}

func (a FieldAccess[T]) SetFieldValue(entity interface{}, fieldValue interface{}) error {
	if fieldValue == nil {
		return nil
	}
	v, ok := fieldValue.(*T)
	if !ok {
		var t *T
		return NewValueConversionError(fieldValue, fmt.Sprintf("%T",t))
	}
	a.fieldWriter(entity, v)
	return nil
}

func (a FieldAccess[T]) ValueToInsert() interface{} {
	return a.valueToInsert
}

func (a FieldAccess[T]) ValueAsSQLOn(w io.Writer) {
	fmt.Fprintf(w, "%v", a.valueToInsert) // TODO
}
