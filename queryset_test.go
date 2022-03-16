package pgtalk

import (
	"testing"
)

func TestQuerySetSelect(t *testing.T) {
	q := MakeQuerySet[poly](polyTable, polyTable.Columns)
	q.selectors = polyColumns
	q.limit = 1
	q.offset = 2
	q.condition = IsNotNull(polyFUUID)
	q.orderBy = []ColumnAccessor{polyFUUID}
	q = q.Alias("ppp")
	if got, want := SQL(q), "SELECT ppp.ftime,ppp.ffloat FROM public.polies ppp WHERE (ppp.fuuid IS NOT NULL) ORDER BY ppp.fuuid LIMIT 1 OFFSET 2"; got != want {
		t.Errorf("got [%v]:%T want [%v]:%T", got, got, want, want)
	}
}
