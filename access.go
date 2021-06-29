package pgtalk

import (
	"bytes"
	"fmt"
	"io"
)

type TableInfo struct {
	Name  string
	Alias string
}

var EmptyColumnAccessor = []ColumnAccessor{}

type ValuePrinter struct {
	v interface{}
}

func (p ValuePrinter) SQL() string { return fmt.Sprintf("%v", p.v) }

type ScanToWrite struct {
	RW     ColumnAccessor
	Entity interface{}
}

func (s ScanToWrite) Scan(fieldValue interface{}) error {
	s.RW.WriteInto(s.Entity, fieldValue)
	return nil
}

type LiteralString string

func (l LiteralString) SQL() string {
	b := new(bytes.Buffer)
	io.WriteString(b, "'")
	io.WriteString(b, string(l))
	io.WriteString(b, "'")
	return b.String()
}

type NoCondition struct{}

var EmptyCondition = NoCondition{}

func (n NoCondition) SQL() string { return "" }

type columnInfo struct {
	tableInfo  TableInfo
	columnName string
}

func (c columnInfo) Name() string { return c.columnName }

func (c columnInfo) SQL() string {
	return fmt.Sprintf("%s.%s", c.tableInfo.Alias, c.columnName)
}
