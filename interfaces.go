package pgtalk

import (
	"io"
)

// assertEnabled, if true then perform extra runtime assertion that may panic
var assertEnabled = false

// EnableAssert will enable running extra, potentially more expensive, assertion checks.
// Use this for running your test code.
func EnableAssert() { assertEnabled = true }

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
	// SQLOn writes a valid SQL on a Writer
	SQLOn(w io.Writer)
}

type SQLExpression interface {
	SQLWriter
	// Collect returns all ColumnAccessor that are used in the expression. It exists for assertion.
	Collect(list []ColumnAccessor) []ColumnAccessor
}
