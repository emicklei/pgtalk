package pgtalk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"

	"github.com/jackc/pgtype"
)

type TableInfo struct {
	Name    string
	Schema  string
	Alias   string
	Columns []ColumnAccessor
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

type TableAccessor[T any] struct {
	TableInfo
	AllColumns []ColumnAccessor
}

var EmptyColumnAccessor = []ColumnAccessor{}

type valuePrinter struct {
	v interface{}
}

func MakeValuePrinter(v interface{}) valuePrinter { return valuePrinter{v: v} }

func (p valuePrinter) SQLOn(b io.Writer) {
	if e, ok := p.v.(pgtype.UUID); ok {
		fmt.Fprintf(b, "'%s'::uuid", encodeUUID(e.Bytes))
		return
	}
	if e, ok := p.v.(pgtype.Date); ok {
		fmt.Fprintf(b, "'%s'::date", toJSON(e))
		return
	}
	fmt.Fprintf(b, "%v", p.v)
}

// hack
func toJSON(m json.Marshaler) string {
	data, _ := m.MarshalJSON()
	return strings.Trim(string(data), "\"")
}

// encodeUUID converts a uuid byte array to UUID standard string form.
func encodeUUID(src [16]byte) string {
	return fmt.Sprintf("%x-%x-%x-%x-%x", src[0:4], src[4:6], src[6:8], src[8:10], src[10:16])
}

// Collect is part of SQLExpression
func (p valuePrinter) Collect(list []ColumnAccessor) []ColumnAccessor {
	return list
}

type valuesPrinter struct {
	vs []interface{}
}

func (p valuesPrinter) SQLOn(b io.Writer) {
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

func (p valuesPrinter) Collect(list []ColumnAccessor) []ColumnAccessor {
	return list
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

// SQL returns the full SQL string for any SQLWriter implementation.
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

const HideNilValues = true

func StringWithFields(v interface{}, includePresent bool) string {
	vt := reflect.TypeOf(v)
	if vt.Kind() == reflect.Ptr {
		vt = vt.Elem()
	}
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	b := new(bytes.Buffer)
	fmt.Fprint(b, vt.PkgPath())
	fmt.Fprint(b, ".")
	fmt.Fprint(b, vt.Name())
	fmt.Fprint(b, "{")
	for i := 0; i < vt.NumField(); i++ {
		f := vt.Field(i)
		fv := rv.Field(i)
		if fv.IsZero() {
			continue
		}
		var fi interface{}
		// check fields that have pointer type
		if fv.Kind() == reflect.Pointer {
			fi = fv.Elem().Interface()
		} else {
			fi = fv.Interface()
		}
		fmt.Fprintf(b, "%s:%v ", f.Name, fi)
	}
	fmt.Fprint(b, "}")
	return b.String()
}
