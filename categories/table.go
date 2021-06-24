package categories

import (
	"github.com/emicklei/pgtalk/xs"
	"github.com/jackc/pgx/v4"
)

var tableInfo = xs.TableInfo{Name: "categories", Alias: "t2"}

type Category struct {
	ID    *int64
	Title *string
}

var ID = xs.NewInt8Access(
	tableInfo,
	"id",
	func(dest interface{}, i *int64) {
		e := dest.(*Category)
		e.ID = i
	})

var Title = xs.NewTextAccess(
	tableInfo,
	"title",
	func(dest interface{}, i *string) {
		e := dest.(*Category)
		e.Title = i
	})

func Select(as ...xs.ReadWrite) CategorysQuerySet {
	return CategorysQuerySet{xs.MakeQuerySet(tableInfo, as, func() interface{} {
		return new(Category)
	})}
}

type CategorysQuerySet struct {
	xs.QuerySet
}

func (s CategorysQuerySet) Unwrap() xs.QuerySet { return s.QuerySet }

// Where is
func (s CategorysQuerySet) Where(condition xs.SQLWriter) CategorysQuerySet {
	return CategorysQuerySet{QuerySet: s.QuerySet.Where(condition)}
}

// Limit is
func (s CategorysQuerySet) Limit(limit int) CategorysQuerySet {
	return CategorysQuerySet{QuerySet: s.QuerySet.Limit(limit)}
}

// Exec is
func (s CategorysQuerySet) Exec(conn *pgx.Conn) (list []*Category, err error) {
	err = s.QuerySet.Exec(conn, func(each interface{}) {
		list = append(list, each.(*Category))
	})
	return
}
