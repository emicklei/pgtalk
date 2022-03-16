package pgtalk

import "testing"

func TestTextOperators(t *testing.T) {
	a := NewTextAccess(ColumnInfo{ti, "col", false, false, false, 1}, nil)
	if got, want := testSQL(a.Equals("help")), "(t1.col = 'help')"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	if got, want := testSQL(a.In("a", "b")), "(t1.col IN ('a','b'))"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	if got, want := testSQL(a.Compare("<", "b")), "(t1.col < 'b')"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	if got, want := testSQL(a.Like("*b")), "(t1.col LIKE '*b')"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	if got, want := testSQL(IsNotNull(a)), "(t1.col IS NOT NULL)"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}
