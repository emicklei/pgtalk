package things

// generated by github.com/emicklei/pgtalk/cmd/pgtalk-gen from version: (devel)
// DO NOT EDIT

import (
	p "github.com/emicklei/pgtalk"
	c "github.com/emicklei/pgtalk/convert"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shopspring/decimal"
	"strings"
	"time"
)

// Thing is generated from the public.things table.
type Thing struct {
	ID         pgtype.UUID         // id : uuid
	Tdate      pgtype.Date         // tdate : date
	Ttimestamp pgtype.Timestamp    // ttimestamp : timestamp without time zone
	Tjsonb     p.NullJSON          // tjsonb : jsonb
	Tjson      p.NullJSON          // tjson : json
	Ttext      pgtype.Text         // ttext : text
	Tnumeric   decimal.NullDecimal // tnumeric : numeric
	Tdecimal   decimal.NullDecimal // tdecimal : numeric
	// for storing custom field expression result values
	expressionResults map[string]any
}

var (
	// ID represents the column "id" of with type "uuid", nullable:true, primary:false
	ID = p.NewFieldAccess[pgtype.UUID](p.MakeColumnInfo(tableInfo, "id", p.NotPrimary, p.Nullable, 0),
		func(dest any) any { return &dest.(*Thing).ID })
	// Tdate represents the column "tdate" of with type "date", nullable:true, primary:false
	Tdate = p.NewFieldAccess[pgtype.Date](p.MakeColumnInfo(tableInfo, "tdate", p.NotPrimary, p.Nullable, 0),
		func(dest any) any { return &dest.(*Thing).Tdate })
	// Ttimestamp represents the column "ttimestamp" of with type "timestamp without time zone", nullable:true, primary:false
	Ttimestamp = p.NewFieldAccess[pgtype.Timestamp](p.MakeColumnInfo(tableInfo, "ttimestamp", p.NotPrimary, p.Nullable, 0),
		func(dest any) any { return &dest.(*Thing).Ttimestamp })
	// Tjsonb represents the column "tjsonb" of with type "jsonb", nullable:true, primary:false
	Tjsonb = p.NewJSONAccess(p.MakeColumnInfo(tableInfo, "tjsonb", p.NotPrimary, p.Nullable, 0),
		func(dest any) any { return &dest.(*Thing).Tjsonb })
	// Tjson represents the column "tjson" of with type "json", nullable:true, primary:false
	Tjson = p.NewJSONAccess(p.MakeColumnInfo(tableInfo, "tjson", p.NotPrimary, p.Nullable, 0),
		func(dest any) any { return &dest.(*Thing).Tjson })
	// Ttext represents the column "ttext" of with type "text", nullable:true, primary:false
	Ttext = p.NewFieldAccess[pgtype.Text](p.MakeColumnInfo(tableInfo, "ttext", p.NotPrimary, p.Nullable, 0),
		func(dest any) any { return &dest.(*Thing).Ttext })
	// Tnumeric represents the column "tnumeric" of with type "numeric", nullable:true, primary:false
	Tnumeric = p.NewFieldAccess[decimal.NullDecimal](p.MakeColumnInfo(tableInfo, "tnumeric", p.NotPrimary, p.Nullable, 0),
		func(dest any) any { return &dest.(*Thing).Tnumeric })
	// Tdecimal represents the column "tdecimal" of with type "numeric", nullable:true, primary:false
	Tdecimal = p.NewFieldAccess[decimal.NullDecimal](p.MakeColumnInfo(tableInfo, "tdecimal", p.NotPrimary, p.Nullable, 0),
		func(dest any) any { return &dest.(*Thing).Tdecimal })
	// package private
	_         = c.UUID // for the occasional unused import from convert
	_         = time.Now
	_         = pgtype.Empty // for the occasional unused import from pgtype
	_         = decimal.Decimal{}
	tableInfo = p.TableInfo{Schema: "public", Name: "things", Alias: "t1"}
)

func init() {
	// after var initialization (to prevent cycle) we need to update the tableInfo to set all columns
	tableInfo.Columns = []p.ColumnAccessor{ID, Tdate, Ttimestamp, Tjsonb, Tjson, Ttext, Tnumeric, Tdecimal}
}

// TableInfo returns meta information about the table.
func TableInfo() p.TableInfo {
	return tableInfo
}

// SetID sets the value to the field value and returns the receiver.
func (e *Thing) SetID(v pgtype.UUID) *Thing { e.ID = v; return e }

// SetTdate sets the value to the field value and returns the receiver.
func (e *Thing) SetTdate(v time.Time) *Thing { e.Tdate = c.TimeToDate(v); return e }

// SetTtimestamp sets the value to the field value and returns the receiver.
func (e *Thing) SetTtimestamp(v time.Time) *Thing { e.Ttimestamp = c.TimeToTimestamp(v); return e }

// SetTjsonb sets the value to the field value and returns the receiver.
func (e *Thing) SetTjsonb(v p.NullJSON) *Thing { e.Tjsonb = v; return e }

// SetTjson sets the value to the field value and returns the receiver.
func (e *Thing) SetTjson(v p.NullJSON) *Thing { e.Tjson = v; return e }

// SetTtext sets the value to the field value and returns the receiver.
func (e *Thing) SetTtext(v string) *Thing { e.Ttext = c.StringToText(v); return e }

// SetTnumeric sets the value to the field value and returns the receiver.
func (e *Thing) SetTnumeric(v decimal.NullDecimal) *Thing { e.Tnumeric = v; return e }

// SetTdecimal sets the value to the field value and returns the receiver.
func (e *Thing) SetTdecimal(v decimal.NullDecimal) *Thing { e.Tdecimal = v; return e }

// Setters returns the list of changes to a Thing for which updates/inserts need to be processed.
// Can be used in Insert,Update,Select. Cannot be used to set null values for columns.
func (e *Thing) Setters() (list []p.ColumnAccessor) {
	if e.ID.Valid {
		list = append(list, ID.Set(e.ID))
	}
	if e.Tdate.Valid {
		list = append(list, Tdate.Set(e.Tdate))
	}
	if e.Ttimestamp.Valid {
		list = append(list, Ttimestamp.Set(e.Ttimestamp))
	}
	if e.Tjsonb.Valid {
		list = append(list, Tjsonb.Set(e.Tjsonb))
	}
	if e.Tjson.Valid {
		list = append(list, Tjson.Set(e.Tjson))
	}
	if e.Ttext.Valid {
		list = append(list, Ttext.Set(e.Ttext))
	}
	if e.Tnumeric.Valid {
		list = append(list, Tnumeric.Set(e.Tnumeric))
	}
	if e.Tdecimal.Valid {
		list = append(list, Tdecimal.Set(e.Tdecimal))
	}
	return
}

// String returns the debug string for *Thing with all non-nil field values.
func (e *Thing) String() string {
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
			if strings.HasPrefix(n, "'") { // mixed case names are quoted
				n = strings.Trim(n, "'")
			}
			if n == each {
				list = append(list, other)
			}
		}
	}
	return
}

// AddExpressionResult puts a value into the custom expression results
func (e *Thing) AddExpressionResult(key string, value any) {
	if e.expressionResults == nil {
		// lazy initialize
		e.expressionResults = map[string]any{}
	}
	e.expressionResults[key] = value
}

// GetExpressionResult gets a value from the custom expression results. Returns nil if absent.
func (e *Thing) GetExpressionResult(key string) any {
	v, ok := e.expressionResults[key]
	if !ok {
		return nil
	}
	pv := v.(*any)
	return *pv
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
