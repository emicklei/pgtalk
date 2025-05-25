package pgtalk

import (
	"testing"
)

func TestUnionSQL(t *testing.T) {
	left := MakeQuerySet[poly](polyTable, polyTable.Columns)
	right := MakeQuerySet[poly](polyTable, polyTable.Columns)
	u := left.Union(right, true)
	if got, want := SQL(u), "((SELECT FROM public.polies p1) UNION ALL (SELECT FROM public.polies p1))"; got != want {
		t.Errorf("got [%[1]v:%[1]T] want [%[2]v:%[2]T]", got, want)
	}
	middle := MakeQuerySet[poly](polyTable, polyTable.Columns)
	w := u.Except(middle)
	if got, want := w.String(), "((((SELECT FROM public.polies p1) UNION ALL (SELECT FROM public.polies p1))) EXCEPT (SELECT FROM public.polies p1))"; got != want {
		t.Errorf("got [%[1]v:%[1]T] want [%[2]v:%[2]T]", got, want)
	}
}
