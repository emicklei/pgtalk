package pgtalk

import (
	"time"
)

type TimeAccess struct {
	ColumnInfo
	valueFieldWriter fieldAccessFunc
	valueToInsert    time.Time
}

func NewTimeAccess(info ColumnInfo,
	valueWriter fieldAccessFunc) TimeAccess {
	return TimeAccess{ColumnInfo: info, valueFieldWriter: valueWriter}
}

func (a TimeAccess) ValueToInsert() any {
	return a.valueToInsert
}

func (a TimeAccess) Set(v time.Time) TimeAccess {
	a.valueToInsert = v
	return a
}

func (a TimeAccess) Column() ColumnInfo { return a.ColumnInfo }

func (a TimeAccess) FieldValueToScan(entity any) any {
	return a.valueFieldWriter(entity)
}

func (a TimeAccess) TableAlias(alias string) TimeAccess {
	a.ColumnInfo = a.ColumnInfo.TableAlias(alias)
	return a
}

// AppendScannable is part of ColumnAccessor
func (a TimeAccess) AppendScannable(list []any) []any {
	return append(list, &a.valueToInsert)
}

// Get returns the value for its columnName from a map (row).
func (a TimeAccess) Get(values map[string]any) any {
	v, ok := values[a.columnName]
	if !ok {
		return time.Time{}
	}
	return v
}
