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
	leftSet      querySet
	rightSet     querySet
	onLeft       ColumnAccessor
	onRight      ColumnAccessor
	condition    SQLExpression
	joinType     JoinType
	limit        int
}

func (i Join) SQLOn(w WriteContext) {
	fmt.Fprint(w, "SELECT ")
	left := i.leftSet.selectAccessors()
	wl := i.leftSet.augmentedContext(w)
	wr := i.rightSet.augmentedContext(w)
	writeAccessOn(left, wl)
	if len(left) > 0 {
		fmt.Fprint(wl, ",")
	}
	writeAccessOn(i.rightSet.selectAccessors(), wr)
	fmt.Fprint(w, " FROM ")
	i.leftSet.fromSectionOn(wl)
	writeJoinType(i.joinType, w)
	i.rightSet.fromSectionOn(wr)
	fmt.Fprint(w, " ON ")
	i.condition.SQLOn(w) // TODO which tableInfo to use?
	if _, ok := i.leftSet.whereCondition().(NoCondition); !ok {
		fmt.Fprint(wl, " WHERE ")
		i.leftSet.whereCondition().SQLOn(wl)
	}
	if i.limit > 0 {
		fmt.Fprintf(wl, " LIMIT %d", i.limit)
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

func (i Join) LeftOuterJoin(q querySet) (m MultiJoin) {
	m.sets = append(m.sets, i.leftSet, i.rightSet, q)
	m.joinTypes = append(m.joinTypes, i.joinType, LeftOuterJoinType)
	m.conditions = append(m.conditions, i.condition)
	return
}

func (i Join) Exec(ctx context.Context, conn *pgx.Conn) (it JoinResultIterator, err error) {
	sql := SQL(i)
	if i.preparedName != "" {
		_, err := conn.Prepare(ctx, i.preparedName, sql)
		if err != nil {
			return JoinResultIterator{queryError: err}, err
		}
	}
	rows, err := conn.Query(ctx, sql)
	return JoinResultIterator{queryError: err, leftSet: i.leftSet, rightSet: i.rightSet, rows: rows}, nil
}

type JoinResultIterator struct {
	queryError error
	leftSet    querySet
	rightSet   querySet
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
	if i.queryError != nil {
		return i.queryError
	}
	sw := []interface{}{}
	// left
	for _, each := range i.leftSet.selectAccessors() {
		sw = append(sw, each.FieldToScan(left))
	}
	// right
	for _, each := range i.rightSet.selectAccessors() {
		sw = append(sw, each.FieldToScan(right))
	}
	return i.rows.Scan(sw...)
}

type MultiJoin struct {
	preparedName string
	sets         []querySet
	joinTypes    []JoinType
	conditions   []SQLExpression
}

func (m MultiJoin) On(condition SQLExpression) MultiJoin {
	m.conditions = append(m.conditions, condition)
	return m
}

func (m MultiJoin) LeftOuterJoin(q querySet) MultiJoin {
	m.sets = append(m.sets, q)
	m.joinTypes = append(m.joinTypes, LeftOuterJoinType)
	return m
}

func (m MultiJoin) Exec(ctx context.Context, conn *pgx.Conn) (*MultiJoinResultIterator, error) {
	sql := SQL(m)
	if m.preparedName != "" {
		_, err := conn.Prepare(ctx, m.preparedName, sql)
		if err != nil {
			return &MultiJoinResultIterator{queryError: err}, nil
		}
	}
	rows, err := conn.Query(ctx, sql)
	return &MultiJoinResultIterator{queryError: err, querySets: m.sets, rows: rows}, nil
}

type tableWhere struct {
	expression SQLExpression
	tableInfo  TableInfo
}

func (m MultiJoin) SQLOn(w WriteContext) {
	fmt.Fprint(w, "SELECT ")
	for i, each := range m.sets {
		if i > 0 && len(each.selectAccessors()) > 0 {
			fmt.Fprint(w, ",")
		}
		writeAccessOn(each.selectAccessors(), w)
	}
	fmt.Fprint(w, " FROM ")
	first := m.sets[0]
	first.fromSectionOn(w)
	// collect all conditions from all sets
	wheres := []SQLExpression{}
	for _, each := range m.sets {
		if each.whereCondition() != EmptyCondition {
			wheres = append(wheres, each.whereCondition())
		}
	}
	for j := 0; j < len(m.joinTypes); j++ {
		jt := m.joinTypes[j]
		writeJoinType(jt, w)
		set := m.sets[j+1]
		set.fromSectionOn(w)
		if j < len(m.conditions) {
			fmt.Fprint(w, " ON ")
			m.conditions[j].SQLOn(w)
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
}

type MultiJoinResultIterator struct {
	queryError error
	querySets  []querySet
	rows       pgx.Rows
}

func (i *MultiJoinResultIterator) Err() error {
	if i.queryError != nil {
		return i.queryError
	}
	return i.rows.Err()
}

func (i *MultiJoinResultIterator) HasNext() bool {
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

func (i *MultiJoinResultIterator) Next(models ...interface{}) error {
	if i.queryError != nil {
		return i.queryError
	}
	// count non-empty querysets
	countNonEmpty := 0
	for _, each := range i.querySets {
		if len(each.selectAccessors()) != 0 {
			countNonEmpty++
		}
	}
	// check models count matches
	if mc, qc := len(models), countNonEmpty; mc != qc {
		return fmt.Errorf("number of models [%d] does not match select count [%d]", mc, qc)
	}
	// TODO how to check model types?
	sw := []interface{}{}
	// all sets
	for m, eachSet := range i.querySets {
		for _, each := range eachSet.selectAccessors() {
			sw = append(sw, each.FieldToScan(models[m]))
		}
	}
	return i.rows.Scan(sw...)
}
