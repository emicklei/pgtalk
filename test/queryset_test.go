package test

import (
	"log"
	"testing"

	"github.com/emicklei/pgtalk/test/categories"
	"github.com/emicklei/pgtalk/test/products"
)

func TestSelectProductsWhere(t *testing.T) {
	q := products.
		Select(products.ID, products.Code).
		Where(products.Code.Equals("F42").
			And(products.ID.Equals(1))).
		Limit(1)
	if got, want := q.SQL(), `SELECT t1.id,t1.code FROM products t1 WHERE ((t1.code = 'F42') AND (t1.id = 1)) LIMIT 1`; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	products, err := q.Exec(testConnect)
	log.Printf("%v,%v,%v", *products[0].ID, *products[0].Code, err)
}

func TestSelectAllColumns(t *testing.T) {
	q := products.
		Select(products.AllColumns...).
		Limit(2)
	if got, want := q.SQL(), `SELECT t1.id,t1.code,t1.category_id FROM products t1 LIMIT 2`; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestInnerJoin(t *testing.T) {
	q := products.Select(products.Code).Where(products.Code.Equals("F42")).
		Join(categories.Select(categories.Title)).
		On(products.ID, categories.ID)
	if got, want := q.SQL(), `SELECT t1.code,t2.title FROM products t1 INNER JOIN categories t2 ON (t1.id = t2.id) WHERE (t1.code = 'F42')`; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}

	it, _ := q.Exec(testConnect)
	for it.HasNext() {
		p := new(products.Product)
		c := new(categories.Category)
		_ = it.Next(p, c)
		t.Logf("%#v,%#v", *p.Code, *c.Title)
	}
}

func TestLeftJoin(t *testing.T) {
	t.Skip() // TODO
	q := products.Select(products.Code).Where(products.Code.Equals("F42")).
		LeftJoin(categories.Select(categories.Title)).
		On(products.ID, categories.ID)
	if got, want := q.SQL(), `SELECT t1.code,t2.title FROM products t1 LEFT JOIN categories t2 ON (t1.id = t2.id) WHERE (t1.code = 'F42')`; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

type ProductWithCount struct {
	*products.Product
	Count int
}

func TestSelectProductWithCount(t *testing.T) {
	t.Skip()
	q := products.Select(products.Code).Count(products.ID)
	if got, want := q.SQL(), `SELECT t1.code,COUNT(t1.id) FROM products t1`; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	it := q.Exec(testConnect)
	for it.HasNext() {
		pc := new(ProductWithCount)
		_ = it.Next(pc)
	}
}
