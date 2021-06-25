package xs

import (
	"bytes"
	"fmt"
	"io"
)

type MutationSet struct {
	tableInfo TableInfo
	selectors []ReadWrite
}

func MakeMutationSet(tableInfo TableInfo, selectors []ReadWrite) MutationSet {
	return MutationSet{
		tableInfo: tableInfo,
		selectors: selectors}
}

// SQL returns the full SQL mutation query
func (m MutationSet) SQL() string {
	// temp
	return "INSERT INTO " + m.tableInfo.Name + " (" + m.ColumnsSection() + ") values (" + m.ValuesSection() + ")"
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
