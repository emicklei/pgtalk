package pgtalk

import (
	"context"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

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

type SQLExpression interface {
	SQLWriter
	And(expr SQLExpression) SQLExpression
	Or(expr SQLExpression) SQLExpression
}

type querySet interface {
	fromSectionOn(w WriteContext)
	selectAccessors() []ColumnAccessor
	whereCondition() SQLWriter
	augmentedContext(w WriteContext) WriteContext
}

type querier interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
}

type preparer interface {
	Prepare(ctx context.Context, name, sql string) (sd *pgconn.StatementDescription, err error)
}

type fieldAccessFunc = func(entity any) any

type expressionValueHolder interface {
	AddExpressionResult(key string, value any)
}
