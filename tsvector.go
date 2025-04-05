package pgtalk

import "fmt"

// NewTSQuery returns a condition SQL Expression to match @@ a search query with a search vector (column).
func NewTSQuery(info TableInfo, columnName, query string) SQLExpression {
	return tsqueryReader{tableInfo: info, columnName: columnName, query: query}
}

type tsqueryReader struct {
	tableInfo  TableInfo
	columnName string
	query      string
}

func (a tsqueryReader) SQLOn(w WriteContext) {
	fmt.Fprint(w, "(")
	fmt.Fprint(w, w.TableAlias(a.tableInfo.Name, a.tableInfo.Alias))
	fmt.Fprint(w, ".")
	fmt.Fprint(w, a.columnName)
	fmt.Fprint(w, " @@ ")
	fmt.Fprint(w, "to_tsquery('")
	fmt.Fprint(w, a.query)
	fmt.Fprint(w, "'))")
}
func (a tsqueryReader) And(expr SQLExpression) SQLExpression {
	return makeBinaryOperator(a, "AND", expr)
}
func (a tsqueryReader) Or(expr SQLExpression) SQLExpression {
	return makeBinaryOperator(a, "OR", expr)
}

// NewTSQuery returns a ColumnAccessor for reading the value of tsquery typed column.

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
