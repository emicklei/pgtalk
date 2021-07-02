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

func (p ValuePrinter) SQLOn(b io.Writer) { fmt.Fprintf(b, "%v", p.v) }

type ValuesPrinter struct {
	vs []interface{}
}

func (p ValuesPrinter) SQLOn(b io.Writer) {
	fmt.Fprintf(b, "(")
	for i, each := range p.vs {
		if i > 0 {
			fmt.Fprintf(b, ",")
		}
		if s, ok := each.(string); ok {
			fmt.Fprintf(b, "'%s'", s)
		}
	}
	fmt.Fprintf(b, ")")
}

type ScanToWrite struct {
	RW     ColumnAccessor
	Entity interface{}
}

func (s ScanToWrite) Scan(fieldValue interface{}) error {
	s.RW.WriteInto(s.Entity, fieldValue)
	return nil
}

type LiteralString string

func (l LiteralString) SQLOn(b io.Writer) {
	io.WriteString(b, "'")
	io.WriteString(b, string(l))
	io.WriteString(b, "'")
}

type NoCondition struct{}

var EmptyCondition = NoCondition{}

func (n NoCondition) SQLOn(b io.Writer) {}

type columnInfo struct {
	tableInfo  TableInfo
	columnName string
}

func (c columnInfo) Name() string { return c.columnName }

func (c columnInfo) SQLOn(w io.Writer) {
	fmt.Fprintf(w, "%s.%s", c.tableInfo.Alias, c.columnName)
}

func SQL(w SQLWriter) string {
	b := new(bytes.Buffer)
	w.SQLOn(b)
	return b.String()
}
