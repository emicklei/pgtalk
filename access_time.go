package pgtalk

import (
	"time"
)

type timeAccess struct {
	ColumnInfo
	valueFieldWriter fieldAccessFunc
	valueToInsert    time.Time
}

func NewTimeAccess(info ColumnInfo,
	valueWriter fieldAccessFunc) timeAccess {
	return timeAccess{ColumnInfo: info, valueFieldWriter: valueWriter}
}

func (a timeAccess) ValueToInsert() any {
	return a.valueToInsert
}

func (a timeAccess) Set(v time.Time) timeAccess {
	a.valueToInsert = v
	return a
}

func (a timeAccess) Column() ColumnInfo { return a.ColumnInfo }

func (a timeAccess) FieldValueToScan(entity any) any {
	return a.valueFieldWriter(entity)
}

// TableAlias changes the table alias for this column accessor.
func (a timeAccess) TableAlias(alias string) timeAccess {
	a.ColumnInfo = a.ColumnInfo.TableAlias(alias)
	return a
}

// AppendScannable is part of ColumnAccessor
func (a timeAccess) AppendScannable(list []any) []any {
	return append(list, &a.valueToInsert)
}

// Get returns the value for its columnName from a map (row).
func (a timeAccess) Get(values map[string]any) any {
	v, ok := values[a.columnName]
	if !ok {
		return time.Time{}
	}
	return v
}
