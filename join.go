package pgtalk

import (
	"context"
	"fmt"
	"io"

	"github.com/jackc/pgx/v4"
)

type JoinType int

const (
	InnerJoinType JoinType = iota
	LeftOuterJoinType
	RightOuterJoinType
	FullOuterJoinType
)

type Join struct {
	preparedName string
	LeftSet      QuerySet
	RightSet     QuerySet
	OnLeft       ColumnAccessor
	OnRight      ColumnAccessor
	Condition    SQLWriter
	Type         JoinType
}

func (i Join) SQLOn(w io.Writer) {
	fmt.Fprint(w, "SELECT ")
	writeAccessOn(i.LeftSet.selectors, w)
	fmt.Fprint(w, ",")
	writeAccessOn(i.RightSet.selectors, w)
	fmt.Fprint(w, " FROM ")
	i.LeftSet.fromSectionOn(w)
	writeJoinType(i.Type, w)
	i.RightSet.fromSectionOn(w)
	fmt.Fprint(w, " ON ")
	i.Condition.SQLOn(w)
	if _, ok := i.LeftSet.condition.(NoCondition); !ok {
		fmt.Fprint(w, " WHERE ")
		i.LeftSet.condition.SQLOn(w)
	}
	// TODO RightSet where
}

func writeJoinType(t JoinType, w io.Writer) {
	switch t {
	case InnerJoinType:
		fmt.Fprint(w, " INNER JOIN ")
	case LeftOuterJoinType:
		fmt.Fprint(w, " LEFT OUTER JOIN ")
	case RightOuterJoinType:
		fmt.Fprint(w, " RIGHT OUTER JOIN ")
	case FullOuterJoinType:
		fmt.Fprint(w, " FULL OUTER JOIN ")
	}
}

func (i Join) Named(preparedName string) Join {
	return Join{
		preparedName: preparedName,
		LeftSet:      i.LeftSet,
		RightSet:     i.RightSet,
		Condition:    i.Condition,
		Type:         i.Type,
	}
}

func (i Join) On(condition SQLWriter) Join {
	return Join{
		preparedName: i.preparedName,
		LeftSet:      i.LeftSet,
		RightSet:     i.RightSet,
		Condition:    condition,
		Type:         i.Type,
	}
}

func (i Join) LeftOuterJoin(q Unwrappable) (m MultiJoin) {
	m.Sets = append(m.Sets, i.LeftSet, i.RightSet, q.Unwrap())
	m.JoinTypes = append(m.JoinTypes, i.Type, LeftOuterJoinType)
	m.Conditions = append(m.Conditions, i.Condition)
	return
}

func (i Join) Exec(ctx context.Context, conn *pgx.Conn) (it JoinResultIterator, err error) {
	sql := SQL(i)
	if i.preparedName != "" {
		_, err := conn.Prepare(ctx, i.preparedName, sql)
		if err != nil {
			return JoinResultIterator{queryError: err}, nil
		}
	}
	rows, err := conn.Query(ctx, sql)
	if err != nil {
		return
	}
	return JoinResultIterator{leftSet: i.LeftSet, rightSet: i.RightSet, rows: rows}, nil
}

type JoinResultIterator struct {
	queryError error
	leftSet    QuerySet
	rightSet   QuerySet
	rows       pgx.Rows
}

func (i *JoinResultIterator) HasNext() bool {
	if i.rows.Next() {
		return true
	} else {
		i.rows.Close()
	}
	return false
}

func (i *JoinResultIterator) Err() error {
	if i.queryError != nil {
		return i.queryError
	}
	return i.rows.Err()
}

func (i *JoinResultIterator) Next(left interface{}, right interface{}) error {
	sw := []interface{}{}
	// left
	for _, each := range i.leftSet.selectors {
		rw := scanToWrite{
			access: each,
			entity: left,
		}
		sw = append(sw, rw)
	}
	// right
	for _, each := range i.rightSet.selectors {
		rw := scanToWrite{
			access: each,
			entity: right,
		}
		sw = append(sw, rw)
	}
	return i.rows.Scan(sw...)
}

type MultiJoin struct {
	Sets       []QuerySet
	JoinTypes  []JoinType
	Conditions []SQLWriter
}

func (m MultiJoin) On(condition SQLWriter) MultiJoin {
	m.Conditions = append(m.Conditions, condition)
	return m
}

func (m MultiJoin) SQLOn(w io.Writer) {
	fmt.Fprint(w, "SELECT ")
	for i, each := range m.Sets {
		if i > 0 {
			fmt.Fprint(w, ",")
		}
		writeAccessOn(each.selectors, w)
	}
	fmt.Fprint(w, " FROM ")
	first := m.Sets[0]
	first.fromSectionOn(w)
	for j := 0; j < len(m.JoinTypes); j++ {
		jt := m.JoinTypes[j]
		writeJoinType(jt, w)
		set := m.Sets[j+1]
		set.fromSectionOn(w)
		fmt.Fprint(w, " ON ")
		m.Conditions[j].SQLOn(w)
	}
}
