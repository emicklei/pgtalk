package test

import (
	"log"
	"testing"

	"github.com/emicklei/pgtalk"
	"github.com/emicklei/pgtalk/test/categories"
	"github.com/emicklei/pgtalk/test/products"
)

func TestSelectProductsWhere(t *testing.T) {
	q := products.
		Select(products.ID, products.Code).
		Where(products.Code.Equals("F42").
			And(products.ID.Equals(1))).
		Limit(1)
	if got, want := pgtalk.SQL(q), `SELECT p1.id,p1.code FROM products p1 WHERE ((p1.code = 'F42') AND (p1.id = 1)) LIMIT 1`; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	if testConnect == nil {
		return
	}
	products, err := q.Exec(testConnect)
	log.Printf("%v,%v,%v", *products[0].ID, *products[0].Code, err)
}

func TestSelectAllColumns(t *testing.T) {
	q := products.
		Select(products.AllColumns()...).
		Limit(2)
	if got, want := pgtalk.SQL(q), `SELECT p1.id,p1.created_at,p1.updated_at,p1.deleted_at,p1.code,p1.price,p1.category_id FROM products p1 LIMIT 2`; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestIn(t *testing.T) {
	q := products.
		Select(products.AllColumns()...).
		Where(products.Code.In("F42", "f42"))
	if got, want := pgtalk.SQL(q), `SELECT p1.id,p1.created_at,p1.updated_at,p1.deleted_at,p1.code,p1.price,p1.category_id FROM products p1 WHERE (p1.code IN ('F42','f42'))`; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestInnerJoin(t *testing.T) {
	q := products.Select(products.Code).Where(products.Code.Equals("F42")).
		Join(categories.Select(categories.Title)).
		On(products.ID.Equals(categories.ID))
	if got, want := pgtalk.SQL(q), `SELECT p1.code,c1.title FROM products p1 INNER JOIN categories c1 ON (p1.id = c1.id) WHERE (p1.code = 'F42')`; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	if testConnect == nil {
		return
	}
	it, _ := q.Exec(testConnect)
	for it.HasNext() {
		p := new(products.Product)
		c := new(categories.Categorie)
		_ = it.Next(p, c)
		t.Logf("%#v,%#v", *p.Code, *c.Title)
	}
}

func TestLeftJoin(t *testing.T) {
	q := products.Select(products.Code).Where(products.Code.Equals("F42")).
		LeftOuterJoin(categories.Select(categories.Title)).
		On(products.ID.Equals(categories.ID))
	if got, want := pgtalk.SQL(q), `SELECT p1.code,c1.title FROM products p1 LEFT OUTER JOIN categories c1 ON (p1.id = c1.id) WHERE (p1.code = 'F42')`; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestMultiLeftJoin(t *testing.T) {
	q := products.Select(products.Code).Where(products.Code.Equals("F42")).
		LeftOuterJoin(categories.Select(categories.Title)).
		On(products.ID.Equals(categories.ID)).
		LeftOuterJoin(categories.Select(categories.Title)).
		On(products.ID.Equals(categories.ID))
	if got, want := pgtalk.SQL(q), `SELECT p1.code,c1.title,c1.title FROM products p1 LEFT OUTER JOIN categories c1 ON (p1.id = c1.id) LEFT OUTER JOIN categories c1 ON (p1.id = c1.id)`; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestFullSelect(t *testing.T) {
	q := products.
		Select(products.AllColumns()...).
		Distinct().
		Where(products.Code.Compare(">", "A").And(products.Category_id.NotNull())).
		GroupBy(products.Category_id).
		OrderBy(products.Category_id).
		Ascending()
	if got, want := pgtalk.SQL(q), `SELECT DISTINCT p1.id,p1.created_at,p1.updated_at,p1.deleted_at,p1.code,p1.price,p1.category_id FROM products p1 WHERE ((p1.code > 'A') AND (p1.category_id IS NOT NULL)) GROUP BY p1.category_id ORDER BY p1.category_id`; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

type ProductWithCount struct {
	*products.Product
	Count int
}

func TestSelectProductWithCount(t *testing.T) {
	q := products.Select(products.Code).Count(products.ID)
	if got, want := pgtalk.SQL(q), `SELECT p1.code,COUNT(p1.id) FROM products p1`; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	if testConnect == nil {
		return
	}
	it := q.Exec(testConnect)
	for it.HasNext() {
		pc := new(ProductWithCount)
		_ = it.Next(pc)
		t.Logf("%#v", pc)
	}
}
