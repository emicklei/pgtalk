package products

import (
	"github.com/emicklei/pgtalk"
)

var tableInfo = pgtalk.TableInfo{Name: "products", Alias: "t1"}

type Product struct {
	ID         *int64
	Code       *string
	Price      *int64
	CategoryID *int64
}

var (
	ID         = pgtalk.NewInt8Access(tableInfo, "id", func(dest interface{}, i *int64) { dest.(*Product).ID = i })
	Code       = pgtalk.NewTextAccess(tableInfo, "code", func(dest interface{}, i *string) { dest.(*Product).Code = i })
	CategoryID = pgtalk.NewInt8Access(tableInfo, "category_id", func(dest interface{}, i *int64) { dest.(*Product).CategoryID = i })
	// or make this func?
	AllColumns = []pgtalk.ColumnAccessor{ID, Code, CategoryID}
)

func Select(as ...pgtalk.ColumnAccessor) ProductsQuerySet {
	return ProductsQuerySet{pgtalk.MakeQuerySet(tableInfo, as, func() interface{} {
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

func Insert(as ...pgtalk.ColumnAccessor) pgtalk.MutationSet {
	return pgtalk.MakeMutationSet(tableInfo, as, pgtalk.MutationInsert)
}

func Delete() pgtalk.MutationSet {
	return pgtalk.MakeMutationSet(tableInfo, pgtalk.EmptyColumnAccessor, pgtalk.MutationDelete)
}

func Update(as ...pgtalk.ColumnAccessor) pgtalk.MutationSet {
	return pgtalk.MakeMutationSet(tableInfo, as, pgtalk.MutationUpdate)
}
