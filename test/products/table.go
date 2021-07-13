package products

// generated by pgtalk-gen on 2021-07-05 09:56:03.800909 &#43;0200 CEST m=&#43;0.025147062 DO NOT EDIT

import (
	"github.com/emicklei/pgtalk"
	"time"
)

var (
	_         = time.Now()
	tableInfo = pgtalk.TableInfo{Name: "products", Alias: "p1"}
)

type Product struct {
	ID          *int64
	Created_at  *time.Time
	Updated_at  *time.Time
	Deleted_at  *time.Time
	Code        *string
	Price       *int64
	Category_id *int64
}

var (
	/**
		ID = pgtalk.NewInt64Access(tableInfo, "id", func(dest interface{}, v *int64) error {
		if e, ok := dest.(*Product); ok { e.ID = v; return nil }
		return pgtalk.EntityTypeError(dest, "products.Product")
	})
	**/
	ID          = pgtalk.NewInt64Access(tableInfo, "id", func(dest interface{}, v *int64) { dest.(*Product).ID = v })
	Created_at  = pgtalk.NewTimeAccess(tableInfo, "created_at", func(dest interface{}, v *time.Time) { dest.(*Product).Created_at = v })
	Updated_at  = pgtalk.NewTimeAccess(tableInfo, "updated_at", func(dest interface{}, v *time.Time) { dest.(*Product).Updated_at = v })
	Deleted_at  = pgtalk.NewTimeAccess(tableInfo, "deleted_at", func(dest interface{}, v *time.Time) { dest.(*Product).Deleted_at = v })
	Code        = pgtalk.NewTextAccess(tableInfo, "code", func(dest interface{}, v *string) { dest.(*Product).Code = v })
	Price       = pgtalk.NewInt64Access(tableInfo, "price", func(dest interface{}, v *int64) { dest.(*Product).Price = v })
	Category_id = pgtalk.NewInt64Access(tableInfo, "category_id", func(dest interface{}, v *int64) { dest.(*Product).Category_id = v })
)

func AllColumns() (all []pgtalk.ColumnAccessor) {
	return append(all, ID, Created_at, Updated_at, Deleted_at, Code, Price, Category_id)
}

func Select(cas ...pgtalk.ColumnAccessor) ProductsQuerySet {
	return ProductsQuerySet{pgtalk.MakeQuerySet(tableInfo, cas, func() interface{} {
		return new(Product)
	})}
}

type ProductsQuerySet struct {
	pgtalk.QuerySet
}

func (s ProductsQuerySet) Unwrap() pgtalk.QuerySet { return s.QuerySet }

// Where is
func (s ProductsQuerySet) Where(condition pgtalk.SQLWriter) ProductsQuerySet {
	return ProductsQuerySet{QuerySet: s.QuerySet.Where(condition)}
}

// Limit is
func (s ProductsQuerySet) Limit(limit int) ProductsQuerySet {
	return ProductsQuerySet{QuerySet: s.QuerySet.Limit(limit)}
}

// GroupBy is
func (s ProductsQuerySet) GroupBy(cas ...pgtalk.ColumnAccessor) ProductsQuerySet {
	return ProductsQuerySet{QuerySet: s.QuerySet.GroupBy(cas...)}
}

// GroupBy is
func (s ProductsQuerySet) OrderBy(cas ...pgtalk.ColumnAccessor) ProductsQuerySet {
	return ProductsQuerySet{QuerySet: s.QuerySet.OrderBy(cas...)}
}

// Exec is
func (s ProductsQuerySet) Exec(conn pgtalk.Connection) (list []*Product, err error) {
	err = s.QuerySet.ExecWithAppender(conn, func(each interface{}) {
		list = append(list, each.(*Product))
	})
	return
}

func Insert(cas ...pgtalk.ColumnAccessor) pgtalk.MutationSet {
	return pgtalk.MakeMutationSet(tableInfo, cas, pgtalk.MutationInsert)
}

func Delete() pgtalk.MutationSet {
	return pgtalk.MakeMutationSet(tableInfo, pgtalk.EmptyColumnAccessor, pgtalk.MutationDelete)
}

func Update(cas ...pgtalk.ColumnAccessor) pgtalk.MutationSet {
	return pgtalk.MakeMutationSet(tableInfo, cas, pgtalk.MutationUpdate)
}
