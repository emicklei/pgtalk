package pgtalk

import (
	"github.com/jackc/pgx/v5"
)

type resultIterator[T any] struct {
	queryError error
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
	return i.rows.Err()
}

// HasNext returns true if there are more rows to scan.
// If none are left, it closes the rows.
func (i *resultIterator[T]) HasNext() bool {
	if i.queryError != nil {
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
func (i *resultIterator[T]) GetParams() []any {
	return i.params
}
