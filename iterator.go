package pgtalk

import (
	"github.com/jackc/pgx/v4"
)

type ResultIterator[T any] struct {
	queryError error
	rows       pgx.Rows
	selectors  []ColumnAccessor
}

func (i *ResultIterator[T]) Err() error {
	return i.queryError
}

func (i *ResultIterator[T]) HasNext() bool {
	if i.queryError != nil {
		return false
	}
	if i.rows.Next() {
		return true
	}
	// is Next returns false we can close the rows
	i.rows.Close()
	return false
}

func (i *ResultIterator[T]) Next() (*T, error) {
	entity := new(T)
	list := i.rows.FieldDescriptions()
	// order of list is not the same as selectors?
	toScan := []any{}
	for _, each := range list {
		for _, other := range i.selectors {
			if other.Column().tableAttributeNumber == each.TableAttributeNumber {
				toScan = append(toScan, other.FieldToScan(entity))
			}
		}
	}
	if err := i.rows.Scan(toScan...); err != nil {
		return nil, err
	}
	return entity, nil
}
