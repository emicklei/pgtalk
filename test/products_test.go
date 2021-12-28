package test

import (
	"context"
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
	if got, want := pgtalk.SQL(q), `SELECT p1.id,p1.code FROM public.products p1 WHERE ((p1.code = 'F42') AND (p1.id = 1)) LIMIT 1`; got != want {
		t.Log(diff(got, want))
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	if testConnect == nil {
		return
	}
	products, err := q.Exec(context.Background(), testConnect)
	if len(products) != 1 {
		t.Log("empty results")
		return
	}
	log.Printf("%v,%v,%v", products[0].ID, products[0].Code, err)
}

func TestSelectAllColumns(t *testing.T) {
	q := products.
		Select(products.AllColumns()...).
		Limit(2)
	if got, want := pgtalk.SQL(q), `SELECT p1.id,p1.created_at,p1.updated_at,p1.deleted_at,p1.code,p1.price,p1.category_id FROM public.products p1 LIMIT 2`; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestIn(t *testing.T) {
	q := products.
		Select(products.AllColumns()...).
		Where(products.Code.In("F42", "f42"))
	if got, want := pgtalk.SQL(q), `SELECT p1.id,p1.created_at,p1.updated_at,p1.deleted_at,p1.code,p1.price,p1.category_id FROM public.products p1 WHERE (p1.code IN ('F42','f42'))`; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestExists(t *testing.T) {
	q := products.
		Select(products.ID).
		Where(categories.Select(categories.ID).Exists())
	if got, want := pgtalk.SQL(q), `SELECT p1.id FROM public.products p1 WHERE EXISTS (SELECT c1.id FROM public.categories c1)`; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestInnerJoin(t *testing.T) {
	q := products.Select(products.Code).Where(products.Code.Equals("F42")).
		Join(categories.Select(categories.Title)).
		On(products.ID.Equals(categories.ID))
	if got, want := pgtalk.SQL(q), `SELECT p1.code,c1.title FROM public.products p1 INNER JOIN public.categories c1 ON (p1.id = c1.id) WHERE (p1.code = 'F42')`; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	if testConnect == nil {
		return
	}
	it, _ := q.Exec(context.Background(), testConnect)
	for it.HasNext() {
		p := new(products.Product)
		c := new(categories.Category)
		_ = it.Next(p, c)
		t.Logf("%#v,%#v", p.Code, c.Title)
	}
}

func TestLeftJoin(t *testing.T) {
	q := products.Select(products.Code).Where(products.Code.Equals("F42")).
		LeftOuterJoin(categories.Select(categories.Title)).
		On(products.ID.Equals(categories.ID))
	if got, want := pgtalk.SQL(q), `SELECT p1.code,c1.title FROM public.products p1 LEFT OUTER JOIN public.categories c1 ON (p1.id = c1.id) WHERE (p1.code = 'F42')`; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestMultiLeftJoin(t *testing.T) {
	q := products.Select(products.Code).Where(products.Code.Equals("F42")).
		LeftOuterJoin(categories.Select(categories.Title)).
		On(products.CategoryId.Equals(categories.ID)).
		LeftOuterJoin(categories.Select(categories.Title)).
		On(products.CategoryId.Equals(categories.ID))
	if got, want := pgtalk.SQL(q), `SELECT p1.code,c1.title,c1.title FROM public.products p1 LEFT OUTER JOIN public.categories c1 ON (p1.category_id = c1.id) LEFT OUTER JOIN public.categories c1 ON (p1.category_id = c1.id) WHERE (p1.code = 'F42')`; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
		t.Log(diff(got, want))
	}
}

func TestFullSelect(t *testing.T) {
	q := products.
		Select(products.AllColumns()...).
		Distinct().
		Where(products.Code.Compare(">", "A").And(pgtalk.IsNotNull(products.CategoryId))).
		GroupBy(products.CategoryId).
		OrderBy(products.CategoryId).
		Ascending().
		Limit(3)
	if got, want := pgtalk.SQL(q), `SELECT DISTINCT p1.id,p1.created_at,p1.updated_at,p1.deleted_at,p1.code,p1.price,p1.category_id FROM public.products p1 WHERE ((p1.code > 'A') AND (p1.category_id IS NOT NULL)) GROUP BY p1.category_id ORDER BY p1.category_id LIMIT 3`; got != want {
		t.Log(diff(got, want))
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}
