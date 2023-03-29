package pgtalk

import (
	"context"
	"fmt"
	"io"

	"github.com/jackc/pgx/v5"
)

type joinType int

const (
	innerJoinType joinType = iota
	leftOuterJoinType
	rightOuterJoinType
	fullOuterJoinType
)

type join struct {
	preparedName string
	leftSet      querySet
	rightSet     querySet
	onLeft       ColumnAccessor
	onRight      ColumnAccessor
	condition    SQLExpression
	joinType     joinType
	limit        int
	offset       int
}

func (i join) SQLOn(w WriteContext) {
	fmt.Fprint(w, "SELECT\n")
	left := i.leftSet.selectAccessors()
	wl := i.leftSet.augmentedContext(w)
	wr := i.rightSet.augmentedContext(w)
	writeAccessOn(left, wl)
	right := i.rightSet.selectAccessors()
	if len(right) > 0 {
		fmt.Fprint(wl, ",")
	}
	writeAccessOn(right, wr)
	fmt.Fprint(w, "\nFROM ")
	i.leftSet.fromSectionOn(wl)
	writeJoinType(i.joinType, w)
	i.rightSet.fromSectionOn(wr)
	fmt.Fprint(w, "\nON ")
	i.condition.SQLOn(w) // TODO which tableInfo to use?
	if _, ok := i.leftSet.whereCondition().(noCondition); !ok {
		fmt.Fprint(wl, "\nWHERE ")
		i.leftSet.whereCondition().SQLOn(wl)
	}
	if i.limit > 0 {
		fmt.Fprintf(wl, "\nLIMIT %d", i.limit)
	}
	if i.offset > 0 {
		fmt.Fprintf(wl, "\nOFFSET %d", i.offset)
	}
	// TODO RightSet where
}

func writeJoinType(t joinType, w io.Writer) {
	switch t {
	case innerJoinType:
		fmt.Fprint(w, "\nINNER JOIN ")
	case leftOuterJoinType:
		fmt.Fprint(w, "\nLEFT OUTER JOIN ")
	case rightOuterJoinType:
		fmt.Fprint(w, "\nRIGHT OUTER JOIN ")
	case fullOuterJoinType:
		fmt.Fprint(w, "\nFULL OUTER JOIN ")
	}
}

func (i join) Named(preparedName string) join {
	i.preparedName = preparedName
	return i
}

func (i join) On(condition SQLExpression) join {
	i.condition = condition
	return i
}

func (i join) Limit(limit int) join {
	i.limit = limit
	return i
}

func (i join) Offset(offset int) join {
	i.offset = offset
	return i
}

func (i join) LeftOuterJoin(q querySet) (m multiJoin) {
	m.sets = append(m.sets, i.leftSet, i.rightSet, q)
	m.joinTypes = append(m.joinTypes, i.joinType, leftOuterJoinType)
	m.conditions = append(m.conditions, i.condition)
	return
}

func (i join) Exec(ctx context.Context, conn querier, parameters ...*QueryParameter) (it joinResultIterator, err error) {
	params := argumentValues(parameters)
	sql := SQL(i)
	if i.preparedName != "" {
		if p, ok := conn.(preparer); ok {
			_, err := p.Prepare(ctx, i.preparedName, sql)
			if err != nil {
				return joinResultIterator{queryError: err}, err
			}
		}
	}
	rows, err := conn.Query(ctx, sql, params...)
	return joinResultIterator{queryError: err, leftSet: i.leftSet, rightSet: i.rightSet, rows: rows}, nil
}

type joinResultIterator struct {
	queryError error
	leftSet    querySet
	rightSet   querySet
	rows       pgx.Rows
}

func (i *joinResultIterator) HasNext() bool {
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

func (i *joinResultIterator) Err() error {
	if i.queryError != nil {
		return i.queryError
	}
	return i.rows.Err()
}

func (i *joinResultIterator) Next(left any, right any) error {
	if i.queryError != nil {
		return i.queryError
	}
	sw := []any{}
	// left
	for _, each := range i.leftSet.selectAccessors() {
		sw = append(sw, each.FieldValueToScan(left))
	}
	// right
	for _, each := range i.rightSet.selectAccessors() {
		sw = append(sw, each.FieldValueToScan(right))
	}
	return i.rows.Scan(sw...)
}

type multiJoin struct {
	preparedName string
	sets         []querySet
	joinTypes    []joinType
	conditions   []SQLExpression
}

func (m multiJoin) On(condition SQLExpression) multiJoin {
	m.conditions = append(m.conditions, condition)
	return m
}

func (m multiJoin) LeftOuterJoin(q querySet) multiJoin {
	m.sets = append(m.sets, q)
	m.joinTypes = append(m.joinTypes, leftOuterJoinType)
	return m
}

func (m multiJoin) Exec(ctx context.Context, conn querier) (*multiJoinResultIterator, error) {
	sql := SQL(m)
	if m.preparedName != "" {
		if p, ok := conn.(preparer); ok {
			_, err := p.Prepare(ctx, m.preparedName, sql)
			if err != nil {
				return &multiJoinResultIterator{queryError: err}, nil
			}
		}
	}
	rows, err := conn.Query(ctx, sql)
	return &multiJoinResultIterator{queryError: err, querySets: m.sets, rows: rows}, nil
}

type tableWhere struct {
	expression SQLExpression
	tableInfo  TableInfo
}

func (m multiJoin) SQLOn(w WriteContext) {
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
	wheres := []SQLWriter{}
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

type multiJoinResultIterator struct {
	queryError error
	querySets  []querySet
	rows       pgx.Rows
}

func (i *multiJoinResultIterator) Err() error {
	if i.queryError != nil {
		return i.queryError
	}
	return i.rows.Err()
}

func (i *multiJoinResultIterator) HasNext() bool {
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

func (i *multiJoinResultIterator) Next(models ...any) error {
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
	sw := []any{}
	// all sets
	for m, eachSet := range i.querySets {
		for _, each := range eachSet.selectAccessors() {
			sw = append(sw, each.FieldValueToScan(models[m]))
		}
	}
	return i.rows.Scan(sw...)
}
