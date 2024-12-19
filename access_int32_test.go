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
