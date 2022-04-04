package pgtalk

import (
	"fmt"
	"strings"
	"testing"
)

func TestQuerySetSelect(t *testing.T) {
	q := MakeQuerySet[poly](polyTable, polyTable.Columns)
	q.selectors = polyColumns
	q.limit = 1
	q.offset = 2
	q.condition = IsNotNull(polyFUUID)
	q.orderBy = []ColumnAccessor{polyFUUID}
	q = q.TableAlias("ppp")
	fmt.Println(SQL(q))
	if got, want := oneliner(SQL(q)), "SELECT ppp.ftime, ppp.ffloat FROM public.polies ppp WHERE (ppp.fuuid IS NOT NULL) ORDER BY ppp.fuuid LIMIT 1 OFFSET 2"; got != want {
		t.Log(diff(got, want))
		t.Errorf("got [%v]:%T want [%v]:%T", got, got, want, want)
	}
}

func oneliner(s string) string {
	return strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(s, "\t", " "), "\n", " "), "  ", " ")
}

func TestQueryWithParameter(t *testing.T) {
	q := MakeQuerySet[poly](polyTable, polyTable.Columns)

	s := NewParameterSet()
	i42 := s.NewParameter(42)

	q.selectors = []ColumnAccessor{polyFUUID}
	q = q.Where(polyFUUID.Equals(i42))
	if got, want := oneliner(SQL(q)), "SELECT p1.fuuid FROM public.polies p1 WHERE (p1.fuuid = ?)"; got != want {
		t.Errorf("got [%v]:%T want [%v]:%T", got, got, want, want)
	}
	mock := newMockConnection(t)
	q.Exec(mock.ctx(), mock, s.Parameters()...)
}
