package categories

// generated by pgtalk-gen on Wed, 05 Jan 2022 10:22:40 CET
// DO NOT EDIT

import (
	p "github.com/emicklei/pgtalk"
	"github.com/jackc/pgtype"
	"time"
)

// Category is generated from the public.categories table.
type Category struct {
	ID    int64       // id : integer
	Title pgtype.Text // title : text
}

var (
	// ID represents the column "id" of with type "integer", nullable:false, primary:true
	ID = p.NewInt64Access(p.MakeColumnInfo(tableInfo, "id", p.IsPrimary, p.NotNull, 1),
		func(dest interface{}, v int64) { dest.(*Category).ID = v }, nil)
	// Title represents the column "title" of with type "text", nullable:true, primary:false
	Title = p.NewTextAccess(p.MakeColumnInfo(tableInfo, "title", p.NotPrimary, p.Nullable, 2),
		nil, func(dest interface{}, v pgtype.Text) { dest.(*Category).Title = v })
	// package private
	_         = time.Now
	_         = pgtype.Empty // for the occasional unused import from pgtype
	tableInfo = p.TableInfo{Schema: "public", Name: "categories", Alias: "c1"}
)

func init() {
	// after var initialization (to prevent cycle) we need to update the tableInfo to set all columns
	tableInfo.Columns = []p.ColumnAccessor{ID, Title}
}

// SetID sets the value to the field value and returns the receiver.
func (e *Category) SetID(v int64) *Category { e.ID = v; return e }

// SetTitle sets the value to the field value and returns the receiver.
func (e *Category) SetTitle(v pgtype.Text) *Category { e.Title = v; return e }

// Setters returns the list of changes to a Category for which updates/inserts need to be processed.
// Can be used in Insert,Update,Select. Cannot be used to set null values for columns.
func (e *Category) Setters() (list []p.ColumnAccessor) {
	list = append(list, ID.Set(e.ID))
	if e.Title.Status == pgtype.Present {
		list = append(list, Title.Set(e.Title.String))
	}
	return
}

// String returns the debug string for *Category with all non-nil field values.
func (e *Category) String() string {
	return p.StringWithFields(e, p.HideNilValues)
}

// AllColumns returns the list of all column accessors for usage in e.g. Select.
func AllColumns() []p.ColumnAccessor {
	return tableInfo.Columns
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
