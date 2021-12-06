package pgtalk

import (
	"bytes"
	"fmt"
	"io"
	"testing"
	"time"
)

type poly struct {
	fTime  *time.Time
	fFloat *float64
	fBool  *bool
}

func TestPoly_float64(t *testing.T) {
	p := new(poly)
	a := NewFieldAccess[float64](
		MakeColumnInfo(TableInfo{}, "TestPoly", false, false, 1),
		func(dest interface{}, f *float64) { dest.(*poly).fFloat = f })
	if err := a.SetFieldValue(p, 42.0); err != nil {
		t.Fatal(err)
	}
	if got, want := *p.fFloat, 42.0; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	forty2 := 42.0
	if err := a.SetFieldValue(p, &forty2); err != nil {
		t.Fatal(err)
	}
	if got, want := *p.fFloat, 42.0; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestPoly_time(t *testing.T) {
	p := new(poly)
	a := NewFieldAccess[time.Time](
		MakeColumnInfo(TableInfo{}, "TestPoly", false, false, 1),
		func(dest interface{}, v *time.Time) { dest.(*poly).fTime = v })
	n := time.Now()
	if err := a.SetFieldValue(p, n); err != nil {
		t.Fatal(err)
	}
	if got, want := *p.fTime, n; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestPoly_bool(t *testing.T) {
	p := new(poly)
	a := NewFieldAccess[bool](
		MakeColumnInfo(TableInfo{}, "TestPoly", false, false, 1),
		func(dest interface{}, v *bool) { dest.(*poly).fBool = v })
	if err := a.SetFieldValue(p, true); err != nil {
		t.Fatal(err)
	}
	if got, want := *p.fBool, true; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestLiteral_String(t *testing.T) {
	l := LiteralString("literal")
	ls := SQL(l)
	if got, want := ls, "'literal'"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func diff(left, right string) string {
	//assume one line
	b := new(bytes.Buffer)
	io.WriteString(b, "\n")
	io.WriteString(b, left)
	io.WriteString(b, "\n")
	leftRunes := []rune(left)
	rightRunes := []rune(right)
	size := len(leftRunes)
	if l := len(rightRunes); l < size {
		size = l
	}
	for c := 0; c < size; c++ {
		l := leftRunes[c]
		r := rightRunes[c]
		if l == r {
			b.WriteRune(l)
		} else {
			fmt.Fprintf(b, "^(%s)...", string(r))
			break
		}
	}
	return b.String()
}
