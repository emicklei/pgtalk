package pgtalk

import (
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type resultIterator[T any] struct {
	queryError error
	commandTag pgconn.CommandTag
	rows       pgx.Rows
	selectors  []ColumnAccessor
	params     []any
}

// Close closes the rows, making the connection ready for use again. It is safe
// to call Close after rows is already closed.
func (i *resultIterator[T]) Close() {
	i.rows.Close()
}

func (i *resultIterator[T]) Err() error {
	if i.queryError != nil {
		return i.queryError
	}
	// happens when query is executed with Exec
	if i.rows == nil {
		return nil
	}
	return i.rows.Err()
}

// CommandTag is valid if the query is an Exec query, i.e. not returning rows.
func (i *resultIterator[T]) CommandTag() pgconn.CommandTag {
	return i.commandTag
}

// HasNext returns true if there are more rows to scan.
// If none are left, it closes the rows.
func (i *resultIterator[T]) HasNext() bool {
	if i.queryError != nil {
		return false
	}
	// happens when query is executed with Exec
	if i.rows == nil {
		return false
	}
	if i.rows.Next() {
		return true
	}
	// if Next returns false we can close the rows
	i.rows.Close()
	return false
}

func (i *resultIterator[T]) Next() (*T, error) {
	entity := new(T)
	list := i.rows.FieldDescriptions()
	// order of list is not the same as selectors?
	toScan := []any{}
	for _, each := range list {
		for _, other := range i.selectors {
			if other.Column().columnName == each.Name {
				toScan = append(toScan, other.FieldValueToScan(entity))
			}
		}
	}
	if err := i.rows.Scan(toScan...); err != nil {
		return nil, err
	}
	return entity, nil
}

// GetParams returns all the parameters used in the query. Can be used for debugging or logging
func (i *resultIterator[T]) GetParams() map[int]any {
	ret := make(map[int]any, len(i.params))
	for i, each := range i.params {
		ret[i+1] = each
	}
	return ret
}
