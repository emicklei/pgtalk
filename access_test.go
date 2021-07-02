package pgtalk

import "testing"

func TestLiteral_String(t *testing.T) {
	l := LiteralString("literal")
	ls := SQL(l)
	if got, want := ls, "'literal'"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}
