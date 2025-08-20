package pgtalk

import (
	"testing"
)

func TestNewBooleanAccess(t *testing.T) {
	info := MakeColumnInfo(boolTable, "fbool", NotPrimary, Nullable, 1)
	a := NewBooleanAccess(info, func(dest any) any { return &dest.(*poly).FBool })
	if got, want := a.Column().columnName, "fbool"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

var boolTable = TableInfo{Name: "boolies", Schema: "public", Alias: "b1"}

var polyFBool = NewBooleanAccess(MakeColumnInfo(boolTable, "fbool", NotPrimary, Nullable, 1),
	func(dest any) any { return &dest.(*poly).FBool })

func TestBooleanAccess_ValueToInsert(t *testing.T) {
	a := polyFBool.Set(true)
	if got, want := a.ValueToInsert(), true; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestBooleanAccess_Equals(t *testing.T) {
	a := polyFBool.Equals(true)
	sql := SQL(a)
	if got, want := sql, `(b1.fbool = true)`; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestBooleanAccess_And(t *testing.T) {
	a := polyFBool.Equals(true)
	b := polyFBool.Equals(false)
	c := a.And(b)
	sql := SQL(c)
	if got, want := sql, `((b1.fbool = true) AND (b1.fbool = false))`; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestBooleanAccess_Get(t *testing.T) {
	a := polyFBool
	values := map[string]any{
		"fbool": true,
	}
	if got, want := a.Get(values), true; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestBooleanAccess_AppendScannable(t *testing.T) {
	a := polyFBool
	list := []any{}
	list = a.AppendScannable(list)
	if got, want := len(list), 1; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestBooleanAccess_FieldValueToScan(t *testing.T) {
	p := new(poly)
	a := polyFBool
	field := a.FieldValueToScan(p)
	if _, ok := field.(*bool); !ok {
		t.Errorf("got [%T] want [*bool]", field)
	}
}

func TestBooleanAccess_TableAlias(t *testing.T) {
	a := polyFBool.TableAlias("p")
	if got, want := a.SQL(), "p.fbool"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}
