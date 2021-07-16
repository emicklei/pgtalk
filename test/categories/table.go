package categories

// generated by pgtalk-gen on 2021-07-16 12:02:38.673445 &#43;0200 CEST m=&#43;0.026874930
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
	tableInfo = pgtalk.TableInfo{Schema: "public", Name: "categories", Alias: "c1"}
)

type Category struct {
	ID    *int64  // bigint
	Title *string // text
}

var (
	ID = pgtalk.NewInt64Access(pgtalk.MakeColumnInfo(tableInfo, "id", true, true),
		func(dest interface{}, v *int64) { dest.(*Category).ID = v })
	Title = pgtalk.NewTextAccess(pgtalk.MakeColumnInfo(tableInfo, "title", false, false),
		func(dest interface{}, v *string) { dest.(*Category).Title = v })
)

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
func AllColumns() (all []pgtalk.ColumnAccessor) {
	return append(all, ID, Title)
}

// Select returns a new CategorysQuerySet for fetching column data.
func Select(cas ...pgtalk.ColumnAccessor) CategorysQuerySet {
	return CategorysQuerySet{pgtalk.MakeQuerySet(tableInfo, cas, func() interface{} {
		return new(Category)
	})}
}

// CategorysQuerySet can query for *Category values.
type CategorysQuerySet struct {
	pgtalk.QuerySet
}

func (s CategorysQuerySet) Unwrap() pgtalk.QuerySet { return s.QuerySet }

// Where returns a new QuerySet with WHERE clause.
func (s CategorysQuerySet) Where(condition pgtalk.SQLWriter) CategorysQuerySet {
	return CategorysQuerySet{QuerySet: s.QuerySet.Where(condition)}
}

// Limit returns a new QuerySet with the maximum number of results set.
func (s CategorysQuerySet) Limit(limit int) CategorysQuerySet {
	return CategorysQuerySet{QuerySet: s.QuerySet.Limit(limit)}
}

// GroupBy returns a new QuerySet with the GROUP BY clause.
func (s CategorysQuerySet) GroupBy(cas ...pgtalk.ColumnAccessor) CategorysQuerySet {
	return CategorysQuerySet{QuerySet: s.QuerySet.GroupBy(cas...)}
}

// OrderBy returns a new QuerySet with the ORDER BY clause.
func (s CategorysQuerySet) OrderBy(cas ...pgtalk.ColumnAccessor) CategorysQuerySet {
	return CategorysQuerySet{QuerySet: s.QuerySet.OrderBy(cas...)}
}

// Exec runs the query and returns the list of *Category.
func (s CategorysQuerySet) Exec(ctx context.Context, conn *pgx.Conn) (list []*Category, err error) {
	err = s.QuerySet.ExecWithAppender(ctx, conn, func(each interface{}) {
		list = append(list, each.(*Category))
	})
	return
}

// Insert creates a MutationSet for inserting data with zero or more columns.
func Insert(cas ...pgtalk.ColumnAccessor) pgtalk.MutationSet {
	return pgtalk.MakeMutationSet(tableInfo, cas, pgtalk.MutationInsert)
}

// Delete creates a MutationSet for deleting data.
func Delete() pgtalk.MutationSet {
	return pgtalk.MakeMutationSet(tableInfo, pgtalk.EmptyColumnAccessor, pgtalk.MutationDelete)
}

// Update creates a MutationSet to update zero or more columns.
func Update(cas ...pgtalk.ColumnAccessor) pgtalk.MutationSet {
	return pgtalk.MakeMutationSet(tableInfo, cas, pgtalk.MutationUpdate)
}
