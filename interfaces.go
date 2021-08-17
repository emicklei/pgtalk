package pgtalk

import (
	"fmt"
	"io"
)

// AssertEnabled, if true then perform extra runtime assertion that may panic
var AssertEnabled = true

type NewEntityFunc func() interface{}

type Unwrappable interface {
	Unwrap() QuerySet
}

type ColumnAccessor interface {
	Name() string
	SQLOn(w io.Writer)
	SetFieldValue(entity interface{}, fieldValue interface{}) error
	ValueAsSQLOn(w io.Writer)
	ValueToInsert() interface{}
	Column() ColumnInfo
}

type SQLWriter interface {
	SQLOn(w io.Writer)
}

type ValueConversionError struct {
	got, want string
}

func NewValueConversionError(got interface{}, want string) error {
	return ValueConversionError{fmt.Sprintf("%T", got), want}
}

func (e ValueConversionError) Error() string {
	return fmt.Sprintf("field value conversion error, got %s expected %s", e.got, e.want)
}

type HasColumn interface {
	Column() ColumnInfo
}
