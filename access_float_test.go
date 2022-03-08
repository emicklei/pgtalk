package pgtalk

import "testing"

// gotip test -v -gcflags=-G=3
func TestFloat64Access_WriteInto(t *testing.T) {
	type price struct {
		amount float64
	}
	p := new(price)
	a := NewFloat64Access(
		MakeColumnInfo(TableInfo{}, "TestFloat64Access_WriteInto", false, false, 1),
		func(dest interface{}, f float64) { dest.(*price).amount = f }, nil)
	a.SetFieldValue(p, 42.0)
	if got, want := p.amount, 42.0; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestGenericFloat64Access_WriteInto(t *testing.T) {
	type price struct {
		amount float64
	}
	p := new(price)
	a := NewFieldAccess(
		MakeColumnInfo(TableInfo{}, "TestGenericFloat64Access_WriteInto", false, false, 1),
		nil, func(dest interface{}, f float64) { dest.(*price).amount = f })
	if err := a.SetFieldValue(p, 42.0); err != nil {
		t.Fatal(err)
	}
	if got, want := p.amount, 42.0; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	forty2 := 42.0
	if err := a.SetFieldValue(p, forty2); err != nil {
		t.Fatal(err)
	}
	if got, want := p.amount, 42.0; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}
