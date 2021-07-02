package pgtalk

import (
	"context"
	"io"

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
	SQLOn(w io.Writer)
	// temp name
	WriteInto(entity interface{}, fieldValue interface{})
	// temp name
	ValueAsSQLOn(w io.Writer)
}

type SQLWriter interface {
	SQLOn(w io.Writer)
}
