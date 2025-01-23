package pgtalk

import "testing"

func TestInt32InEmpty(t *testing.T) {
	ta := NewInt32Access(ci, nil)
	in := ta.In()
	sql := SQL(in)
	if got, want := sql, "false"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}
func TestInt32Euals(t *testing.T) {
	ta := NewInt32Access(ci, nil)
	in := ta.Equals(int32(1))
	sql := SQL(in)
	if got, want := sql, "(t1.label = 1)"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}
func TestIntEuals(t *testing.T) {
	ta := NewInt32Access(ci, nil)
	in := ta.Equals(int(1))
	sql := SQL(in)
	if got, want := sql, "(t1.label = 1)"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}
