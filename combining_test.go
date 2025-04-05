package pgtalk

import (
	"testing"
)

func TestUnionSQL(t *testing.T) {
	u := queryCombination{
		leftSet:  MakeQuerySet[poly](polyTable, polyTable.Columns),
		operator: "UNION ALL",
		rightSet: MakeQuerySet[poly](polyTable, polyTable.Columns),
	}
	if got, want := SQL(u), "((SELECT FROM public.polies p1) UNION ALL (SELECT FROM public.polies p1))"; got != want {
		t.Errorf("got [%[1]v:%[1]T] want [%[2]v:%[2]T]", got, want)
	}
}
