package products
// generated by github.com/emicklei/pgtalk/cmd/pgtalk-gen from version: (devel) 
// DO NOT EDIT

import (
	p "github.com/emicklei/pgtalk"
	c "github.com/emicklei/pgtalk/convert"
	"time"
	"strings"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shopspring/decimal" 
)

// Product is generated from the public.products table.
type Product struct {
	CategoryId	pgtype.Int4 // category_id : integer
	Code	pgtype.Text // code : text
	CreatedAt	pgtype.Timestamptz // created_at : timestamp with time zone
	DeletedAt	pgtype.Timestamptz // deleted_at : timestamp with time zone
	ID	int32 // id : integer
	Price	pgtype.Int8 // price : bigint
	UpdatedAt	pgtype.Timestamptz // updated_at : timestamp with time zone
	// for storing custom field expression result values
	expressionResults map[string]any
}

var (	
	// CategoryId represents the column "category_id" of with type "integer", nullable:true, primary:false
	CategoryId = p.NewFieldAccess[pgtype.Int4](p.MakeColumnInfo(tableInfo, "category_id", p.NotPrimary, p.Nullable, 0),
		func(dest any) any { return &dest.(*Product).CategoryId })	
	// Code represents the column "code" of with type "text", nullable:true, primary:false
	Code = p.NewFieldAccess[pgtype.Text](p.MakeColumnInfo(tableInfo, "code", p.NotPrimary, p.Nullable, 0),
		func(dest any) any { return &dest.(*Product).Code })	
	// CreatedAt represents the column "created_at" of with type "timestamp with time zone", nullable:true, primary:false
	CreatedAt = p.NewFieldAccess[pgtype.Timestamptz](p.MakeColumnInfo(tableInfo, "created_at", p.NotPrimary, p.Nullable, 0),
		func(dest any) any { return &dest.(*Product).CreatedAt })	
	// DeletedAt represents the column "deleted_at" of with type "timestamp with time zone", nullable:true, primary:false
	DeletedAt = p.NewFieldAccess[pgtype.Timestamptz](p.MakeColumnInfo(tableInfo, "deleted_at", p.NotPrimary, p.Nullable, 0),
		func(dest any) any { return &dest.(*Product).DeletedAt })	
	// ID represents the column "id" of with type "integer", nullable:false, primary:true
	ID = p.NewInt32Access(p.MakeColumnInfo(tableInfo, "id", p.IsPrimary, p.NotNull, 0),
		func(dest any) any { return &dest.(*Product).ID })	
	// Price represents the column "price" of with type "bigint", nullable:true, primary:false
	Price = p.NewFieldAccess[pgtype.Int8](p.MakeColumnInfo(tableInfo, "price", p.NotPrimary, p.Nullable, 0),
		func(dest any) any { return &dest.(*Product).Price })	
	// UpdatedAt represents the column "updated_at" of with type "timestamp with time zone", nullable:true, primary:false
	UpdatedAt = p.NewFieldAccess[pgtype.Timestamptz](p.MakeColumnInfo(tableInfo, "updated_at", p.NotPrimary, p.Nullable, 0),
		func(dest any) any { return &dest.(*Product).UpdatedAt })
	// package private
	_ = c.UUID // for the occasional unused import from convert
	_ = time.Now 
	_ = pgtype.Empty // for the occasional unused import from pgtype
	_ = decimal.Decimal{}
	tableInfo = p.TableInfo{Schema: "public", Name: "products", Alias: "p1" }
)

func init() {
	// after var initialization (to prevent cycle) we need to update the tableInfo to set all columns
	tableInfo.Columns = []p.ColumnAccessor{CategoryId,Code,CreatedAt,DeletedAt,ID,Price,UpdatedAt, }
}

// TableInfo returns meta information about the table.
func TableInfo() p.TableInfo {
	return tableInfo
}


// SetCategoryId sets the value to the field value and returns the receiver.
func (e *Product) SetCategoryId(v int32) *Product { e.CategoryId = c.Int32ToInt4(v) ; return e }


// SetCode sets the value to the field value and returns the receiver.
func (e *Product) SetCode(v string) *Product { e.Code = c.StringToText(v) ; return e }


// SetCreatedAt sets the value to the field value and returns the receiver.
func (e *Product) SetCreatedAt(v time.Time) *Product { e.CreatedAt = c.TimeToTimestamptz(v) ; return e }


// SetDeletedAt sets the value to the field value and returns the receiver.
func (e *Product) SetDeletedAt(v time.Time) *Product { e.DeletedAt = c.TimeToTimestamptz(v) ; return e }

// SetID sets the value to the field value and returns the receiver.
func (e *Product) SetID(v int32) *Product { e.ID = v ; return e }


// SetPrice sets the value to the field value and returns the receiver.
func (e *Product) SetPrice(v int64) *Product { e.Price = c.Int64ToInt8(v) ; return e }


// SetUpdatedAt sets the value to the field value and returns the receiver.
func (e *Product) SetUpdatedAt(v time.Time) *Product { e.UpdatedAt = c.TimeToTimestamptz(v) ; return e }

// Setters returns the list of changes to a Product for which updates/inserts need to be processed.
// Can be used in Insert,Update,Select. Cannot be used to set null (or empty array) values for columns.
func (e *Product) Setters() (list []p.ColumnAccessor) {
	if e.CategoryId.Valid {
		list = append(list, CategoryId.Set(e.CategoryId))
	}
	if e.Code.Valid {
		list = append(list, Code.Set(e.Code))
	}
	if e.CreatedAt.Valid {
		list = append(list, CreatedAt.Set(e.CreatedAt))
	}
	if e.DeletedAt.Valid {
		list = append(list, DeletedAt.Set(e.DeletedAt))
	}
	list = append(list, ID.Set(e.ID))
	if e.Price.Valid {
		list = append(list, Price.Set(e.Price))
	}
	if e.UpdatedAt.Valid {
		list = append(list, UpdatedAt.Set(e.UpdatedAt))
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
			n := other.Column().Name()
			if strings.HasPrefix(n,"'") { // mixed case names are quoted
				n = strings.Trim(n,"'")
			} 
			if n == each {
				list = append(list, other)
			}
		}
	}
	return
}

// AddExpressionResult puts a value into the custom expression results
func (e *Product) AddExpressionResult(key string, value any) {
	if e.expressionResults == nil {
		// lazy initialize
		e.expressionResults = map[string]any{}
	}
	e.expressionResults[key]=value
}

// GetExpressionResult gets a value from the custom expression results. Returns nil if absent.
func (e *Product) GetExpressionResult(key string) any {
	v, ok := e.expressionResults[key]
	if !ok {
		return nil
	}
	pv := v.(*any)
	return *pv
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
