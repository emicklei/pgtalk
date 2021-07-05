package test

import (
	"testing"

	"github.com/emicklei/pgtalk"
	"github.com/emicklei/pgtalk/test/products"
)

func TestUpdate(t *testing.T) {
	m := products.Update(
		products.ID.Set(10),
		products.Code.Set("test"),
		products.Category_id.Set(1)).
		Where(products.ID.Equals(10))
	if got, want := pgtalk.SQL(m), `UPDATE products SET id = 10,code = 'test',category_id = 1 WHERE (p1.id = 10)`; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestDelete(t *testing.T) {
	m := products.Delete().Where(products.ID.Equals(10))
	if got, want := pgtalk.SQL(m), `DELETE FROM products WHERE (p1.id = 10)`; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestInsert(t *testing.T) {
	m := products.Insert(
		products.ID.Set(10),
		products.Code.Set("test"),
		products.Category_id.Set(1))
	if got, want := pgtalk.SQL(m), `INSERT INTO products (id,code,category_id) values ($1,$2,$3)`; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}
