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

type BooleanAccess struct {
	ColumnInfo
	valueFieldWriter FieldAccessFunc
	valueToInsert    bool
}

func NewBooleanAccess(info ColumnInfo, writer FieldAccessFunc) BooleanAccess {
	return BooleanAccess{ColumnInfo: info, valueFieldWriter: writer}
}

func (a BooleanAccess) Collect(list []ColumnAccessor) []ColumnAccessor {
	return append(list, a)
}

func (a BooleanAccess) Column() ColumnInfo { return a.ColumnInfo }

func (a BooleanAccess) Set(v bool) BooleanAccess {
	a.valueToInsert = v
	return a
}
func (a BooleanAccess) ValueToInsert() interface{} {
	return a.valueToInsert
}

func (a BooleanAccess) Equals(b bool) SQLExpression {
	return MakeBinaryOperator(a, "=", valuePrinter{b})
}

func (a BooleanAccess) FieldToScan(entity any) any {
	return a.valueFieldWriter(entity)
}
