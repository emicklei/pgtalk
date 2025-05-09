package pgtalk

import (
	"testing"

	"github.com/emicklei/pgtalk/convert"
	"github.com/jackc/pgx/v5/pgtype"
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

func TestUUID_IN_SQL(t *testing.T) {
	a := NewFieldAccess[pgtype.UUID](ColumnInfo{ti, "col", false, false, false}, nil)
	ids := []pgtype.UUID{
		convert.StringToUUID("b344a1918d0cbd1542de669644dd1bfd"),
	}
	ex := a.In(ids...)
	sql := SQL(ex)
	t.Log(sql)
}

func TestNewColumnsAdd(t *testing.T) {
	cols := NewColumns(polyColumns...)
	if len(cols) != 2 {
		t.Error("expected 2 columns")
	}
	cols.Add(polyColumns[0])
	if len(cols) != 3 {
		t.Error("expected 3 columns")
	}
}
