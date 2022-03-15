package pgtalk

import (
	"testing"
)

func TestPretty(t *testing.T) {
	sql := `SELECT DISTINCT p1.id,p1.created_at,p1.updated_at,p1.deleted_at,p1.code,p1.price,p1.category_id FROM products p1 WHERE ((p1.code > 'A') AND (p1.category_id IS NOT NULL)) GROUP BY p1.category_id ORDER BY p1.category_id`
	t.Log(PrettySQL(MakeValuePrinter(sql)))
}

func TestQuerySetSelect(t *testing.T) {
	q := MakeQuerySet[poly](polyTable, polyTable.Columns)
	q.selectors = polyColumns
	q.limit = 1
	q.offset = 2
	q.condition = IsNotNull(polyFUUID)
	q.orderBy = []ColumnAccessor{polyFUUID}
	if got, want := SQL(q), "SELECT p1.ftime,p1.ffloat FROM public.polies p1 WHERE (p1.fuuid IS NOT NULL) ORDER BY p1.fuuid LIMIT 1 OFFSET 2"; got != want {
		t.Errorf("got [%v]:%T want [%v]:%T", got, got, want, want)
	}
}
