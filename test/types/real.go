package types

import (
	"database/sql"

	"github.com/emicklei/pgtalk"
)

var _ pgtalk.ColumnAccessor = RealAccess{}

type Real struct {
	Valid bool
	Float float32
}

// Scan implements the database/sql Scanner interface.
func (f *Real) Scan(src any) error {
	if src == nil {
		*f = Real{}
		return nil
	}
	switch src := src.(type) {
	case float64:
		*f = Real{Float: float32(src), Valid: true}
	}
	return nil
}

type RealAccess struct {
	pgtalk.ColumnInfo
	valueToInsert float32
	fieldWriter   func(entity any) any
}

func NewRealAccess(
	info pgtalk.ColumnInfo,
	valueWriter func(dest any) any) RealAccess {
	return RealAccess{
		ColumnInfo:  info,
		fieldWriter: valueWriter}
}

// ColumnInfo is part of ColumnAccessor
func (a RealAccess) Column() pgtalk.ColumnInfo { return a.ColumnInfo }

// AppendScannable is part of ColumnAccessor
func (a RealAccess) AppendScannable(list []any) []any {
	return append(list, &a.valueToInsert)
}

// FieldValueToScan is part of ColumnAccessor
func (a RealAccess) FieldValueToScan(entity any) any {
	return a.fieldWriter(entity)
}

// Get is part of ColumnAccessor
func (a RealAccess) Get(values map[string]any) any {
	v, ok := values[a.Name()] // TODO RawColumnName
	if !ok {
		return sql.NullFloat64{Valid: false}
	}
	return v
}

// ValueToInsert is part of ColumnAccessor
func (a RealAccess) ValueToInsert() any {
	return a.valueToInsert
}

// Set is not (yet) part of ColumnAccessor
func (a RealAccess) Set(value Real) RealAccess {
	a.valueToInsert = value.Float
	return a
}
