package products

// generated by pgtalk-gen on 2021-08-10 11:59:41.195569 &#43;0200 CEST m=&#43;0.024146111
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
	ID = pgtalk.NewInt64Access(pgtalk.MakeColumnInfo(tableInfo, "id", true, true),
		func(dest interface{}, v *int64) { dest.(*Product).ID = v })
	Created_at = pgtalk.NewTimeAccess(pgtalk.MakeColumnInfo(tableInfo, "created_at", false, false),
		func(dest interface{}, v *time.Time) { dest.(*Product).Created_at = v })
	Updated_at = pgtalk.NewTimeAccess(pgtalk.MakeColumnInfo(tableInfo, "updated_at", false, false),
		func(dest interface{}, v *time.Time) { dest.(*Product).Updated_at = v })
	Deleted_at = pgtalk.NewTimeAccess(pgtalk.MakeColumnInfo(tableInfo, "deleted_at", false, false),
		func(dest interface{}, v *time.Time) { dest.(*Product).Deleted_at = v })
	Code = pgtalk.NewTextAccess(pgtalk.MakeColumnInfo(tableInfo, "code", false, false),
		func(dest interface{}, v *string) { dest.(*Product).Code = v })
	Price = pgtalk.NewInt64Access(pgtalk.MakeColumnInfo(tableInfo, "price", false, false),
		func(dest interface{}, v *int64) { dest.(*Product).Price = v })
	Category_id = pgtalk.NewInt64Access(pgtalk.MakeColumnInfo(tableInfo, "category_id", false, false),
		func(dest interface{}, v *int64) { dest.(*Product).Category_id = v })
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

// setColumnValueTo sets the field of a *Product to the non-nil value.
func setColumnValueTo(e *Product, tableAttributeNumber uint16, value interface{}) error {
	if value == nil {
		return nil
	}
	switch tableAttributeNumber {
	case 1:
		if tvalue, ok := value.(int64); !ok {
			return fmt.Errorf("unable to assert value of type [*%T] to field [ID] of type [%s] in entity with type [%s]", value, "*int64", "*products.Product")
		} else {
			e.ID = &tvalue
		}
	case 2:
		if tvalue, ok := value.(time.Time); !ok {
			return fmt.Errorf("unable to assert value of type [*%T] to field [Created_at] of type [%s] in entity with type [%s]", value, "*time.Time", "*products.Product")
		} else {
			e.Created_at = &tvalue
		}
	case 3:
		if tvalue, ok := value.(time.Time); !ok {
			return fmt.Errorf("unable to assert value of type [*%T] to field [Updated_at] of type [%s] in entity with type [%s]", value, "*time.Time", "*products.Product")
		} else {
			e.Updated_at = &tvalue
		}
	case 4:
		if tvalue, ok := value.(time.Time); !ok {
			return fmt.Errorf("unable to assert value of type [*%T] to field [Deleted_at] of type [%s] in entity with type [%s]", value, "*time.Time", "*products.Product")
		} else {
			e.Deleted_at = &tvalue
		}
	case 5:
		if tvalue, ok := value.(string); !ok {
			return fmt.Errorf("unable to assert value of type [*%T] to field [Code] of type [%s] in entity with type [%s]", value, "*string", "*products.Product")
		} else {
			e.Code = &tvalue
		}
	case 6:
		if tvalue, ok := value.(int64); !ok {
			return fmt.Errorf("unable to assert value of type [*%T] to field [Price] of type [%s] in entity with type [%s]", value, "*int64", "*products.Product")
		} else {
			e.Price = &tvalue
		}
	case 7:
		if tvalue, ok := value.(int64); !ok {
			return fmt.Errorf("unable to assert value of type [*%T] to field [Category_id] of type [%s] in entity with type [%s]", value, "*int64", "*products.Product")
		} else {
			e.Category_id = &tvalue
		}
	default:
		return fmt.Errorf("unable to set value [%v] to field of [%v] with table attribute number [%d] in entity with type [%s]", value, e, tableAttributeNumber, "*products.Product")
	}
	return nil
}

var fieldSetter = func(entityPointer interface{}, tableAttributeNumber uint16, value interface{}) error {
	return setColumnValueTo(entityPointer.(*Product), tableAttributeNumber, value)
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
func AllColumns() (all []pgtalk.ColumnAccessor) {
	return append(all, ID, Created_at, Updated_at, Deleted_at, Code, Price, Category_id)
}

// Select returns a new ProductsQuerySet for fetching column data.
func Select(cas ...pgtalk.ColumnAccessor) ProductsQuerySet {
	return ProductsQuerySet{pgtalk.MakeQuerySet(tableInfo, cas, func() interface{} {
		return new(Product)
	}, fieldSetter)}
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
	return pgtalk.MakeMutationSet(tableInfo, cas, pgtalk.MutationInsert, fieldSetter)
}

// Delete creates a MutationSet for deleting data.
func Delete() pgtalk.MutationSet {
	return pgtalk.MakeMutationSet(tableInfo, pgtalk.EmptyColumnAccessor, pgtalk.MutationDelete, fieldSetter)
}

// Update creates a MutationSet to update zero or more columns.
func Update(cas ...pgtalk.ColumnAccessor) pgtalk.MutationSet {
	return pgtalk.MakeMutationSet(tableInfo, cas, pgtalk.MutationUpdate, fieldSetter)
}
