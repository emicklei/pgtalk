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

func (t TableInfo) SQLOn(w io.Writer) {
	fmt.Fprintf(w, "%s.%s %s", t.Schema, t.Name, t.Alias)
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

type ColumnInfo struct {
	tableInfo   TableInfo
	columnName  string
	isPrimary   bool
	notNull     bool
	isMixedCase bool
}

func MakeColumnInfo(t TableInfo, name string, isPrimary bool, isNotNull bool) ColumnInfo {
	return ColumnInfo{
		tableInfo:   t,
		columnName:  name,
		notNull:     isNotNull,
		isPrimary:   isPrimary,
		isMixedCase: strings.ToLower(name) != name,
	}
}

func (c ColumnInfo) Name() string {
	if c.isMixedCase {
		return strconv.Quote(c.columnName)
	}
	return c.columnName
}

func (c ColumnInfo) SQLOn(w io.Writer) {
	fmt.Fprintf(w, "%s.%s", c.tableInfo.Alias, c.Name())
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

// PrettySQL returns a multiline SQL statement with line breaks before each next uppercase token
func PrettySQL(sql string) string {
	b := new(bytes.Buffer)
	wasUpper := false
	for i, each := range strings.Fields(sql) {
		if i > 0 { // skip first
			if len(each) > 1 { // sql token are multi-char
				if !strings.HasPrefix(each, "'") {
					if strings.ToUpper(each) == each {
						if !wasUpper {
							io.WriteString(b, "\n")
							wasUpper = true
						}
					} else {
						wasUpper = false
					}
				}
			} else {
				wasUpper = false
			}
		}
		fmt.Fprintf(b, "%s ", each)
	}
	return b.String()
}
