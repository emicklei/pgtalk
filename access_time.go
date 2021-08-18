package pgtalk

import (
	"fmt"
	"io"
	"time"
)

type TimeAccess struct {
	ColumnInfo
	fieldWriter   func(dest interface{}, i *time.Time)
	valueToInsert time.Time
}

func NewTimeAccess(info ColumnInfo, writer func(dest interface{}, i *time.Time)) TimeAccess {
	return TimeAccess{ColumnInfo: info, fieldWriter: writer}
}

func (a TimeAccess) Collect(list []ColumnAccessor) []ColumnAccessor {
	return append(list, a)
}

func (a TimeAccess) SetFieldValue(entity interface{}, fieldValue interface{}) error {
	if fieldValue == nil {
		return nil
	}
	v, ok := fieldValue.(time.Time)
	if !ok {
		return NewValueConversionError(fieldValue, "time.Time")
	}
	a.fieldWriter(entity, &v)
	return nil
}

func (a TimeAccess) ValueAsSQLOn(w io.Writer) {
	fmt.Fprintf(w, "%v", a.valueToInsert)
}

func (a TimeAccess) ValueToInsert() interface{} {
	return a.valueToInsert
}

func (a TimeAccess) Set(v time.Time) TimeAccess {
	a.valueToInsert = v
	return a
}

func (a TimeAccess) Column() ColumnInfo { return a.ColumnInfo }

type BooleanAccess struct {
	ColumnInfo
	fieldWriter   func(dest interface{}, b *bool)
	valueToInsert bool
}

func NewBooleanAccess(info ColumnInfo, writer func(dest interface{}, b *bool)) BooleanAccess {
	return BooleanAccess{ColumnInfo: info, fieldWriter: writer}
}

func (a BooleanAccess) Collect(list []ColumnAccessor) []ColumnAccessor {
	return append(list, a)
}

func (a BooleanAccess) SetFieldValue(entity interface{}, fieldValue interface{}) error {
	if fieldValue == nil {
		return nil
	}
	v, ok := fieldValue.(bool)
	if !ok {
		return NewValueConversionError(fieldValue, "bool")
	}
	a.fieldWriter(entity, &v)
	return nil
}

func (a BooleanAccess) Column() ColumnInfo { return a.ColumnInfo }

func (a BooleanAccess) Set(v bool) BooleanAccess {
	a.valueToInsert = v
	return a
}
func (a BooleanAccess) ValueToInsert() interface{} {
	return a.valueToInsert
}

func (a BooleanAccess) ValueAsSQLOn(w io.Writer) {
	fmt.Fprintf(w, "%v", a.valueToInsert)
}

func (a BooleanAccess) Equals(b bool) SQLExpression {
	return MakeBinaryOperator(a, "=", ValuePrinter{b})
}
