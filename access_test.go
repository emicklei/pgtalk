package pgtalk

import (
	"testing"
)

func TestLiteral_String(t *testing.T) {
	l := newLiteralString("literal")
	ls := testSQL(l)
	if got, want := ls, "'literal'"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestStringWithNonNilFields_poly(t *testing.T) {
	p := poly{}
	f := 42.0
	p.FFloat = f
	if got, want := StringWithFields(p, false), "github.com/emicklei/pgtalk.poly{FFloat:42 }"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestStringWithNonNilFields_pointerpoly(t *testing.T) {
	p := new(poly)
	f := 42.0
	p.FFloat = f
	if got, want := StringWithFields(p, false), "github.com/emicklei/pgtalk.poly{FFloat:42 }"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestGetStringOfPoly(t *testing.T) {
	th := new(poly)
	th.expressionResults = map[string]any{
		"k": 42,
	}
	t.Log(StringWithFields(th, HideNilValues))
}
