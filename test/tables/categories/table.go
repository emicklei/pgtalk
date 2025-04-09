package categories
// generated by github.com/emicklei/pgtalk/cmd/pgtalk-gen from version: v1.11.1+dirty 
// DO NOT EDIT

import (
	p "github.com/emicklei/pgtalk"
	c "github.com/emicklei/pgtalk/convert"
	"time"
	"strings"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shopspring/decimal"
)

// Category is generated from the public.categories table.
type Category struct {
	ID	int32 // id : integer
	Title	pgtype.Text // title : character varying(255)
	// for storing custom field expression result values
	expressionResults map[string]any
}

var (	
	// ID represents the column "id" of with type "integer", nullable:false, primary:true
	ID = p.NewInt32Access(p.MakeColumnInfo(tableInfo, "id", p.IsPrimary, p.NotNull, 0),
		func(dest any) any { return &dest.(*Category).ID })	
	// Title represents the column "title" of with type "character varying(255)", nullable:true, primary:false
	Title = p.NewFieldAccess[pgtype.Text](p.MakeColumnInfo(tableInfo, "title", p.NotPrimary, p.Nullable, 0),
		func(dest any) any { return &dest.(*Category).Title })
	// unmapped fields
	//
	// TitleTokens holds information about the column "title_tokens" of with type "tsvector", nullable:true, primary:false
	TitleTokens = p.MakeColumnInfo(tableInfo, "title_tokens", p.NotPrimary, p.Nullable, 0)
	// package private
	_ = c.UUID // for the occasional unused import from convert
	_ = time.Now 
	_ = pgtype.Empty // for the occasional unused import from pgtype
	_ = decimal.Decimal{}
	tableInfo = p.TableInfo{Schema: "public", Name: "categories", Alias: "c1" }
)

func init() {
	// after var initialization (to prevent cycle) we need to update the tableInfo to set all columns
	tableInfo.Columns = []p.ColumnAccessor{ID,Title, }
}

// TableInfo returns meta information about the table.
func TableInfo() p.TableInfo {
	return tableInfo
}

// SetID sets the value to the field value and returns the receiver.
func (e *Category) SetID(v int32) *Category { e.ID = v ; return e }


// SetTitle sets the value to the field value and returns the receiver.
func (e *Category) SetTitle(v string) *Category { e.Title = c.StringToText(v) ; return e }

// Setters returns the list of changes to a Category for which updates/inserts need to be processed.
// Can be used in Insert,Update,Select. Cannot be used to set null (or empty array) values for columns.
func (e *Category) Setters() (list []p.ColumnAccessor) {
	list = append(list, ID.Set(e.ID))
	if e.Title.Valid {
		list = append(list, Title.Set(e.Title))
	}	
	return
}

// String returns the debug string for *Category with all non-nil field values.
func (e *Category) String() string {
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
func (e *Category) AddExpressionResult(key string, value any) {
	if e.expressionResults == nil {
		// lazy initialize
		e.expressionResults = map[string]any{}
	}
	e.expressionResults[key]=value
}

// GetExpressionResult gets a value from the custom expression results. Returns nil if absent.
func (e *Category) GetExpressionResult(key string) any {
	v, ok := e.expressionResults[key]
	if !ok {
		return nil
	}
	pv := v.(*any)
	return *pv
}

// Select returns a new QuerySet[Category] for fetching column data.
func Select(cas ...p.ColumnAccessor) p.QuerySet[Category] {
	return p.MakeQuerySet[Category](tableInfo, cas)
}

// Insert creates a MutationSet for inserting data with zero or more columns.
func Insert(cas ...p.ColumnAccessor) p.MutationSet[Category] {
	return p.MakeMutationSet[Category](tableInfo, cas, p.MutationInsert)
}

// Delete creates a MutationSet for deleting data.
func Delete() p.MutationSet[Category] {
	return p.MakeMutationSet[Category](tableInfo, p.EmptyColumnAccessor, p.MutationDelete)
}

// Update creates a MutationSet to update zero or more columns.
func Update(cas ...p.ColumnAccessor) p.MutationSet[Category] {
	return p.MakeMutationSet[Category](tableInfo, cas, p.MutationUpdate)
}
