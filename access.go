package pgtalk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/jackc/pgtype"
)

var EmptyColumnAccessor = []ColumnAccessor{}

type valuePrinter struct {
	v any
}

func MakeValuePrinter(v any) valuePrinter { return valuePrinter{v: v} }

func (p valuePrinter) SQLOn(w WriteContext) {
	if e, ok := p.v.(SQLWriter); ok {
		e.SQLOn(w)
		return
	}
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
	fmt.Fprintf(w, "%v", p.v)
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

type valuesPrinter struct {
	vs []any
}

func (p valuesPrinter) SQLOn(w WriteContext) {
	fmt.Fprintf(w, "(")
	for i, each := range p.vs {
		if i > 0 {
			fmt.Fprintf(w, ",")
		}
		valuePrinter{each}.SQLOn(w)
	}
	fmt.Fprintf(w, ")")
}

type LiteralString string

func (l LiteralString) SQLOn(w WriteContext) {
	io.WriteString(w, "'")
	io.WriteString(w, string(l))
	io.WriteString(w, "'")
}

type NoCondition struct{}

var EmptyCondition = NoCondition{}

func (n NoCondition) SQLOn(w WriteContext) {}

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
		fmt.Fprintf(b, "%s:%v ", f.Name, fi)
	}
	fmt.Fprint(b, "}")
	return b.String()
}
