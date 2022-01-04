package things

// generated by pgtalk-gen on Tue, 04 Jan 2022 15:42:32 CET
// DO NOT EDIT

import (
	p "github.com/emicklei/pgtalk"
	"github.com/jackc/pgtype"
	"time"
)

// Thing is generated from the public.things table.
type Thing struct {
	ID         pgtype.UUID      // id : uuid
	Tdate      pgtype.Date      // tdate : date
	Ttimestamp pgtype.Timestamp // ttimestamp : timestamp without time zone
	Tjson      pgtype.JSONB     // tjson : jsonb
}

var (
	// ID represents the column "id" of with type "uuid", nullable:true, primary:false
	ID = p.NewFieldAccess[pgtype.UUID](p.MakeColumnInfo(tableInfo, "id", p.NotPrimary, p.Nullable, 1),
		nil, func(dest interface{}, v pgtype.UUID) { dest.(*Thing).ID = v })
	// Tdate represents the column "tdate" of with type "date", nullable:true, primary:false
	Tdate = p.NewFieldAccess[pgtype.Date](p.MakeColumnInfo(tableInfo, "tdate", p.NotPrimary, p.Nullable, 2),
		nil, func(dest interface{}, v pgtype.Date) { dest.(*Thing).Tdate = v })
	// Ttimestamp represents the column "ttimestamp" of with type "timestamp without time zone", nullable:true, primary:false
	Ttimestamp = p.NewFieldAccess[pgtype.Timestamp](p.MakeColumnInfo(tableInfo, "ttimestamp", p.NotPrimary, p.Nullable, 3),
		nil, func(dest interface{}, v pgtype.Timestamp) { dest.(*Thing).Ttimestamp = v })
	// Tjson represents the column "tjson" of with type "jsonb", nullable:true, primary:false
	Tjson = p.NewJSONBAccess(p.MakeColumnInfo(tableInfo, "tjson", p.NotPrimary, p.Nullable, 4),
		nil, func(dest interface{}, v pgtype.JSONB) { dest.(*Thing).Tjson = v })
	// package private
	_         = time.Now
	_         = pgtype.Empty // for the occasional unused import from pgtype
	tableInfo = p.TableInfo{Schema: "public", Name: "things", Alias: "t1"}
)

func init() {
	// after var initialization (to prevent cycle) we need to update the tableInfo to set all columns
	tableInfo.Columns = []p.ColumnAccessor{ID, Tdate, Ttimestamp, Tjson}
}

// SetID sets the value to the field value and returns the receiver.
func (e *Thing) SetID(v pgtype.UUID) *Thing { e.ID = v; return e }

// SetTdate sets the value to the field value and returns the receiver.
func (e *Thing) SetTdate(v pgtype.Date) *Thing { e.Tdate = v; return e }

// SetTtimestamp sets the value to the field value and returns the receiver.
func (e *Thing) SetTtimestamp(v pgtype.Timestamp) *Thing { e.Ttimestamp = v; return e }

// SetTjson sets the value to the field value and returns the receiver.
func (e *Thing) SetTjson(v pgtype.JSONB) *Thing { e.Tjson = v; return e }

// Setters returns the list of changes to a Thing for which updates/inserts need to be processed.
// Can be used in Insert,Update,Select. Cannot be used to set null values for columns.
func (e *Thing) Setters() (list []p.ColumnAccessor) {
	if e.ID.Status == pgtype.Present {
		list = append(list, ID.Set(e.ID))
	}
	if e.Tdate.Status == pgtype.Present {
		list = append(list, Tdate.Set(e.Tdate))
	}
	if e.Ttimestamp.Status == pgtype.Present {
		list = append(list, Ttimestamp.Set(e.Ttimestamp))
	}
	if e.Tjson.Status == pgtype.Present {
		list = append(list, Tjson.Set(e.Tjson.Bytes))
	}
	return
}

// String returns the debug string for *Thing with all non-nil field values.
func (e *Thing) String() string {
	return p.StringWithFields(e, p.HideNilValues)
}

// AllColumns returns the list of all column accessors for usage in e.g. Select.
func AllColumns() []p.ColumnAccessor {
	return tableInfo.Columns
}

// Select returns a new QuerySet[Thing] for fetching column data.
func Select(cas ...p.ColumnAccessor) p.QuerySet[Thing] {
	return p.MakeQuerySet[Thing](tableInfo, cas)
}

// Insert creates a MutationSet for inserting data with zero or more columns.
func Insert(cas ...p.ColumnAccessor) p.MutationSet[Thing] {
	return p.MakeMutationSet[Thing](tableInfo, cas, p.MutationInsert)
}

// Delete creates a MutationSet for deleting data.
func Delete() p.MutationSet[Thing] {
	return p.MakeMutationSet[Thing](tableInfo, p.EmptyColumnAccessor, p.MutationDelete)
}

// Update creates a MutationSet to update zero or more columns.
func Update(cas ...p.ColumnAccessor) p.MutationSet[Thing] {
	return p.MakeMutationSet[Thing](tableInfo, cas, p.MutationUpdate)
}
