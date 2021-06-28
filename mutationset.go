package pgtalk

import (
	"bytes"
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

func (m MutationSet) WhereSection() string {
	return m.condition.SQL()
}

// SQL returns the full SQL mutation query
func (m MutationSet) SQL() string {
	// TODO

	if m.operationType == MutationInsert {
		return "INSERT INTO " + m.tableInfo.Name + " (" + m.ColumnsSection() + ") values (" + m.ValuesSection() + ")"
	}
	if m.operationType == MutationDelete {
		return "DELETE FROM " + m.tableInfo.Name + " WHERE " + m.WhereSection()
	}
	if m.operationType == MutationUpdate {
		return "UPDATE " + m.tableInfo.Name + " SET " + m.WhereSection() + " WHERE " + m.WhereSection()
	}
	return "-- unknown operation type"
}

func (m MutationSet) Where(condition SQLWriter) MutationSet {
	m.condition = condition
	return m
}

// todo
func (m MutationSet) On() MutationSet {
	return m
}

func (m MutationSet) ColumnsSection() string {
	buf := new(bytes.Buffer)
	for i, each := range m.selectors {
		if i > 0 {
			io.WriteString(buf, ",")
		}
		io.WriteString(buf, each.Name())
	}
	return buf.String()
}

func (m MutationSet) ValuesSection() string {
	buf := new(bytes.Buffer)
	for i := range m.selectors {
		if i > 0 {
			io.WriteString(buf, ",")
		}
		fmt.Fprintf(buf, "$%d", i+1)
	}
	return buf.String()
}

func (m MutationSet) SetSection() string {
	buf := new(bytes.Buffer)
	for i, each := range m.selectors {
		if i > 0 {
			io.WriteString(buf, ",")
		}
		fmt.Fprintf(buf, "%s = %s", each.Name(), each.ValueAsSQL())
	}
	return buf.String()
}
