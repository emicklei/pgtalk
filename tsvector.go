package pgtalk

import "fmt"

// NewTSQuery returns a condition SQL Expression to match @@ a search query with a search vector (column).
func NewTSQuery(columnInfo ColumnInfo, query string) SQLExpression {
	return tsqueryReader{columnInfo: columnInfo, query: query}
}

// NewTSQuery returns a condition SQL Expression to match @@ a search query with a search vector (column) and a config.
func NewTSQueryWithConfig(columnInfo ColumnInfo, regconfig, query string) SQLExpression {
	return tsqueryReader{columnInfo: columnInfo, regconfig: regconfig, query: query}
}

type tsqueryReader struct {
	columnInfo ColumnInfo
	regconfig  string
	query      string
}

func (a tsqueryReader) SQLOn(w WriteContext) {
	fmt.Fprint(w, "(")
	fmt.Fprint(w, w.TableAlias(a.columnInfo.tableInfo.Name, a.columnInfo.tableInfo.Alias))
	fmt.Fprint(w, ".")
	fmt.Fprint(w, a.columnInfo.columnName)
	fmt.Fprint(w, " @@ ")
	fmt.Fprint(w, "to_tsquery('")
	if a.regconfig != "" {
		fmt.Fprint(w, a.regconfig)
		fmt.Fprint(w, "','")
	}
	fmt.Fprint(w, a.query)
	fmt.Fprint(w, "'))")
}
func (a tsqueryReader) And(expr SQLExpression) SQLExpression {
	return makeBinaryOperator(a, "AND", expr)
}
func (a tsqueryReader) Or(expr SQLExpression) SQLExpression {
	return makeBinaryOperator(a, "OR", expr)
}

type tsvectorWriter struct {
	ColumnInfo
	regconfig string
	value     string
}

// NewTSVector returns a ColumnAccessor for writing the value of tsvector typed column.
// Cannot be used for reading the value of such a column.
func NewTSVector(columnInfo ColumnInfo, value string) ColumnAccessor {
	return tsvectorWriter{ColumnInfo: columnInfo, value: value}
}

// NewTSVectorWithConfig returns a ColumnAccessor for writing the value of tsvector typed column.
// Cannot be used for reading the value of such a column.
func NewTSVectorWithConfig(columnInfo ColumnInfo, regconfig, value string) ColumnAccessor {
	return tsvectorWriter{ColumnInfo: columnInfo, regconfig: regconfig, value: value}
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
	if a.regconfig != "" {
		return fmt.Sprintf("to_tsvector('%s', $%d)", a.regconfig, parameterIndex+1)
	}
	return fmt.Sprintf("to_tsvector($%d)", parameterIndex)
}
