package things

// generated by pgtalk-gen on 2021-08-23 17:30:49.416417 &#43;0200 CEST m=&#43;0.027916047
// DO NOT EDIT

import (
	"bytes"
	"fmt"
	"github.com/emicklei/pgtalk"
	"time"
)

var (
	_         = time.Now()
	tableInfo = pgtalk.TableInfo{Schema: "public", Name: "things", Alias: "t1"} 
)

type Thing struct {
	TDate      *time.Time // date
	TTimestamp *time.Time // timestamp without time zone
	TJSON      *string    // jsonb
	ID         *int64     // bigint
}

var (
	TDate = pgtalk.NewTimeAccess(pgtalk.MakeColumnInfo(tableInfo, "tDate", false, false, 1),
		func(dest interface{}, v *time.Time) { dest.(*Thing).TDate = v })
	TTimestamp = pgtalk.NewTimeAccess(pgtalk.MakeColumnInfo(tableInfo, "tTimestamp", false, false, 2),
		func(dest interface{}, v *time.Time) { dest.(*Thing).TTimestamp = v })
	TJSON = pgtalk.NewJSONBAccess(pgtalk.MakeColumnInfo(tableInfo, "tJSON", false, false, 3),
		func(dest interface{}, v *string) { dest.(*Thing).TJSON = v })
	ID = pgtalk.NewInt64Access(pgtalk.MakeColumnInfo(tableInfo, "id", true, true, 4),
		func(dest interface{}, v *int64) { dest.(*Thing).ID = v })
	tableAccess = pgtalk.TableAccessor{TableInfo: tableInfo,
		Factory: func() interface{} { return new(Thing) }, AllColumns: []pgtalk.ColumnAccessor{TDate, TTimestamp, TJSON, ID}}
)

// ColumnUpdatesFrom returns the list of changes to a Thing for which updates need to be processed.
// Cannot be used to set null values for columns.
func ColumnUpdatesFrom(e *Thing) (list []pgtalk.SQLExpression) {
	if e.TDate != nil {
		list = append(list, TDate.Set(*e.TDate))
	}
	if e.TTimestamp != nil {
		list = append(list, TTimestamp.Set(*e.TTimestamp))
	}
	if e.TJSON != nil {
		list = append(list, TJSON.Set(*e.TJSON))
	}
	if e.ID != nil {
		list = append(list, ID.Set(*e.ID))
	}
	return
}

// Next returns the next *Thing from the iterator data.
// Use err to check for failure.
func Next(it *pgtalk.ResultIterator) (e *Thing, err error) {
	var each = new(Thing)
	// first check for query error in case caller forgot
	if err = it.Err(); err != nil {
		return nil, err
	}
	err = it.Next(each)
	return each, err
}

// String returns the debug string for *Thing with all non-nil field values.
func (e *Thing) String() string {
	b := new(bytes.Buffer)
	fmt.Fprint(b, "things.Thing{")
	if e.TDate != nil {
		fmt.Fprintf(b, "TDate:%v ", *e.TDate)
	}
	if e.TTimestamp != nil {
		fmt.Fprintf(b, "TTimestamp:%v ", *e.TTimestamp)
	}
	if e.TJSON != nil {
		fmt.Fprintf(b, "TJSON:%v ", *e.TJSON)
	}
	if e.ID != nil {
		fmt.Fprintf(b, "ID:%v ", *e.ID)
	}
	fmt.Fprint(b, "}")
	return b.String()
}

// AllColumns returns the list of all column accessors for usage in e.g. Select.
func AllColumns() []pgtalk.ColumnAccessor {
	return tableAccess.AllColumns
}

// Select returns a new QuerySet[Thing] for fetching column data.
func Select(cas ...pgtalk.ColumnAccessor) pgtalk.QuerySet[Thing] {
	return pgtalk.MakeQuerySet[Thing](tableAccess, cas)
}

// Insert creates a MutationSet for inserting data with zero or more columns.
func Insert(cas ...pgtalk.ColumnAccessor) pgtalk.MutationSet[Thing] {
	return pgtalk.MakeMutationSet[Thing](tableAccess, cas, pgtalk.MutationInsert)
}

// Delete creates a MutationSet for deleting data.
func Delete() pgtalk.MutationSet[Thing] {
	return pgtalk.MakeMutationSet[Thing](tableAccess, pgtalk.EmptyColumnAccessor, pgtalk.MutationDelete)
}

// Update creates a MutationSet to update zero or more columns.
func Update(cas ...pgtalk.ColumnAccessor) pgtalk.MutationSet[Thing] {
	return pgtalk.MakeMutationSet[Thing](tableAccess, cas, pgtalk.MutationUpdate)
}
