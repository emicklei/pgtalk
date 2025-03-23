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
	condition    SQLExpression
	joinType     joinType
	limit        int
	offset       int
}

func (i join) SQLOn(w WriteContext) {
	fmt.Fprint(w, "SELECT\n")
	if i.leftSet == nil {
		fmt.Fprint(w, "ERROR: no left queryset set")
		return
	}
	if i.rightSet == nil {
		fmt.Fprint(w, "ERROR: no right queryset set")
		return
	}
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
	if i.condition == nil {
		fmt.Fprint(w, "ERROR: no condition set")
		return
	}
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

func (i join) Join(q querySet) (m multiJoin) {
	m.sets = append(m.sets, i.leftSet, i.rightSet, q)
	m.joinTypes = append(m.joinTypes, i.joinType, innerJoinType)
	m.conditions = append(m.conditions, i.condition)
	return
}

func (i join) RightJoin(q querySet) (m multiJoin) {
	m.sets = append(m.sets, i.leftSet, i.rightSet, q)
	m.joinTypes = append(m.joinTypes, i.joinType, rightOuterJoinType)
	m.conditions = append(m.conditions, i.condition)
	return
}

func (i join) LeftOuterJoin(q querySet) (m multiJoin) {
	m.sets = append(m.sets, i.leftSet, i.rightSet, q)
	m.joinTypes = append(m.joinTypes, i.joinType, leftOuterJoinType)
	m.conditions = append(m.conditions, i.condition)
	return
}

func (i join) FullOuterJoin(q querySet) (m multiJoin) {
	m.sets = append(m.sets, i.leftSet, i.rightSet, q)
	m.joinTypes = append(m.joinTypes, i.joinType, fullOuterJoinType)
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
	querySets := []querySet{i.leftSet, i.rightSet}
	models := []any{left, right}

	// TODO remove duplicate code

	// we cannot scan with all fields for all models because
	// in a JOIN values for fields per model may all be NULL
	// so we collect all raw values and inspect them before
	// scanning them into the fields per model
	raw := i.rows.RawValues()
	// where in the raw values are we taking values per model
	offset := 0
	// complete type map
	typeMap := i.rows.Conn().TypeMap()
	// all fields descriptions
	fieldDefs := i.rows.FieldDescriptions()
	for m, eachSet := range querySets {
		// each set has its own set of accessors
		modelAccessors := eachSet.selectAccessors()
		// take a slice of values and fields
		subvalues := raw[offset : offset+len(modelAccessors)]
		subFields := fieldDefs[offset : offset+len(modelAccessors)]
		offset += len(modelAccessors)
		// check if there is data to scan
		if hasNullsOnly(subvalues) {
			continue
		}
		// at least one non-zero value is available
		dest := []any{}
		// collect the destinations
		for _, each := range modelAccessors {
			dest = append(dest, each.FieldValueToScan(models[m]))
		}
		// scan the subvalues into all destinations.
		err := pgx.ScanRow(typeMap, subFields, subvalues, dest...)
		if err != nil {
			// abort on the first error
			return err
		}
	}
	return nil
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

func (m multiJoin) Join(q querySet) multiJoin {
	m.sets = append(m.sets, q)
	m.joinTypes = append(m.joinTypes, innerJoinType)
	return m
}

func (m multiJoin) RightOuterJoin(q querySet) multiJoin {
	m.sets = append(m.sets, q)
	m.joinTypes = append(m.joinTypes, rightOuterJoinType)
	return m
}

func (m multiJoin) FullOuterJoin(q querySet) multiJoin {
	m.sets = append(m.sets, q)
	m.joinTypes = append(m.joinTypes, fullOuterJoinType)
	return m
}

func (m multiJoin) Named(preparedName string) multiJoin {
	m.preparedName = preparedName
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
	// we cannot scan with all fields for all models because
	// in a JOIN values for fields per model may all be NULL
	// so we collect all raw values and inspect them before
	// scanning them into the fields per model
	raw := i.rows.RawValues()
	// where in the raw values are we taking values per model
	offset := 0
	// complete type map
	typeMap := i.rows.Conn().TypeMap()
	// all fields descriptions
	fieldDefs := i.rows.FieldDescriptions()
	for m, eachSet := range i.querySets {
		// each set has its own set of accessors
		subAccessors := eachSet.selectAccessors()
		// take a slice of values and fields
		subvalues := raw[offset : offset+len(subAccessors)]
		subFields := fieldDefs[offset : offset+len(subAccessors)]
		offset += len(subAccessors)
		// check if there is data to scan
		if hasNullsOnly(subvalues) {
			continue
		}
		// at least one non-zero value is available
		dest := []any{}
		// collect the destinations
		for _, each := range subAccessors {
			dest = append(dest, each.FieldValueToScan(models[m]))
		}
		// scan the subvalues into all destinations.
		if err := pgx.ScanRow(typeMap, subFields, subvalues, dest...); err != nil {
			// abort on the first error
			return err
		}
	}
	return nil
}

func hasNullsOnly(sub [][]byte) bool {
	for _, each := range sub {
		if len(each) != 0 {
			return false
		}
	}
	return true
}
