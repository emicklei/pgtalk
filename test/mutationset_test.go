package test

import (
	"context"
	"testing"

	"github.com/emicklei/pgtalk"
	"github.com/emicklei/pgtalk/convert"
	"github.com/emicklei/pgtalk/test/tables/products"
)

func TestUpdate(t *testing.T) {
	m := products.Update(
		products.Code.Set(convert.StringToText("test")),
		products.CategoryId.Set(convert.Int4(1))).
		Where(products.ID.Equals(10))
	if got, want := oneliner(pgtalk.SQL(m)), `UPDATE public.products p1 SET code = $1,category_id = $2 WHERE (p1.id = 10)`; got != want {
		t.Log(diff(got, want))
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestUpdateReturning(t *testing.T) {
	m := products.Update(
		products.Code.Set(convert.StringToText("F42")),
		products.CategoryId.Set(convert.Int4(1))).
		Where(products.ID.Equals(1)).
		Returning(products.Code)
	if got, want := oneliner(pgtalk.SQL(m)), `UPDATE public.products p1 SET code = $1,category_id = $2 WHERE (p1.id = 1) RETURNING code`; got != want {
		t.Log(diff(got, want))
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	if testConnect == nil {
		return
	}
	it := m.Exec(context.Background(), testConnect)
	for it.HasNext() {
		p, err := it.Next()
		if err != nil {
			t.Fatal(it.Err())
		}
		t.Logf("%v", p.Code)
	}
}

func TestDelete(t *testing.T) {
	m := products.Delete().Where(products.ID.Equals(10))
	if got, want := oneliner(pgtalk.SQL(m)), `DELETE FROM public.products p1 WHERE (p1.id = 10)`; got != want {
		t.Log(diff(got, want))
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestDeleteWithParameter(t *testing.T) {
	par := pgtalk.NewParameter(10)

	m := products.Delete().Where(products.ID.Equals(par))
	if got, want := oneliner(pgtalk.SQL(m)), `DELETE FROM public.products p1 WHERE (p1.id = ?)`; got != want {
		t.Log(diff(got, want))
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestInsert(t *testing.T) {
	m := products.Insert(
		products.ID.Set(10),
		products.Code.Set(convert.StringToText("test")),
		products.CategoryId.Set(convert.Int4(1)))
	if got, want := oneliner(pgtalk.SQL(m)), `INSERT INTO public.products (id,code,category_id) VALUES ($1,$2,$3)`; got != want {
		t.Log(diff(got, want))
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}
