package tables

import (
	"database/sql"

	"github.com/emicklei/pgtalk"
)

var _ pgtalk.ColumnAccessor = JSONPathAccess{}

type JSONPath struct {
	Valid  bool
	String string
}

type JSONPathAccess struct {
	pgtalk.ColumnInfo
	valueToInsert string
	fieldWriter   func(entity any) any
}

func NewJSONPathAccess(
	info pgtalk.ColumnInfo,
	valueWriter func(dest any) any) JSONPathAccess {
	return JSONPathAccess{
		ColumnInfo:  info,
		fieldWriter: valueWriter}
}

// ColumnInfo is part of ColumnAccessor
func (a JSONPathAccess) Column() pgtalk.ColumnInfo { return a.ColumnInfo }

// AppendScannable is part of ColumnAccessor
func (a JSONPathAccess) AppendScannable(list []any) []any {
	return append(list, &a.valueToInsert)
}

// FieldValueToScan is part of ColumnAccessor
func (a JSONPathAccess) FieldValueToScan(entity any) any {
	return a.fieldWriter(entity)
}

// Get is part of ColumnAccessor
func (a JSONPathAccess) Get(values map[string]any) any {
	v, ok := values[a.Name()] // TODO RawColumnName
	if !ok {
		return sql.NullString{Valid: true, String: "?"}
	}
	return v
}

// ValueToInsert is part of ColumnAccessor
func (a JSONPathAccess) ValueToInsert() any {
	return a.valueToInsert
}

func (a JSONPathAccess) Set(value JSONPath) JSONPathAccess {
	a.valueToInsert = value.String
	return a
}
