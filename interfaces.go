package pgtalk

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
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
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
}

type preparer interface {
	Prepare(ctx context.Context, name, sql string) (sd *pgconn.StatementDescription, err error)
}

type fieldAccessFunc = func(entity any) any

type expressionValueHolder interface {
	AddExpressionResult(key string, value any)
}

// ResultIterator is returned from executing a Query (or Mutation).
type ResultIterator[T any] interface {
	// Close closes the rows of the iterator, making the connection ready for use again. It is safe
	// to call Close after rows is already closed.
	// Close is called implicitly when no return results are expected.
	Close()
	// Err returns the Query error if any
	Err() error
	// HasNext returns true if a more results are available. If not then Close is called implicitly.
	HasNext() bool
	// Next returns the next row populated in a T.
	Next() (*T, error)
	// GetParams returns all the parameters used in the query. Can be used for debugging or logging
	GetParams() map[int]any
}
