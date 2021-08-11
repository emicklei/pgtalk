package pgtalk

import (
	"fmt"

	"github.com/jackc/pgx/v4"
)

type ResultIterator struct {
	queryError error
	rows       pgx.Rows
	selectors  []ColumnAccessor
}

func (i *ResultIterator) Err() error {
	return i.queryError
}

func (i *ResultIterator) HasNext() bool {
	if i.queryError != nil {
		return false
	}
	if i.rows.Next() {
		return true
	}
	i.rows.Close()
	return false
}

func (i *ResultIterator) Next(entity interface{}) error {
	list := i.rows.FieldDescriptions()
	vals, err := i.rows.Values()
	if err != nil {
		return fmt.Errorf("unable to get values:%v", err)
	}
	// order of list is not the same as selectors?
	for f, each := range list {
		for _, other := range i.selectors {
			if other.Column().tableAttributeNumber == each.TableAttributeNumber {
				// TODO error handling
				other.WriteInto(entity, vals[f])
			}
		}
	}
	return nil
}
