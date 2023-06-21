package pgtalk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/emicklei/pgtalk/convert"
	"github.com/jackc/pgx/v5/pgtype"
)

var EmptyColumnAccessor = []ColumnAccessor{}

type valuePrinter struct {
	unimplementedBooleanExpression
	v any
}

func makeValuePrinter(v any) valuePrinter { return valuePrinter{v: v} }

func (p valuePrinter) SQLOn(w WriteContext) {
	if e, ok := p.v.(SQLWriter); ok {
		e.SQLOn(w)
		return
	}
	// TODO rewrite using type switch
	if e, ok := p.v.(string); ok {
		fmt.Fprintf(w, "'%s'", e)
		return
	}
	if e, ok := p.v.(pgtype.UUID); ok {
		fmt.Fprintf(w, "'%s'::uuid", encodeUUID(e.Bytes))
		return
	}
	if e, ok := p.v.(pgtype.Date); ok {
		fmt.Fprintf(w, "'%s'::date", toJSON(e))
		return
	}
	if e, ok := p.v.(pgtype.Text); ok {
		fmt.Fprintf(w, "'%s'", e.String)
		return
	}
	if e, ok := p.v.(pgtype.Timestamptz); ok {
		// '2015-07-16 00:00:00'
		// RFC3339     = "2006-01-02T15:04:05Z07:00"
		fmt.Fprintf(w, "timestamp '%s'", e.Time.Format("2006-01-02 15:04:05"))
		return
	}
	if e, ok := p.v.(pgtype.Timestamp); ok {
		fmt.Fprintf(w, "timestamp '%s'", e.Time.Format("2006-01-02 15:04:05"))
		return
	}
	fmt.Fprintf(w, "%v", p.v)
}

type unimplementedBooleanExpression struct{}

func (unimplementedBooleanExpression) And(e SQLExpression) SQLExpression { panic("unsupported And") }
func (unimplementedBooleanExpression) Or(e SQLExpression) SQLExpression  { panic("unsupported Or") }

// hack
func toJSON(m json.Marshaler) string {
	data, _ := m.MarshalJSON()
	return strings.Trim(string(data), "\"")
}

// encodeUUID converts a uuid byte array to UUID standard string form.
func encodeUUID(src [16]byte) string {
	return fmt.Sprintf("%x-%x-%x-%x-%x", src[0:4], src[4:6], src[6:8], src[8:10], src[10:16])
}

type valuesPrinter struct {
	unimplementedBooleanExpression
	vs []any
}

func (p valuesPrinter) SQLOn(w WriteContext) {
	fmt.Fprintf(w, "(")
	for i, each := range p.vs {
		if i > 0 {
			fmt.Fprintf(w, ",")
		}
		valuePrinter{v: each}.SQLOn(w)
	}
	fmt.Fprintf(w, ")")
}

type literalString struct {
	unimplementedBooleanExpression
	value string
}

func newLiteralString(s string) literalString {
	return literalString{value: s}
}

func (l literalString) SQLOn(w WriteContext) {
	io.WriteString(w, "'")
	io.WriteString(w, l.value)
	io.WriteString(w, "'")
}

type noCondition struct{}

var EmptyCondition SQLExpression = noCondition{}

func (n noCondition) SQLOn(w WriteContext) {}

// And returns the argument as the receiver is a no operation
func (n noCondition) And(ex SQLExpression) SQLExpression {
	return ex
}

// And returns the argument as the receiver is a no operation
func (n noCondition) Or(ex SQLExpression) SQLExpression {
	return ex
}

const (
	IsPrimary  = true
	NotPrimary = false
	NotNull    = true
	Nullable   = false
)

func writeAccessOn(list []ColumnAccessor, w WriteContext) {
	for i, each := range list {
		if i > 0 {
			io.WriteString(w, ",\n")
		}
		io.WriteString(w, "\t")
		each.SQLOn(w)
	}
}

func writeListOn(list []SQLWriter, w WriteContext) {
	for i, each := range list {
		if i > 0 {
			io.WriteString(w, ",\n")
		}
		io.WriteString(w, "\t")
		each.SQLOn(w)
	}
}

const HideNilValues = true

func StringWithFields(v any, includePresent bool) string {
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
		if !f.IsExported() {
			continue
		}
		fv := rv.Field(i)
		if fv.IsZero() {
			continue
		}
		var fi any
		// check fields that have pointer type
		if fv.Kind() == reflect.Pointer {
			fi = fv.Elem().Interface()
		} else {
			fi = fv.Interface()
		}
		// StringWithFields is for debugging purposes so we try to display the xxx-xxx-xxx representation of a UUID
		if uid, ok := fi.(pgtype.UUID); ok {
			fmt.Fprintf(b, "%s:%s ", f.Name, convert.UUIDToString(uid))
		} else {
			fmt.Fprintf(b, "%s:%v ", f.Name, fi)
		}
	}
	fmt.Fprint(b, "}")
	return b.String()
}

type SQLLiteral struct {
	Literal string
}

func (l SQLLiteral) SQLOn(w WriteContext) {
	w.Write([]byte(l.Literal))
}
