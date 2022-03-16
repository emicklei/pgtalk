package pgtalk

import (
	"time"
)

type TimeAccess struct {
	ColumnInfo
	valueFieldWriter FieldAccessFunc
	valueToInsert    time.Time
}

func NewTimeAccess(info ColumnInfo,
	valueWriter FieldAccessFunc) TimeAccess {
	return TimeAccess{ColumnInfo: info, valueFieldWriter: valueWriter}
}

// Collect is part of SQLExpression
func (a TimeAccess) Collect(list []ColumnAccessor) []ColumnAccessor {
	return append(list, a)
}

func (a TimeAccess) ValueToInsert() interface{} {
	return a.valueToInsert
}

func (a TimeAccess) Set(v time.Time) TimeAccess {
	a.valueToInsert = v
	return a
}

func (a TimeAccess) Column() ColumnInfo { return a.ColumnInfo }

func (a TimeAccess) FieldToScan(entity any) any {
	return a.valueFieldWriter(entity)
}

func (a TimeAccess) TableAlias(alias string) TimeAccess {
	a.ColumnInfo = a.ColumnInfo.TableAlias(alias)
	return a
}
