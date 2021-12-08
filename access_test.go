package pgtalk

import (
	"testing"
	"time"
)

func TestPoly_float64(t *testing.T) {
	p := new(poly)
	a := NewFieldAccess[float64](
		MakeColumnInfo(TableInfo{}, "TestPoly", false, false, 1),
		func(dest interface{}, f *float64) { dest.(*poly).FFloat = f })
	if err := a.SetFieldValue(p, 42.0); err != nil {
		t.Fatal(err)
	}
	if got, want := *p.FFloat, 42.0; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	forty2 := 42.0
	if err := a.SetFieldValue(p, &forty2); err != nil {
		t.Fatal(err)
	}
	if got, want := *p.FFloat, 42.0; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestPoly_time(t *testing.T) {
	p := new(poly)
	a := NewFieldAccess[time.Time](
		MakeColumnInfo(TableInfo{}, "TestPoly", false, false, 1),
		func(dest interface{}, v *time.Time) { dest.(*poly).FTime = v })
	n := time.Now()
	if err := a.SetFieldValue(p, n); err != nil {
		t.Fatal(err)
	}
	if got, want := *p.FTime, n; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestPoly_bool(t *testing.T) {
	p := new(poly)
	a := NewFieldAccess[bool](
		MakeColumnInfo(TableInfo{}, "TestPoly", false, false, 1),
		func(dest interface{}, v *bool) { dest.(*poly).FBool = v })
	if err := a.SetFieldValue(p, true); err != nil {
		t.Fatal(err)
	}
	if got, want := *p.FBool, true; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestLiteral_String(t *testing.T) {
	l := LiteralString("literal")
	ls := SQL(l)
	if got, want := ls, "'literal'"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestStringWithNonNilFields_poly(t *testing.T) {
	p := poly{}
	f := 42.0
	p.FFloat = &f
	if got, want := StringWithFields(p, false), "github.com/emicklei/pgtalk.poly{FFloat:42 }"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestStringWithNonNilFields_pointerpoly(t *testing.T) {
	p := new(poly)
	f := 42.0
	p.FFloat = &f
	if got, want := StringWithFields(p, false), "github.com/emicklei/pgtalk.poly{FFloat:42 }"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}
