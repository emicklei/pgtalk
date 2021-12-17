package pgtalk

import (
	"bytes"
	"fmt"
	"io"
)

// Experimental

type Query[T any] struct {
	set       QuerySet[T]
	filter    SQLExpression
	accessors []ColumnAccessor
}

func NewQuery[T any](set QuerySet[T], filter SQLExpression) Query[T] {
	return Query[T]{
		set:    set,
		filter: filter}
}

func (q Query[T]) String() string {
	b := new(bytes.Buffer)
	fmt.Fprintln(b, "filter:", q.filter)
	fmt.Fprintln(b, "map:", q.accessors)
	return b.String()
}

func (q Query[T]) Map(cas ...ColumnAccessor) Query[T] {
	return Query[T]{
		set:       q.set,
		filter:    q.filter,
		accessors: cas}
}

func (q Query[T]) SQLOn(w io.Writer) {
	all := q.set.selectors[:]
	all = q.filter.Collect(all)
	all = append(all, q.accessors...)
	for _, each := range all {
		set := MakeQuerySet[ignoreType](each.Column().tableInfo, EmptyColumnAccessor)
		fmt.Fprintln(w, SQL(set))
	}
}

type ignoreType struct{}
