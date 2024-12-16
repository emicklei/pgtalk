package test

import (
	"github.com/emicklei/pgtalk"
	"github.com/jackc/pgx/v5/pgtype"
)

var _ pgtalk.ColumnAccessor = IntervalAccess{}

type IntervalAccess struct {
	info          pgtalk.ColumnInfo
	valueToInsert pgtype.Interval
	fieldWriter   func(entity any) any
}

func NewIntervalAccess(
	info pgtalk.ColumnInfo,
	valueWriter func(dest any) any) IntervalAccess {
	return IntervalAccess{
		info:        info,
		fieldWriter: valueWriter}
}

// ColumnInfo is part of ColumnAccessor
func (a IntervalAccess) Column() pgtalk.ColumnInfo { return a.info }

// AppendScannable is part of ColumnAccessor
func (a IntervalAccess) AppendScannable(list []any) []any {
	return append(list, &a.valueToInsert)
}

// FieldValueToScan is part of ColumnAccessor
func (a IntervalAccess) FieldValueToScan(entity any) any {
	return a.fieldWriter(entity)
}

// Get is part of ColumnAccessor
func (a IntervalAccess) Get(values map[string]any) any {
	v, ok := values[a.info.Name()] // TODO RawColumnName
	if !ok {
		return int64(0)
	}
	return v
}

// Name is part of ColumnAccessor
func (a IntervalAccess) Name() string { return a.info.Name() }

// Name is part of ColumnAccessor
func (a IntervalAccess) SQLOn(w pgtalk.WriteContext) {
	a.info.SQLOn(w)
}

// ValueToInsert is part of ColumnAccessor
func (a IntervalAccess) ValueToInsert() any {
	return a.valueToInsert
}
