package pgtalk

import (
	"testing"
)

func TestPretty(t *testing.T) {
	sql := `SELECT DISTINCT p1.id,p1.created_at,p1.updated_at,p1.deleted_at,p1.code,p1.price,p1.category_id FROM products p1 WHERE ((p1.code > 'A') AND (p1.category_id IS NOT NULL)) GROUP BY p1.category_id ORDER BY p1.category_id`
	t.Log(PrettySQL(MakeValuePrinter(sql)))
}
