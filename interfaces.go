package pgtalk

import (
	"io"
)

type NewEntityFunc func() interface{}

type Unwrappable interface {
	Unwrap() QuerySet
}

type ColumnAccessor interface {
	Name() string
	SQLOn(w io.Writer)
	// temp name
	WriteInto(entity interface{}, fieldValue interface{})
	// temp name
	ValueAsSQLOn(w io.Writer)
	// temp name
	InsertValue() interface{}
}

type SQLWriter interface {
	SQLOn(w io.Writer)
}
