package products

// generated by pgtalk-gen on 2021-08-17 10:36:41.845302 &#43;0200 CEST m=&#43;0.024933820
// DO NOT EDIT

import (
	"bytes"
	"context"
	"fmt"
	"github.com/emicklei/pgtalk"
	"github.com/jackc/pgx/v4"
	"time"
)

var (
	_         = time.Now()
	tableInfo = pgtalk.TableInfo{Schema: "public", Name: "products", Alias: "p1"}
)

type Product struct {
	ID          *int64     // bigint
	Created_at  *time.Time // timestamp with time zone
	Updated_at  *time.Time // timestamp with time zone
	Deleted_at  *time.Time // timestamp with time zone
	Code        *string    // text
	Price       *int64     // bigint
	Category_id *int64     // bigint
}

var (
	ID = pgtalk.NewInt64Access(pgtalk.MakeColumnInfo(tableInfo, "id", true, true, 1),
		func(dest interface{}, v *int64) { dest.(*Product).ID = v })
	Created_at = pgtalk.NewTimeAccess(pgtalk.MakeColumnInfo(tableInfo, "created_at", false, false, 2),
		func(dest interface{}, v *time.Time) { dest.(*Product).Created_at = v })
	Updated_at = pgtalk.NewTimeAccess(pgtalk.MakeColumnInfo(tableInfo, "updated_at", false, false, 3),
		func(dest interface{}, v *time.Time) { dest.(*Product).Updated_at = v })
	Deleted_at = pgtalk.NewTimeAccess(pgtalk.MakeColumnInfo(tableInfo, "deleted_at", false, false, 4),
		func(dest interface{}, v *time.Time) { dest.(*Product).Deleted_at = v })
	Code = pgtalk.NewTextAccess(pgtalk.MakeColumnInfo(tableInfo, "code", false, false, 5),
		func(dest interface{}, v *string) { dest.(*Product).Code = v })
	Price = pgtalk.NewInt64Access(pgtalk.MakeColumnInfo(tableInfo, "price", false, false, 6),
		func(dest interface{}, v *int64) { dest.(*Product).Price = v })
	Category_id = pgtalk.NewInt64Access(pgtalk.MakeColumnInfo(tableInfo, "category_id", false, false, 7),
		func(dest interface{}, v *int64) { dest.(*Product).Category_id = v })
	tableAccess = pgtalk.TableAccessor{TableInfo: tableInfo, AllColumns: []pgtalk.ColumnAccessor{ID, Created_at, Updated_at, Deleted_at, Code, Price, Category_id}}
)

// ColumnUpdatesFrom returns the list of changes to a Product for which updates need to be processed.
// Cannot be used to set null values for columns.
func ColumnUpdatesFrom(e Product) (list []pgtalk.SQLWriter) {
	if e.ID != nil {
		list = append(list, ID.Set(*e.ID))
	}
	if e.Created_at != nil {
		list = append(list, Created_at.Set(*e.Created_at))
	}
	if e.Updated_at != nil {
		list = append(list, Updated_at.Set(*e.Updated_at))
	}
	if e.Deleted_at != nil {
		list = append(list, Deleted_at.Set(*e.Deleted_at))
	}
	if e.Code != nil {
		list = append(list, Code.Set(*e.Code))
	}
	if e.Price != nil {
		list = append(list, Price.Set(*e.Price))
	}
	if e.Category_id != nil {
		list = append(list, Category_id.Set(*e.Category_id))
	}
	return
}

// Next returns the next *Product from the iterator data.
// Use err to check for failure.
func Next(it *pgtalk.ResultIterator) (e *Product, err error) {
	var each = new(Product)
	// first check for query error in case caller forgot
	if err = it.Err(); err != nil {
		return nil, err
	}
	err = it.Next(each)
	return each, err
}

// String returns the debug string for *Product with all non-nil field values.
func (e *Product) String() string {
	b := new(bytes.Buffer)
	fmt.Fprint(b, "products.Product{")
	if e.ID != nil {
		fmt.Fprintf(b, "ID:%v ", *e.ID)
	}
	if e.Created_at != nil {
		fmt.Fprintf(b, "Created_at:%v ", *e.Created_at)
	}
	if e.Updated_at != nil {
		fmt.Fprintf(b, "Updated_at:%v ", *e.Updated_at)
	}
	if e.Deleted_at != nil {
		fmt.Fprintf(b, "Deleted_at:%v ", *e.Deleted_at)
	}
	if e.Code != nil {
		fmt.Fprintf(b, "Code:%v ", *e.Code)
	}
	if e.Price != nil {
		fmt.Fprintf(b, "Price:%v ", *e.Price)
	}
	if e.Category_id != nil {
		fmt.Fprintf(b, "Category_id:%v ", *e.Category_id)
	}
	fmt.Fprint(b, "}")
	return b.String()
}

// AllColumns returns the list of all column accessors for usage in e.g. Select.
func AllColumns() []pgtalk.ColumnAccessor {
	return tableAccess.AllColumns
}

// Select returns a new ProductsQuerySet for fetching column data.
func Select(cas ...pgtalk.ColumnAccessor) ProductsQuerySet {
	return ProductsQuerySet{pgtalk.MakeQuerySet(tableAccess, cas, func() interface{} {
		return new(Product)
	})}
}

// ProductsQuerySet can query for *Product values.
type ProductsQuerySet struct {
	pgtalk.QuerySet
}

func (s ProductsQuerySet) Unwrap() pgtalk.QuerySet { return s.QuerySet }

// Where returns a new QuerySet with WHERE clause.
func (s ProductsQuerySet) Where(condition pgtalk.SQLWriter) ProductsQuerySet {
	return ProductsQuerySet{QuerySet: s.QuerySet.Where(condition)}
}

// Limit returns a new QuerySet with the maximum number of results set.
func (s ProductsQuerySet) Limit(limit int) ProductsQuerySet {
	return ProductsQuerySet{QuerySet: s.QuerySet.Limit(limit)}
}

// GroupBy returns a new QuerySet with the GROUP BY clause.
func (s ProductsQuerySet) GroupBy(cas ...pgtalk.ColumnAccessor) ProductsQuerySet {
	return ProductsQuerySet{QuerySet: s.QuerySet.GroupBy(cas...)}
}

// OrderBy returns a new QuerySet with the ORDER BY clause.
func (s ProductsQuerySet) OrderBy(cas ...pgtalk.ColumnAccessor) ProductsQuerySet {
	return ProductsQuerySet{QuerySet: s.QuerySet.OrderBy(cas...)}
}

// Exec runs the query and returns the list of *Product.
func (s ProductsQuerySet) Exec(ctx context.Context, conn *pgx.Conn) (list []*Product, err error) {
	err = s.QuerySet.ExecWithAppender(ctx, conn, func(each interface{}) {
		list = append(list, each.(*Product))
	})
	return
}

// Insert creates a MutationSet for inserting data with zero or more columns.
func Insert(cas ...pgtalk.ColumnAccessor) pgtalk.MutationSet {
	return pgtalk.MakeMutationSet(tableAccess, cas, pgtalk.MutationInsert)
}

// Delete creates a MutationSet for deleting data.
func Delete() pgtalk.MutationSet {
	return pgtalk.MakeMutationSet(tableAccess, pgtalk.EmptyColumnAccessor, pgtalk.MutationDelete)
}

// Update creates a MutationSet to update zero or more columns.
func Update(cas ...pgtalk.ColumnAccessor) pgtalk.MutationSet {
	return pgtalk.MakeMutationSet(tableAccess, cas, pgtalk.MutationUpdate)
}
