package pgtalk

import "fmt"

func NewTSQuery(info TableInfo, columnName, query string) SQLExpression {
	return binaryExpression{
		Left:     NewSQLSource(fmt.Sprintf("%s.%s", info.Alias, columnName)),
		Operator: "@@",
		Right:    NewSQLSource(fmt.Sprintf("to_tsquery('%s')", query)),
	}
}

type tsvectorWriter struct {
	ColumnInfo
	value string
}

// NewTSVector returns a ColumnAccessor for writing the value of tsvector typed column.
// Cannot be used for reading the value of such a column.
func NewTSVector(columnName string, value string) ColumnAccessor {
	return tsvectorWriter{ColumnInfo: ColumnInfo{columnName: columnName}, value: value}
}

// AppendScannable is part of ColumnAccessor
func (a tsvectorWriter) AppendScannable(list []any) []any {
	// write only
	return list
}
func (a tsvectorWriter) Column() ColumnInfo { return a.ColumnInfo }

func (a tsvectorWriter) FieldValueToScan(entity any) any {
	// write only
	return nil
}

func (a tsvectorWriter) ValueToInsert() any {
	return a.value
}

// Get returns the value for its columnName from a map (row).
func (a tsvectorWriter) Get(values map[string]any) any {
	// write only
	return a.value
}

func (a tsvectorWriter) SetSource(parameterIndex int) string {
	return fmt.Sprintf("to_tsvector($%d)", parameterIndex)
}
