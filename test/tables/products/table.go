package products

// generated by github.com/emicklei/pgtalk/cmd/pgtalk-gen from version: (devel)
// DO NOT EDIT

import (
	p "github.com/emicklei/pgtalk"
	c "github.com/emicklei/pgtalk/convert"
	"github.com/jackc/pgtype"
	"time"
)

// Product is generated from the public.products table.
type Product struct {
	ID         int64              // id : integer
	CreatedAt  pgtype.Timestamptz // created_at : timestamp with time zone
	UpdatedAt  pgtype.Timestamptz // updated_at : timestamp with time zone
	DeletedAt  pgtype.Timestamptz // deleted_at : timestamp with time zone
	Code       pgtype.Text        // code : text
	Price      pgtype.Int8        // price : bigint
	CategoryId pgtype.Int8        // category_id : bigint
}

var (
	// ID represents the column "id" of with type "integer", nullable:false, primary:true
	ID = p.NewInt64Access(p.MakeColumnInfo(tableInfo, "id", p.IsPrimary, p.NotNull, 1),
		func(dest any) any { return &dest.(*Product).ID })
	// CreatedAt represents the column "created_at" of with type "timestamp with time zone", nullable:true, primary:false
	CreatedAt = p.NewFieldAccess[pgtype.Timestamptz](p.MakeColumnInfo(tableInfo, "created_at", p.NotPrimary, p.Nullable, 2),
		func(dest any) any { return &dest.(*Product).CreatedAt })
	// UpdatedAt represents the column "updated_at" of with type "timestamp with time zone", nullable:true, primary:false
	UpdatedAt = p.NewFieldAccess[pgtype.Timestamptz](p.MakeColumnInfo(tableInfo, "updated_at", p.NotPrimary, p.Nullable, 3),
		func(dest any) any { return &dest.(*Product).UpdatedAt })
	// DeletedAt represents the column "deleted_at" of with type "timestamp with time zone", nullable:true, primary:false
	DeletedAt = p.NewFieldAccess[pgtype.Timestamptz](p.MakeColumnInfo(tableInfo, "deleted_at", p.NotPrimary, p.Nullable, 4),
		func(dest any) any { return &dest.(*Product).DeletedAt })
	// Code represents the column "code" of with type "text", nullable:true, primary:false
	Code = p.NewFieldAccess[pgtype.Text](p.MakeColumnInfo(tableInfo, "code", p.NotPrimary, p.Nullable, 5),
		func(dest any) any { return &dest.(*Product).Code })
	// Price represents the column "price" of with type "bigint", nullable:true, primary:false
	Price = p.NewFieldAccess[pgtype.Int8](p.MakeColumnInfo(tableInfo, "price", p.NotPrimary, p.Nullable, 6),
		func(dest any) any { return &dest.(*Product).Price })
	// CategoryId represents the column "category_id" of with type "bigint", nullable:true, primary:false
	CategoryId = p.NewFieldAccess[pgtype.Int8](p.MakeColumnInfo(tableInfo, "category_id", p.NotPrimary, p.Nullable, 7),
		func(dest any) any { return &dest.(*Product).CategoryId })
	// package private
	_         = c.UUID // for the occasional unused import from convert
	_         = time.Now
	_         = pgtype.Empty // for the occasional unused import from pgtype
	tableInfo = p.TableInfo{Schema: "public", Name: "products", Alias: "p1"}
)

func init() {
	// after var initialization (to prevent cycle) we need to update the tableInfo to set all columns
	tableInfo.Columns = []p.ColumnAccessor{ID, CreatedAt, UpdatedAt, DeletedAt, Code, Price, CategoryId}
}

// SetID sets the value to the field value and returns the receiver.
func (e *Product) SetID(v int64) *Product { e.ID = v; return e }

// SetCreatedAt sets the value to the field value and returns the receiver.
func (e *Product) SetCreatedAt(v time.Time) *Product { e.CreatedAt = c.TimeToTimestamptz(v); return e }

// SetUpdatedAt sets the value to the field value and returns the receiver.
func (e *Product) SetUpdatedAt(v time.Time) *Product { e.UpdatedAt = c.TimeToTimestamptz(v); return e }

// SetDeletedAt sets the value to the field value and returns the receiver.
func (e *Product) SetDeletedAt(v time.Time) *Product { e.DeletedAt = c.TimeToTimestamptz(v); return e }

// SetCode sets the value to the field value and returns the receiver.
func (e *Product) SetCode(v string) *Product { e.Code = c.StringToText(v); return e }

// SetPrice sets the value to the field value and returns the receiver.
func (e *Product) SetPrice(v int64) *Product { e.Price = c.Int64ToInt8(v); return e }

// SetCategoryId sets the value to the field value and returns the receiver.
func (e *Product) SetCategoryId(v int64) *Product { e.CategoryId = c.Int64ToInt8(v); return e }

// Setters returns the list of changes to a Product for which updates/inserts need to be processed.
// Can be used in Insert,Update,Select. Cannot be used to set null values for columns.
func (e *Product) Setters() (list []p.ColumnAccessor) {
	list = append(list, ID.Set(e.ID))
	if e.CreatedAt.Status == pgtype.Present {
		list = append(list, CreatedAt.Set(e.CreatedAt))
	}
	if e.UpdatedAt.Status == pgtype.Present {
		list = append(list, UpdatedAt.Set(e.UpdatedAt))
	}
	if e.DeletedAt.Status == pgtype.Present {
		list = append(list, DeletedAt.Set(e.DeletedAt))
	}
	if e.Code.Status == pgtype.Present {
		list = append(list, Code.Set(e.Code))
	}
	if e.Price.Status == pgtype.Present {
		list = append(list, Price.Set(e.Price))
	}
	if e.CategoryId.Status == pgtype.Present {
		list = append(list, CategoryId.Set(e.CategoryId))
	}
	return
}

// String returns the debug string for *Product with all non-nil field values.
func (e *Product) String() string {
	return p.StringWithFields(e, p.HideNilValues)
}

// Columns returns the ColumnAccessor list for the given column names.
// If the names is empty then return all columns.
func Columns(names ...string) (list []p.ColumnAccessor) {
	if len(names) == 0 {
		return tableInfo.Columns
	}
	for _, each := range names {
		for _, other := range tableInfo.Columns {
			if other.Column().Name() == each {
				list = append(list, other)
			}
		}
	}
	return
}

// Select returns a new QuerySet[Product] for fetching column data.
func Select(cas ...p.ColumnAccessor) p.QuerySet[Product] {
	return p.MakeQuerySet[Product](tableInfo, cas)
}

// Insert creates a MutationSet for inserting data with zero or more columns.
func Insert(cas ...p.ColumnAccessor) p.MutationSet[Product] {
	return p.MakeMutationSet[Product](tableInfo, cas, p.MutationInsert)
}

// Delete creates a MutationSet for deleting data.
func Delete() p.MutationSet[Product] {
	return p.MakeMutationSet[Product](tableInfo, p.EmptyColumnAccessor, p.MutationDelete)
}

// Update creates a MutationSet to update zero or more columns.
func Update(cas ...p.ColumnAccessor) p.MutationSet[Product] {
	return p.MakeMutationSet[Product](tableInfo, cas, p.MutationUpdate)
}
