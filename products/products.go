package products

import (
	"github.com/emicklei/pgtalk/xs"
	"github.com/jackc/pgx/v4"
)

type Product struct {
	ID         *int64
	Code       *string
	Price      *int64
	CategoryID *int64
}

var ID = xs.NewInt8Access(
	"products",
	"id",
	func(dest interface{}, i *int64) {
		e := dest.(*Product)
		e.ID = i
	})

var Code = xs.NewTextAccess(
	"products",
	"code",
	func(dest interface{}, i *string) {
		e := dest.(*Product)
		e.Code = i
	})

var CategoryID = xs.NewInt8Access(
	"products",
	"categoryID",
	func(dest interface{}, i *int64) {
		e := dest.(*Product)
		e.CategoryID = i
	})

func Select(as ...xs.ReadWrite) ProductsQuerySet {
	return ProductsQuerySet{xs.MakeQuerySet("products", as)}
}

type ProductsQuerySet struct {
	xs.QuerySet
}

func (s ProductsQuerySet) Unwrap() xs.QuerySet { return s.QuerySet }

// Where is
func (s ProductsQuerySet) Where(condition xs.SQLWriter) ProductsQuerySet {
	return ProductsQuerySet{QuerySet: s.QuerySet.Where(condition)}
}

// Limit is
func (s ProductsQuerySet) Limit(limit int) ProductsQuerySet {
	return ProductsQuerySet{QuerySet: s.QuerySet.Limit(limit)}
}

// Exec is
func (s ProductsQuerySet) Exec(conn *pgx.Conn) (list []*Product, err error) {
	err = s.QuerySet.Exec(conn, func() interface{} {
		return new(Product)
	}, func(each interface{}) {
		list = append(list, each.(*Product))
	})
	return
}
