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

func (t TableInfo) Equals(o TableInfo) bool {
	return t.Name == o.Name && t.Schema == o.Schema && t.Alias == o.Alias
}

func (t TableInfo) String() string {
	return fmt.Sprintf("table(%s.%s %s)", t.Schema, t.Name, t.Alias)
}

type TableAccessor struct {
	TableInfo
	Factory    NewEntityFunc
	AllColumns []ColumnAccessor
}

var EmptyColumnAccessor = []ColumnAccessor{}

type ValuePrinter struct {
	v interface{}
}

func MakeValuePrinter(v interface{}) ValuePrinter { return ValuePrinter{v: v} }

func (p ValuePrinter) SQLOn(b io.Writer) { fmt.Fprintf(b, "%v", p.v) }

// Collect is part of SQLExpression
func (p ValuePrinter) Collect(list []ColumnAccessor) []ColumnAccessor {
	return list
}

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

func (p ValuesPrinter) Collect(list []ColumnAccessor) []ColumnAccessor {
	return list
}

type scanToWrite struct {
	access ColumnAccessor
	entity interface{}
}

func (s scanToWrite) Scan(fieldValue interface{}) error {
	s.access.SetFieldValue(s.entity, fieldValue)
	return nil
}

type LiteralString string

func (l LiteralString) SQLOn(b io.Writer) {
	io.WriteString(b, "'")
	io.WriteString(b, string(l))
	io.WriteString(b, "'")
}

func (l LiteralString) Collect(list []ColumnAccessor) []ColumnAccessor { return list }

type NoCondition struct{}

var EmptyCondition = NoCondition{}

func (n NoCondition) SQLOn(b io.Writer)                              {}
func (n NoCondition) Collect(list []ColumnAccessor) []ColumnAccessor { return list }

const (
	IsPrimary  = true
	NotPrimary = false
	NotNull    = true
	Nullable   = false
)

type ColumnInfo struct {
	tableInfo            TableInfo
	columnName           string
	isPrimary            bool
	notNull              bool
	isMixedCase          bool
	tableAttributeNumber uint16
}

func MakeColumnInfo(t TableInfo, name string, isPrimary bool, isNotNull bool, tableAttributeNumber uint16) ColumnInfo {
	return ColumnInfo{
		tableInfo:            t,
		columnName:           name,
		notNull:              isNotNull,
		isPrimary:            isPrimary,
		isMixedCase:          strings.ToLower(name) != name,
		tableAttributeNumber: tableAttributeNumber,
	}
}

func (c ColumnInfo) String() string {
	return fmt.Sprintf("column(%s.%s:%s)", c.tableInfo.Schema, c.tableInfo.Name, c.columnName)
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
func PrettySQL(sql SQLWriter) string {
	b := new(bytes.Buffer)
	wasUpper := false
	for i, each := range strings.Fields(SQL(sql)) {
		if i > 0 { // skip first
			if len(each) > 1 { // sql token are multi-char
				if !strings.HasPrefix(each, "'") {
					// no break after IS,NOT,NULL
					if strings.ToUpper(each) == each && strings.Index("IS NOT NULL", each) == -1 {
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
