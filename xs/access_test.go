package xs

import "testing"

func TestLiteral_String(t *testing.T) {
	l := LiteralString("literal")
	ls := l.String()
	if got, want := ls, "'literal'"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}
