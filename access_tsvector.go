package pgtalk

import (
	"fmt"
	"time"
)

type tsvectorWriteOnly struct {
	ColumnInfo
	valueToInsert string
}

func NewTSVectorWriter(info ColumnInfo, _ fieldAccessFunc) tsvectorWriteOnly {
	return tsvectorWriteOnly{ColumnInfo: info}
}

func (a tsvectorWriteOnly) ValueToInsert() any {
	return a.valueToInsert
}

func (a tsvectorWriteOnly) ToTSVector(v string) tsvectorWriteOnly {
	a.valueToInsert = v
	return a
}

func (a tsvectorWriteOnly) SetSource(parameterIndex int) string {
	return fmt.Sprintf("to_tsvector($%d)", parameterIndex)
}

func (a tsvectorWriteOnly) FieldValueToScan(entity any) any {
	var ignore any
	return &ignore
}

func (a tsvectorWriteOnly) Column() ColumnInfo { return a.ColumnInfo }

// TableAlias changes the table alias for this column accessor.
func (a tsvectorWriteOnly) TableAlias(alias string) tsvectorWriteOnly {
	a.ColumnInfo = a.ColumnInfo.TableAlias(alias)
	return a
}

// AppendScannable is part of ColumnAccessor
func (a tsvectorWriteOnly) AppendScannable(list []any) []any {
	return list
}

// Get returns the value for its columnName from a map (row).
func (a tsvectorWriteOnly) Get(values map[string]any) any {
	v, ok := values[a.columnName]
	if !ok {
		return time.Time{}
	}
	return v
}
