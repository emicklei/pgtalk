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
	SetFieldValue(entity interface{}, fieldValue interface{})
	ValueAsSQLOn(w io.Writer)
	ValueToInsert() interface{}
	Column() ColumnInfo
}

type SQLWriter interface {
	SQLOn(w io.Writer)
}
