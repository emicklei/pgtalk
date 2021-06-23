package xs

import (
	"testing"
)

func TestEquals_String(t *testing.T) {
	e := Equals{Left: TextAccess{name: "test"}, Right: LiteralString("42")}
	if got, want := e.String(), "(test = '42')"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}
