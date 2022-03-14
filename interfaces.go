package pgtalk

import (
	"context"
	"io"

	"github.com/jackc/pgx/v4"
)

// assertEnabled, if true then perform extra runtime assertion that may panic
var assertEnabled = false

// EnableAssert will enable running extra, potentially more expensive, assertion checks.
// Use this for running your test code.
func EnableAssert() { assertEnabled = true }

type NewEntityFunc func() interface{}

type ColumnAccessor interface {
	Name() string
	SQLOn(w io.Writer)
	ValueToInsert() interface{}
	Column() ColumnInfo
	// FieldToScan return the address of the value of the field in the entity
	FieldToScan(entity any) any
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

type querySet interface {
	fromSectionOn(io.Writer)
	selectAccessors() []ColumnAccessor
	whereCondition() SQLExpression
}

type Querier interface {
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
}

type FieldAccessFunc = func(entity any) any
