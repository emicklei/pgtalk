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

func TestQueryWithArguments(t *testing.T) {
	q := MakeQuerySet[poly](polyTable, polyTable.Columns)
	q.selectors = []ColumnAccessor{polyFUUID}
	q, arg := q.NewParameter(42)
	q = q.Where(polyFUUID.Equals(arg))
	fmt.Println(oneliner(SQL(q)))
	if got, want := len(q.queryParameters), 1; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}
