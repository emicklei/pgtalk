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
