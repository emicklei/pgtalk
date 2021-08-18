package test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"testing"

	"github.com/emicklei/pgtalk"
	"github.com/emicklei/pgtalk/test/products"
)

func TestUpdate(t *testing.T) {
	m := products.Update(
		products.Code.Set("test"),
		products.CategoryId.Set(1)).
		Where(products.ID.Equals(10))
	if got, want := pgtalk.SQL(m), `UPDATE public.products p1 SET code = $1,category_id = $2 WHERE (p1.id = 10)`; got != want {
		t.Log(diff(got, want))
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestUpdateReturning(t *testing.T) {
	m := products.Update(
		products.Code.Set("F42"),
		products.CategoryId.Set(1)).
		Where(products.ID.Equals(1)).
		Returning(products.Code)
	if got, want := pgtalk.SQL(m), `UPDATE public.products p1 SET code = $1,category_id = $2 WHERE (p1.id = 1) RETURNING code`; got != want {
		t.Log(diff(got, want))
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	if testConnect == nil {
		return
	}
	it := m.Exec(context.Background(), testConnect)
	for it.HasNext() {
		p, err := products.Next(it)
		if err != nil {
			t.Fatal(it.Err())
		}
		t.Logf("%s", *p.Code)
	}
}

func TestDelete(t *testing.T) {
	m := products.Delete().Where(products.ID.Equals(10))
	if got, want := pgtalk.SQL(m), `DELETE FROM public.products p1 WHERE (p1.id = 10)`; got != want {
		t.Log(diff(got, want))
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestInsert(t *testing.T) {
	m := products.Insert(
		products.ID.Set(10),
		products.Code.Set("test"),
		products.CategoryId.Set(1))
	if got, want := pgtalk.SQL(m), `INSERT INTO public.products (id,code,category_id) VALUES ($1,$2,$3)`; got != want {
		t.Log(diff(got, want))
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func diff(left, right string) string {
	//assume one line
	b := new(bytes.Buffer)
	io.WriteString(b, "\n")
	io.WriteString(b, left)
	io.WriteString(b, "\n")
	leftRunes := []rune(left)
	rightRunes := []rune(right)
	size := len(leftRunes)
	if l := len(rightRunes); l < size {
		size = l
	}
	for c := 0; c < size; c++ {
		l := leftRunes[c]
		r := rightRunes[c]
		if l == r {
			b.WriteRune(l)
		} else {
			fmt.Fprintf(b, "^(%s)...", string(r))
			break
		}
	}
	return b.String()
}
