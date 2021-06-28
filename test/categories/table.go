package categories

import (
	"github.com/emicklei/pgtalk"
)

var tableInfo = pgtalk.TableInfo{Name: "categories", Alias: "t2"}

type Category struct {
	ID    *int64
	Title *string
}

var ID = pgtalk.NewInt8Access(
	tableInfo,
	"id",
	func(dest interface{}, i *int64) {
		e := dest.(*Category)
		e.ID = i
	})

var Title = pgtalk.NewTextAccess(
	tableInfo,
	"title",
	func(dest interface{}, i *string) {
		e := dest.(*Category)
		e.Title = i
	})

func Select(as ...pgtalk.ColumnAccessor) CategorysQuerySet {
	return CategorysQuerySet{pgtalk.MakeQuerySet(tableInfo, as, func() interface{} {
		return new(Category)
	})}
}

type CategorysQuerySet struct {
	pgtalk.QuerySet
}

func (s CategorysQuerySet) Unwrap() pgtalk.QuerySet { return s.QuerySet }

// Where is
func (s CategorysQuerySet) Where(condition pgtalk.SQLWriter) CategorysQuerySet {
	return CategorysQuerySet{QuerySet: s.QuerySet.Where(condition)}
}

// Limit is
func (s CategorysQuerySet) Limit(limit int) CategorysQuerySet {
	return CategorysQuerySet{QuerySet: s.QuerySet.Limit(limit)}
}

// Exec is
func (s CategorysQuerySet) Exec(conn pgtalk.Connection) (list []*Category, err error) {
	err = s.QuerySet.ExecWithAppender(conn, func(each interface{}) {
		list = append(list, each.(*Category))
	})
	return
}
