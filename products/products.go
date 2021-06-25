package products

import (
	"github.com/emicklei/pgtalk/xs"
	"github.com/jackc/pgx/v4"
)

var tableInfo = xs.TableInfo{Name: "products", Alias: "t1"}

type Product struct {
	ID         *int64
	Code       *string
	Price      *int64
	CategoryID *int64
}

var (
	ID         = xs.NewInt8Access(tableInfo, "id", func(dest interface{}, i *int64) { dest.(*Product).ID = i })
	Code       = xs.NewTextAccess(tableInfo, "code", func(dest interface{}, i *string) { dest.(*Product).Code = i })
	CategoryID = xs.NewInt8Access(tableInfo, "category_id", func(dest interface{}, i *int64) { dest.(*Product).CategoryID = i })
	// or make this func?
	AllColumns = []xs.ReadWrite{ID, Code, CategoryID}
)

func Select(as ...xs.ReadWrite) ProductsQuerySet {
	return ProductsQuerySet{xs.MakeQuerySet(tableInfo, as, func() interface{} {
		return new(Product)
	})}
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
	err = s.QuerySet.Exec(conn, func(each interface{}) {
		list = append(list, each.(*Product))
	})
	return
}

func Insert(as ...xs.ReadWrite) xs.MutationSet {
	return xs.MakeMutationSet(tableInfo, as, xs.MutationInsert)
}

func Delete() xs.MutationSet {
	return xs.MakeMutationSet(tableInfo, xs.EmptyReadWrite, xs.MutationDelete)
}

func Update(as ...xs.ReadWrite) xs.MutationSet {
	return xs.MakeMutationSet(tableInfo, as, xs.MutationUpdate)
}
