package things

// generated by pgtalk-gen on 2021-08-10 11:59:41.197675 &#43;0200 CEST m=&#43;0.026252712
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
	tableInfo = pgtalk.TableInfo{Schema: "public", Name: "things", Alias: "t1"}
)

type Thing struct {
	TDate      *time.Time // date
	TTimestamp *time.Time // timestamp without time zone
	TJSON      *string    // jsonb
	ID         *int64     // bigint
}

var (
	TDate = pgtalk.NewTimeAccess(pgtalk.MakeColumnInfo(tableInfo, "tDate", false, false),
		func(dest interface{}, v *time.Time) { dest.(*Thing).TDate = v })
	TTimestamp = pgtalk.NewTimeAccess(pgtalk.MakeColumnInfo(tableInfo, "tTimestamp", false, false),
		func(dest interface{}, v *time.Time) { dest.(*Thing).TTimestamp = v })
	TJSON = pgtalk.NewJSONBAccess(pgtalk.MakeColumnInfo(tableInfo, "tJSON", false, false),
		func(dest interface{}, v *string) { dest.(*Thing).TJSON = v })
	ID = pgtalk.NewInt64Access(pgtalk.MakeColumnInfo(tableInfo, "id", true, true),
		func(dest interface{}, v *int64) { dest.(*Thing).ID = v })
)

// ColumnUpdatesFrom returns the list of changes to a Thing for which updates need to be processed.
// Cannot be used to set null values for columns.
func ColumnUpdatesFrom(e Thing) (list []pgtalk.SQLWriter) {
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

// setColumnValueTo sets the field of a *Thing to the non-nil value.
func setColumnValueTo(e *Thing, tableAttributeNumber uint16, value interface{}) error {
	if value == nil {
		return nil
	}
	switch tableAttributeNumber {
	case 1:
		if tvalue, ok := value.(time.Time); !ok {
			return fmt.Errorf("unable to assert value of type [*%T] to field [TDate] of type [%s] in entity with type [%s]", value, "*time.Time", "*things.Thing")
		} else {
			e.TDate = &tvalue
		}
	case 2:
		if tvalue, ok := value.(time.Time); !ok {
			return fmt.Errorf("unable to assert value of type [*%T] to field [TTimestamp] of type [%s] in entity with type [%s]", value, "*time.Time", "*things.Thing")
		} else {
			e.TTimestamp = &tvalue
		}
	case 3:
		if tvalue, ok := value.(string); !ok {
			return fmt.Errorf("unable to assert value of type [*%T] to field [TJSON] of type [%s] in entity with type [%s]", value, "*string", "*things.Thing")
		} else {
			e.TJSON = &tvalue
		}
	case 4:
		if tvalue, ok := value.(int64); !ok {
			return fmt.Errorf("unable to assert value of type [*%T] to field [ID] of type [%s] in entity with type [%s]", value, "*int64", "*things.Thing")
		} else {
			e.ID = &tvalue
		}
	default:
		return fmt.Errorf("unable to set value [%v] to field of [%v] with table attribute number [%d] in entity with type [%s]", value, e, tableAttributeNumber, "*things.Thing")
	}
	return nil
}

var fieldSetter = func(entityPointer interface{}, tableAttributeNumber uint16, value interface{}) error {
	return setColumnValueTo(entityPointer.(*Thing), tableAttributeNumber, value)
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
func AllColumns() (all []pgtalk.ColumnAccessor) {
	return append(all, TDate, TTimestamp, TJSON, ID)
}

// Select returns a new ThingsQuerySet for fetching column data.
func Select(cas ...pgtalk.ColumnAccessor) ThingsQuerySet {
	return ThingsQuerySet{pgtalk.MakeQuerySet(tableInfo, cas, func() interface{} {
		return new(Thing)
	}, fieldSetter)}
}

// ThingsQuerySet can query for *Thing values.
type ThingsQuerySet struct {
	pgtalk.QuerySet
}

func (s ThingsQuerySet) Unwrap() pgtalk.QuerySet { return s.QuerySet }

// Where returns a new QuerySet with WHERE clause.
func (s ThingsQuerySet) Where(condition pgtalk.SQLWriter) ThingsQuerySet {
	return ThingsQuerySet{QuerySet: s.QuerySet.Where(condition)}
}

// Limit returns a new QuerySet with the maximum number of results set.
func (s ThingsQuerySet) Limit(limit int) ThingsQuerySet {
	return ThingsQuerySet{QuerySet: s.QuerySet.Limit(limit)}
}

// GroupBy returns a new QuerySet with the GROUP BY clause.
func (s ThingsQuerySet) GroupBy(cas ...pgtalk.ColumnAccessor) ThingsQuerySet {
	return ThingsQuerySet{QuerySet: s.QuerySet.GroupBy(cas...)}
}

// OrderBy returns a new QuerySet with the ORDER BY clause.
func (s ThingsQuerySet) OrderBy(cas ...pgtalk.ColumnAccessor) ThingsQuerySet {
	return ThingsQuerySet{QuerySet: s.QuerySet.OrderBy(cas...)}
}

// Exec runs the query and returns the list of *Thing.
func (s ThingsQuerySet) Exec(ctx context.Context, conn *pgx.Conn) (list []*Thing, err error) {
	err = s.QuerySet.ExecWithAppender(ctx, conn, func(each interface{}) {
		list = append(list, each.(*Thing))
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
