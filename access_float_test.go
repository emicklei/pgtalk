package pgtalk

import (
	"testing"
)

func TestNewFloat64Access(t *testing.T) {
	info := MakeColumnInfo(polyTable, "ffloat", NotPrimary, Nullable, 1)
	a := NewFloat64Access(info, func(dest any) any { return &dest.(*poly).FFloat })
	if got, want := a.Column().columnName, "ffloat"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestFloat64Access_Set(t *testing.T) {
	a := polyFFloat.Set(1.23)
	if got, want := a.valueToInsert, 1.23; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestFloat64Access_Equals(t *testing.T) {
	a := polyFFloat.Equals(1.23)
	sql := SQL(a)
	if got, want := sql, `(p1.ffloat = 1.23)`; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestFloat64Access_Compare(t *testing.T) {
	a := polyFFloat.Compare(">", 1.23)
	sql := SQL(a)
	if got, want := sql, `(p1.ffloat > 1.23)`; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestFloat64Access_FieldValueToScan(t *testing.T) {
	p := new(poly)
	a := polyFFloat
	field := a.FieldValueToScan(p)
	if _, ok := field.(*float64); !ok {
		t.Errorf("got [%T] want [*float64]", field)
	}
}

func TestFloat64Access_TableAlias(t *testing.T) {
	a := polyFFloat.TableAlias("p")
	if got, want := a.SQL(), "p.ffloat"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestFloat64Access_AppendScannable(t *testing.T) {
	a := polyFFloat
	list := []any{}
	list = a.AppendScannable(list)
	if got, want := len(list), 1; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestFloat64Access_Get(t *testing.T) {
	a := polyFFloat
	values := map[string]any{
		"ffloat": 1.23,
	}
	if got, want := a.Get(values), 1.23; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}
