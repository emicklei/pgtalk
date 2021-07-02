package pgtalk

import (
	"fmt"
	"io"
)

const (
	MutationDelete = iota
	MutationInsert
	MutationUpdate
)

type MutationSet struct {
	tableInfo     TableInfo
	selectors     []ColumnAccessor
	condition     SQLWriter
	operationType int
}

func MakeMutationSet(tableInfo TableInfo, selectors []ColumnAccessor, operationType int) MutationSet {
	return MutationSet{
		tableInfo:     tableInfo,
		selectors:     selectors,
		condition:     EmptyCondition,
		operationType: operationType}
}

// SQL returns the full SQL mutation query
func (m MutationSet) SQLOn(w io.Writer) {
	if m.operationType == MutationInsert {
		fmt.Fprint(w, "INSERT INTO ")
		fmt.Fprint(w, m.tableInfo.Name)
		fmt.Fprint(w, " (")
		//m.ColumnsSection() + ") values (" + m.ValuesSection() + ")"
		return
	}
	if m.operationType == MutationDelete {
		fmt.Fprint(w, "INSERT INTO ")
		fmt.Fprint(w, m.tableInfo.Name)
		//return "DELETE FROM " + m.tableInfo.Name + " WHERE " + m.WhereSection()
		return
	}
	if m.operationType == MutationUpdate {
		fmt.Fprint(w, "INSERT INTO ")
		fmt.Fprint(w, m.tableInfo.Name)
		//return "UPDATE " + m.tableInfo.Name + " SET " + m.WhereSection() + " WHERE " + m.WhereSection()
		return
	}
}

func (m MutationSet) Where(condition SQLWriter) MutationSet {
	m.condition = condition
	return m
}

// todo
func (m MutationSet) On() MutationSet {
	return m
}

func (m MutationSet) columnsSectionOn(buf io.Writer) {
	for i, each := range m.selectors {
		if i > 0 {
			io.WriteString(buf, ",")
		}
		io.WriteString(buf, each.Name())
	}
}

func (m MutationSet) valuesSectionOn(buf io.Writer) {
	for i := range m.selectors {
		if i > 0 {
			io.WriteString(buf, ",")
		}
		fmt.Fprintf(buf, "$%d", i+1)
	}
}

func (m MutationSet) setSectionOn(w io.Writer) {
	for i, each := range m.selectors {
		if i > 0 {
			io.WriteString(w, ",")
		}
		fmt.Fprintf(w, "%s = ", each.Name())
		each.ValueAsSQLOn(w)
	}
}
