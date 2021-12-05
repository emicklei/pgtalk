package categories
// generated by pgtalk-gen on Sun, 05 Dec 2021 14:02:25 CET 
// DO NOT EDIT

import (
	p "github.com/emicklei/pgtalk"
	"time"
	"fmt"
	"bytes"
	"github.com/jackc/pgtype"
)

// Category is generated from the public.categories table.
type Category struct {
	ID	*int64 // id : integer
	Title	*string // title : text
}

var (	
	ID = p.NewInt64Access(p.MakeColumnInfo(tableInfo, "id", p.IsPrimary, p.NotNull, 1),
		func(dest interface{}, v *int64) { dest.(*Category).ID = v })	
	Title = p.NewTextAccess(p.MakeColumnInfo(tableInfo, "title", p.NotPrimary, p.Nullable, 2),
		func(dest interface{}, v *string) { dest.(*Category).Title = v })
	// package private
	_ = time.Now()
	_ = pgtype.Empty // for the occasional unused import
	tableInfo = p.TableInfo{Schema: "public", Name: "categories", Alias: "c1"}
	tableAccess = p.TableAccessor{TableInfo: tableInfo, 
		Factory: func() interface{}{return new(Category)}, AllColumns: []p.ColumnAccessor{ID,Title,
}}
)

// ColumnUpdatesFrom returns the list of changes to a Category for which updates need to be processed.
// Cannot be used to set null values for columns.
func ColumnUpdatesFrom(e *Category) (list []p.SQLExpression) {
	if e.ID != nil {
		list = append(list, ID.Set(*e.ID))
	}
	if e.Title != nil {
		list = append(list, Title.Set(*e.Title))
	}	
	return
}

// String returns the debug string for *Category with all non-nil field values.
func (e *Category) String() string {
	b := new(bytes.Buffer)
	fmt.Fprint(b, "categories.Category{")
	if e.ID != nil {
		fmt.Fprintf(b, "ID:%v ", *e.ID)
	}
	if e.Title != nil {
		fmt.Fprintf(b, "Title:%v ", *e.Title)
	}
	fmt.Fprint(b, "}")
	return b.String()
}

// AllColumns returns the list of all column accessors for usage in e.g. Select.
func AllColumns() []p.ColumnAccessor {
	return tableAccess.AllColumns
}

// Select returns a new QuerySet[Category] for fetching column data.
func Select(cas ...p.ColumnAccessor) p.QuerySet[Category] {
	return p.MakeQuerySet[Category](tableAccess, cas)
}

// Insert creates a MutationSet for inserting data with zero or more columns.
func Insert(cas ...p.ColumnAccessor) p.MutationSet[Category] {
	return p.MakeMutationSet[Category](tableAccess, cas, p.MutationInsert)
}

// Delete creates a MutationSet for deleting data.
func Delete() p.MutationSet[Category] {
	return p.MakeMutationSet[Category](tableAccess, p.EmptyColumnAccessor, p.MutationDelete)
}

// Update creates a MutationSet to update zero or more columns.
func Update(cas ...p.ColumnAccessor) p.MutationSet[Category] {
	return p.MakeMutationSet[Category](tableAccess, cas, p.MutationUpdate)
}
