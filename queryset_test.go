package pgtalk

import (
	"strings"
	"testing"
)

func TestQuerySetSelect(t *testing.T) {
	q := MakeQuerySet[poly](polyTable, polyTable.Columns)
	q.selectors = polyColumns
	q.limit = 1
	q.offset = 2
	q.condition = IsNotNull(polyFUUID)
	q.orderBy = []SQLWriter{polyFUUID}
	q = q.TableAlias("ppp")
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

	i42 := NewParameter(42)

	q.selectors = []ColumnAccessor{polyFUUID}
	q = q.Where(polyFUUID.Equals(i42))
	if got, want := oneliner(SQL(q)), "SELECT p1.fuuid FROM public.polies p1 WHERE (p1.fuuid = ?)"; got != want {
		t.Errorf("got [%v]:%T want [%v]:%T", got, got, want, want)
	}
	mock := newMockConnection(t)
	q.Exec(mock.ctx(), mock, i42)
}

func TestSelectForSkipLocked(t *testing.T) {
	q := MakeQuerySet[poly](polyTable, polyTable.Columns)
	q = q.For(FOR_UPDATE)
	q = q.SkipLocked()
	q = q.Where(polyFUUID.Equals(1))
	if got, want := oneliner(SQL(q)), "SELECT FROM public.polies p1 WHERE (p1.fuuid = 1) FOR UPDATE SKIP LOCKED"; got != want {
		t.Errorf("got [%v]:%T want [%v]:%T", got, got, want, want)
	}
}
