package pgtalk

import (
	"context"

	"github.com/jackc/pgconn"
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
	// FieldValueToScan returns the address of the value of the field in the entity
	FieldValueToScan(entity any) any
	// AppendScannable collects values for scanning by a result Row
	// Cannot use ValueToInsert because that looses type information such that the Scanner will use default mapping
	AppendScannable(list []any) []any
	// Get accesses the value from a map.
	// (unfortunately, Go methods cannot have additional type parameters:
	// Get[V](values map[string]any) V )
	Get(values map[string]any) any
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

type Preparer interface {
	Prepare(ctx context.Context, name, sql string) (sd *pgconn.StatementDescription, err error)
}

type FieldAccessFunc = func(entity any) any

type ExpressionValueHolder interface {
	AddExpressionResult(key string, value any)
}
