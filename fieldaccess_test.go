package pgtalk

import (
	"testing"
	"time"

	"github.com/jackc/pgtype"
)

// gotip test -v -gcflags=-G=3
func TestFieldAccess_Daterange_WriteInto(t *testing.T) {
	type entity struct {
		dr *pgtype.Daterange
	}
	p := new(entity)
	a := NewFieldAccess[pgtype.Daterange](
		MakeColumnInfo(TableInfo{}, "TestFieldAccess_Daterange_WriteInto", false, false, 1),
		func(dest interface{}, v *pgtype.Daterange) { dest.(*entity).dr = v })
	rv := &pgtype.Daterange{
		Lower:     pgtype.Date{Time: time.Date(1990, 12, 31, 0, 0, 0, 0, time.UTC), Status: pgtype.Present},
		Upper:     pgtype.Date{Time: time.Date(2028, 1, 1, 0, 0, 0, 0, time.UTC), Status: pgtype.Present},
		LowerType: pgtype.Inclusive,
		UpperType: pgtype.Exclusive,
		Status:    pgtype.Present,
	}
	err := a.SetFieldValue(p, rv)
	if err != nil {
		t.Fatal(err)
	}
	if p.dr == nil {
		t.Fatal()
	}
	if got, want := p.dr.Status, pgtype.Present; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}
