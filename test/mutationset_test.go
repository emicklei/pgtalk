package test

import (
	"testing"

	"github.com/emicklei/pgtalk"
	"github.com/emicklei/pgtalk/test/products"
)

func TestUpdate(t *testing.T) {
	m := products.Update(
		products.ID.Value(10),
		products.Code.Value("test"),
		products.CategoryID.Value(1)).
		Where(products.ID.Equals(10))
	if got, want := pgtalk.SQL(m), `UPDATE products SET (t1.id = 10) WHERE (t1.id = 10)`; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestDelete(t *testing.T) {
	m := products.Delete().Where(products.ID.Equals(10))
	if got, want := pgtalk.SQL(m), `DELETE FROM products WHERE (t1.id = 10)`; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestInsert(t *testing.T) {
	m := products.Insert(
		products.ID.Value(10),
		products.Code.Value("test"),
		products.CategoryID.Value(1))
	if got, want := pgtalk.SQL(m), `INSERT INTO products (id,code,category_id) values ($1,$2,$3)`; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}
