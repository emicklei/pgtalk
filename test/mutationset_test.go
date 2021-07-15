package test

import (
	"bytes"
	"fmt"
	"io"
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
	if got, want := pgtalk.SQL(m), `UPDATE public.products SET id = 10,code = 'test',category_id = 1 WHERE (p1.id = 10)`; got != want {
		t.Log(diff(got, want))
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestDelete(t *testing.T) {
	m := products.Delete().Where(products.ID.Equals(10))
	if got, want := pgtalk.SQL(m), `DELETE FROM public.products WHERE (p1.id = 10)`; got != want {
		t.Log(diff(got, want))
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestInsert(t *testing.T) {
	m := products.Insert(
		products.ID.Set(10),
		products.Code.Set("test"),
		products.Category_id.Set(1))
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
