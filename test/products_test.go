package test

import (
	"context"
	"log"
	"testing"

	"github.com/emicklei/pgtalk"
	"github.com/emicklei/pgtalk/convert"
	"github.com/emicklei/pgtalk/test/tables/categories"
	"github.com/emicklei/pgtalk/test/tables/products"
)

func TestSelectProductsWhere(t *testing.T) {
	q := products.
		Select(products.ID, products.Code).
		Where(products.Code.Equals(convert.StringToText("F42")).
			And(products.ID.Equals(1))).
		Limit(1)
	if got, want := oneliner(pgtalk.SQL(q)), `SELECT p1.id, p1.code FROM public.products p1 WHERE ((p1.code = 'F42') AND (p1.id = 1)) LIMIT 1`; got != want {
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

func TestSelectProductsWhereIterator(t *testing.T) {
	createCategory(t, 24)
	createProduct(t, 13, 24)
	q := products.
		Select(products.ID, products.Code).
		Limit(1)
	if testConnect == nil {
		return
	}
	productsIt, err := q.Iterate(context.Background(), testConnect)
	if err != nil {
		t.Fatal(err)
	}
	for productsIt.HasNext() {
		each, err := productsIt.Next()
		if err != nil {
			t.Fatal(err)
		}
		t.Log(each)
	}
}

func TestSelectAllColumns(t *testing.T) {
	q := products.
		Select(products.Columns()...).
		Limit(2)
	if got, want := oneliner(pgtalk.SQL(q)), `SELECT p1.category_id, p1.code, p1.created_at, p1.deleted_at, p1.id, p1.price, p1.updated_at FROM public.products p1 LIMIT 2`; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestIn(t *testing.T) {
	q := products.
		Select(products.Columns()...).
		Where(products.Code.In(convert.StringToText("F42"), convert.StringToText("f42")))
	if got, want := oneliner(pgtalk.SQL(q)), `SELECT p1.category_id, p1.code, p1.created_at, p1.deleted_at, p1.id, p1.price, p1.updated_at FROM public.products p1 WHERE (p1.code IN ('F42','f42'))`; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestExists(t *testing.T) {
	q := products.
		Select(products.ID).
		Where(categories.Select(categories.ID).Exists())
	if got, want := oneliner(pgtalk.SQL(q)), `SELECT p1.id FROM public.products p1 WHERE EXISTS (SELECT c1.id FROM public.categories c1)`; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestInnerJoin(t *testing.T) {
	createCategory(t, 23)
	createProduct(t, 12, 23)

	par := pgtalk.NewParameter("F42")

	q := products.Select(products.Code)
	j := q.Where(products.Code.Equals(par)).
		Join(categories.Select(categories.Title)).
		On(products.CategoryId.Equals(categories.ID))
	sql := oneliner(pgtalk.SQL(j))
	if got, want := sql, `SELECT p1.code, c1.title FROM public.products p1 INNER JOIN public.categories c1 ON (p1.category_id = c1.id) WHERE (p1.code = ?)`; got != want {
		t.Log(diff(got, want))
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	if testConnect == nil {
		return
	}
	it, err := j.Exec(context.Background(), testConnect, par)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := pgtalk.SQL(par), "$1"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	for it.HasNext() {
		p := new(products.Product)
		c := new(categories.Category)
		_ = it.Next(p, c)
		t.Logf("%s,%s\n", p.String(), c.String())
	}
}

func TestLeftJoin(t *testing.T) {
	q := products.Select(products.Code).Where(products.Code.Equals(convert.StringToText("F42"))).
		LeftOuterJoin(categories.Select(categories.Title)).
		On(products.ID.Equals(categories.ID))
	if got, want := oneliner(pgtalk.SQL(q)), `SELECT p1.code, c1.title FROM public.products p1 LEFT OUTER JOIN public.categories c1 ON (p1.id = c1.id) WHERE (p1.code = 'F42')`; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestMultiLeftJoin(t *testing.T) {
	q := products.Select(products.Code).Where(products.Code.Equals(convert.StringToText("F42"))).
		LeftOuterJoin(categories.Select(categories.Title)).
		On(products.CategoryId.Equals(categories.ID)).
		LeftOuterJoin(categories.Select(categories.Title)).
		On(products.CategoryId.Equals(categories.ID))
	if got, want := oneliner(pgtalk.SQL(q)), `SELECT p1.code, c1.title, c1.title FROM public.products p1 LEFT OUTER JOIN public.categories c1 ON (p1.category_id = c1.id) LEFT OUTER JOIN public.categories c1 ON (p1.category_id = c1.id) WHERE (p1.code = 'F42')`; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
		t.Log(diff(got, want))
	}
}

func TestFullSelect(t *testing.T) {
	q := products.
		Select(products.Columns()...).
		Distinct().
		Where(products.Code.Compare(">", "A").And(pgtalk.IsNotNull(products.CategoryId))).
		GroupBy(products.CategoryId).
		OrderBy(products.CategoryId).
		Ascending().
		Limit(3)
	if got, want := oneliner(pgtalk.SQL(q)), `SELECT DISTINCT p1.category_id, p1.code, p1.created_at, p1.deleted_at, p1.id, p1.price, p1.updated_at FROM public.products p1 WHERE ((p1.code > 'A') AND (p1.category_id IS NOT NULL)) GROUP BY p1.category_id ORDER BY p1.category_id ASC LIMIT 3`; got != want {
		t.Log(diff(got, want))
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	//fmt.Println(pgtalk.SQL(q))
}

func TestProductUpperCode(t *testing.T) {
	createProduct(t, 1234, 1)
	createProduct(t, 12345, 1)
	q := products.Select(products.ID, pgtalk.SQLAs("UPPER(p1.Code)", "upper"), pgtalk.SQLAs("p1.id", "rid"))
	t.Log(pgtalk.SQL(q))
	list, err := q.Exec(context.Background(), testConnect)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := list[0].GetExpressionResult("upper"), "F42"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	if got, want := list[1].GetExpressionResult("upper"), "F42"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	if got, want := list[0].GetExpressionResult("rid"), list[0].ID; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	if got, want := list[1].GetExpressionResult("rid"), list[1].ID; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestUpdateProductRowsAffected(t *testing.T) {
	if testConnect == nil {
		t.Skip("no connection")
	}
	createProduct(t, 123456, 1)

	// this is an example of collecting columns before calling products.Update
	cols := pgtalk.NewColumns()
	cols.Add(products.Code.Set(convert.StringToText("G1")))

	m := products.Update(cols...).Where(products.ID.Equals(123456))
	it := m.Exec(context.Background(), testConnect)
	if it.Err() != nil {
		t.Fatal(it.Err())
	}
	if got, want := it.CommandTag().RowsAffected(), int64(1); got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func createProduct(t *testing.T, id int32, categoryId int) {
	q := products.Insert(
		products.ID.Set(id),
		products.Code.Set(convert.StringToText("f42")),
		products.CategoryId.Set(convert.Int4(categoryId)))
	ctx := context.Background()
	tx, err := testConnect.Begin(ctx)
	if err != nil {
		t.Fatal(err)
	}
	it := q.Exec(ctx, testConnect)
	if it.Err() != nil {
		t.Fatal(it.Err())
	}
	if err := tx.Commit(ctx); err != nil {
		t.Fatal(err)
	}
}

func createCategory(t *testing.T, id int32) {
	q := categories.Insert(categories.ID.Set(id), categories.Title.Set(convert.StringToText("one")))
	ctx := context.Background()
	tx, err := testConnect.Begin(ctx)
	if err != nil {
		t.Fatal(err)
	}
	it := q.Exec(context.Background(), testConnect)
	if it.Err() != nil {
		t.Fatal(it.Err())
	}
	if err := tx.Commit(ctx); err != nil {
		t.Fatal(err)
	}
}

func TestSearchCategoryTitle(t *testing.T) {
	if testConnect == nil {
		t.Skip()
	}
	txt := "The lazy quick brown fox jumps over the quick lazy jumping dog"
	mut := categories.Insert(
		categories.ID.Set(1234),
		categories.Title.Set(convert.StringToText(txt)),
		pgtalk.NewTSVector(categories.TitleTokens, txt),
	)
	t.Log(pgtalk.SQL(mut))
	ctx := context.Background()
	tx, err := testConnect.Begin(ctx)
	if err != nil {
		t.Fatal(err)
	}
	it := mut.Exec(context.Background(), testConnect)
	if it.Err() != nil {
		t.Fatal(it.Err())
	}
	if err := tx.Commit(ctx); err != nil {
		t.Fatal(err)
	}
	// query
	q := categories.Select(categories.Columns()...).Where(pgtalk.NewTSQuery(categories.TitleTokens, "quick"))
	t.Log(pgtalk.SQL(q))
	list, err := q.Exec(context.Background(), testConnect)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(list)
}

func TestUnion(t *testing.T) {
	createCategory(t, 123)
	createCategory(t, 456)
	createProduct(t, 345, 123)
	createProduct(t, 678, 123)
	left := categories.Select(categories.ID, pgtalk.SQLAs("'category'", "type"))
	right := products.Select(products.ID, pgtalk.SQLAs("'product'", "type"))
	q := left.Union(right)
	list, err := q.ExecIntoMaps(context.Background(), testConnect)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(list)
}
