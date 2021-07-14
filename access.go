package pgtalk

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type TableInfo struct {
	Name   string
	Schema string
	Alias  string
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

type scanToWrite struct {
	access ColumnAccessor
	entity interface{}
}

func (s scanToWrite) Scan(fieldValue interface{}) error {
	s.access.WriteInto(s.entity, fieldValue)
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
	tableInfo   TableInfo
	columnName  string
	notNull     bool // TODO validate in Go or by Postgres?
	isMixedCase bool
}

func makeColumnInfo(t TableInfo, c string) columnInfo {
	return columnInfo{
		tableInfo:   t,
		columnName:  c,
		isMixedCase: strings.ToLower(c) != c,
	}
}

func (c columnInfo) Name() string {
	if c.isMixedCase {
		return strconv.Quote(c.columnName)
	}
	return c.columnName
}

func (c columnInfo) SQLOn(w io.Writer) {
	fmt.Fprintf(w, "%s.%s", c.tableInfo.Alias, c.columnName)
}

func SQL(w SQLWriter) string {
	b := new(bytes.Buffer)
	w.SQLOn(b)
	return b.String()
}

func writeAccessOn(list []ColumnAccessor, w io.Writer) {
	for i, each := range list {
		if i > 0 {
			io.WriteString(w, ",")
		}
		each.SQLOn(w)
	}
}
