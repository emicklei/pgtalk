package pgtalk

import (
	"bytes"
	"testing"
)

var (
	ti = TableInfo{
		Name:   "things",
		Schema: "public",
		Alias:  "t1",
	}
	ci = ColumnInfo{
		tableInfo:  ti,
		columnName: "label",
	}
)

func TestTextIn(t *testing.T) {
	b := new(bytes.Buffer)
	ta := NewTextAccess(ci, nil)
	in := ta.In("a", "b")
	in.SQLOn(b)
	if got, want := b.String(), "(t1.label IN ('a','b'))"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestTextLike(t *testing.T) {
	b := new(bytes.Buffer)
	ta := NewTextAccess(ci, nil)
	in := ta.Like("*me")
	in.SQLOn(b)
	if got, want := b.String(), "(t1.label LIKE '*me')"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestTextNotNull(t *testing.T) {
	b := new(bytes.Buffer)
	ta := NewTextAccess(ci, nil, nil)
	op := IsNotNull(ta)
	op.SQLOn(b)
	if got, want := b.String(), "(t1.label IS NOT NULL)"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestTextEquals(t *testing.T) {
	b := new(bytes.Buffer)
	ta := NewTextAccess(ci, nil)
	op := ta.Equals("me")
	op.SQLOn(b)
	if got, want := b.String(), "(t1.label = 'me')"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}
