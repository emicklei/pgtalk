package pgtalk

import (
	"fmt"

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
	i.rows.Close()
	return false
}

func (i *ResultIterator[T]) Next() (*T, error) {
	entity := new(T)
	list := i.rows.FieldDescriptions()
	vals, err := i.rows.Values()
	if err != nil {
		return entity, fmt.Errorf("unable to get values:%v", err)
	}
	// order of list is not the same as selectors?
	for f, each := range list {
		for _, other := range i.selectors {
			if other.Column().tableAttributeNumber == each.TableAttributeNumber {
				// TODO error handling
				other.SetFieldValue(entity, vals[f])
			}
		}
	}
	return entity, nil
}
