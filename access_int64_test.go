package pgtalk

import "testing"

func TestInt64InEmpty(t *testing.T) {
	ta := NewInt64Access(ci, nil)
	in := ta.In()
	sql := SQL(in)
	if got, want := sql, "false"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}
