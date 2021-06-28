package pgtalk

import (
	"context"

	"github.com/jackc/pgx/v4"
)

type Connection interface {
	// Query executes sql with args. If there is an error the returned Rows will be returned in an error state. So it is
	// allowed to ignore the error returned from Query and handle it in Rows.
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
}

type NewEntityFunc func() interface{}

type Unwrappable interface {
	Unwrap() QuerySet
}

type ColumnAccessor interface {
	Name() string
	// temp name
	WriteInto(entity interface{}, fieldValue interface{})
	// temp name
	ValueAsSQL() string
}
