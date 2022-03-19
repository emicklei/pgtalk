package pgtalk

import (
	"context"

	"github.com/jackc/pgx/v4"
)

// assertEnabled, if true then perform extra runtime assertion that may panic
var assertEnabled = false

// EnableAssert will enable running extra, potentially more expensive, assertion checks.
// Use this for running your test code.
func EnableAssert() { assertEnabled = true }

type NewEntityFunc func() any

type ColumnAccessor interface {
	SQLWriter
	Name() string
	ValueToInsert() any
	Column() ColumnInfo
	// FieldToScan returns the address of the value of the field in the entity
	FieldToScan(entity any) any
}

type SQLWriter interface {
	// SQLOn writes a valid SQL on a Writer in a context
	SQLOn(w WriteContext)
}

type SQLer interface {
	SQL() string
}

type SQLExpression interface {
	SQLWriter
}

type querySet interface {
	fromSectionOn(w WriteContext)
	selectAccessors() []ColumnAccessor
	whereCondition() SQLWriter
	augmentedContext(w WriteContext) WriteContext
}

type Querier interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
}

type FieldAccessFunc = func(entity any) any
