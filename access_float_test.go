package pgtalk

import "testing"

func TestFloat64Access_WriteInto(t *testing.T) {
	type price struct {
		amount *float64
	}
	p := new(price)
	a := NewFloat64Access(TableInfo{}, "TestFloat64Access_WriteInto", func(dest interface{}, f *float64) { dest.(*price).amount = f })
	a.WriteInto(p, 42.0)
	if got, want := *p.amount, 42.0; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}
