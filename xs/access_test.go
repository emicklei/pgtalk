package xs

import "testing"

func TestLiteral_String(t *testing.T) {
	l := LiteralString("literal")
	ls := l.SQL(SQLContext{})
	if got, want := ls, "'literal'"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}
