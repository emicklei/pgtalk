package pgtalk

import (
	"testing"
)


func TestTextIn(t *testing.T) {
	ta := NewTextAccess(ci, nil)
	in := ta.In("a", "b")
	sql := SQL(in)
	if got, want := sql, "(t1.label IN ('a','b'))"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestTextInEmpty(t *testing.T) {
	ta := NewTextAccess(ci, nil)
	in := ta.In()
	sql := SQL(in)
	if got, want := sql, "false"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestTextLike(t *testing.T) {
	ta := NewTextAccess(ci, nil)
	in := ta.Like("*me")
	sql := SQL(in)
	if got, want := sql, "(t1.label LIKE '*me')"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestTextNotNull(t *testing.T) {
	ta := NewTextAccess(ci, nil)
	op := IsNotNull(ta)
	sql := SQL(op)
	if got, want := sql, "(t1.label IS NOT NULL)"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestTextEquals(t *testing.T) {
	ta := NewTextAccess(ci, nil)
	op := ta.Equals("me")
	sql := SQL(op)
	if got, want := sql, "(t1.label = 'me')"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}
