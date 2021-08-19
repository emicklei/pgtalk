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
	leftSet      QuerySet
	rightSet     QuerySet
	onLeft       ColumnAccessor
	onRight      ColumnAccessor
	condition    SQLExpression
	joinType     JoinType
	limit        int
}

func (i Join) SQLOn(w io.Writer) {
	fmt.Fprint(w, "SELECT ")
	writeAccessOn(i.leftSet.selectors, w)
	fmt.Fprint(w, ",")
	writeAccessOn(i.rightSet.selectors, w)
	fmt.Fprint(w, " FROM ")
	i.leftSet.fromSectionOn(w)
	writeJoinType(i.joinType, w)
	i.rightSet.fromSectionOn(w)
	fmt.Fprint(w, " ON ")
	i.condition.SQLOn(w)
	if _, ok := i.leftSet.condition.(NoCondition); !ok {
		fmt.Fprint(w, " WHERE ")
		i.leftSet.condition.SQLOn(w)
	}
	if i.limit > 0 {
		fmt.Fprintf(w, " LIMIT %d", i.limit)
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
	i.preparedName = preparedName
	return i
}

func (i Join) On(condition SQLExpression) Join {
	i.condition = condition
	return i
}

func (i Join) Limit(limit int) Join {
	i.limit = limit
	return i
}

func (i Join) LeftOuterJoin(q Unwrappable) (m MultiJoin) {
	m.sets = append(m.sets, i.leftSet, i.rightSet, q.Unwrap())
	m.joinTypes = append(m.joinTypes, i.joinType, LeftOuterJoinType)
	m.conditions = append(m.conditions, i.condition)
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
	return JoinResultIterator{queryError: err, leftSet: i.leftSet, rightSet: i.rightSet, rows: rows}, nil
}

type JoinResultIterator struct {
	queryError error
	leftSet    QuerySet
	rightSet   QuerySet
	rows       pgx.Rows
}

func (i *JoinResultIterator) HasNext() bool {
	if i.queryError != nil {
		return false
	}
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
	sets       []QuerySet
	joinTypes  []JoinType
	conditions []SQLExpression
}

func (m MultiJoin) On(condition SQLExpression) MultiJoin {
	m.conditions = append(m.conditions, condition)
	return m
}

func (m MultiJoin) LeftOuterJoin(q Unwrappable) MultiJoin {
	m.sets = append(m.sets, q.Unwrap())
	m.joinTypes = append(m.joinTypes, LeftOuterJoinType)
	return m
}

func (m MultiJoin) SQLOn(w io.Writer) {
	fmt.Fprint(w, "SELECT ")
	for i, each := range m.sets {
		if i > 0 {
			fmt.Fprint(w, ",")
		}
		writeAccessOn(each.selectors, w)
	}
	fmt.Fprint(w, " FROM ")
	first := m.sets[0]
	first.fromSectionOn(w)
	// collect all conditions from all sets
	wheres := []SQLExpression{}
	for _, each := range m.sets {
		if each.condition != EmptyCondition {
			wheres = append(wheres, each.condition)
		}
	}
	if len(wheres) > 0 {
		fmt.Fprint(w, " WHERE ")
		for i, each := range wheres {
			if i > 0 {
				fmt.Fprint(w, " AND ")
			}
			each.SQLOn(w)
		}
	}
	for j := 0; j < len(m.joinTypes); j++ {
		jt := m.joinTypes[j]
		writeJoinType(jt, w)
		set := m.sets[j+1]
		set.fromSectionOn(w)
		fmt.Fprint(w, " ON ")
		m.conditions[j].SQLOn(w)
	}
}
