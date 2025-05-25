package pgtalk

import "testing"

func TestLazySQL(t *testing.T) {
	q := MakeQuerySet[poly](polyTable, polyTable.Columns)
	q.condition = IsNotNull(polyFUUID)
	if got, want := q.String(), "SELECT FROM public.polies p1 WHERE (p1.fuuid IS NOT NULL)"; got != want {
		t.Errorf("got [%[1]v:%[1]T] want [%[2]v:%[2]T]", got, want)
	}
}
