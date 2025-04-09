package pgtalk

import (
	"bytes"
	"testing"
)

func TestNewTSQuery(t *testing.T) {
	columnInfo := ColumnInfo{
		tableInfo: TableInfo{
			Name:  "test_table",
			Alias: "t",
		},
		columnName: "search_column",
	}
	query := "test_query"

	expr := NewTSQuery(columnInfo, query)

	tsquery, ok := expr.(tsqueryReader)
	if !ok {
		t.Fatal()
	}
	buf := new(bytes.Buffer)
	wtx := NewWriteContext(buf)
	tsquery.SQLOn(wtx)
	if got, want := buf.String(), "(t.search_column @@ to_tsquery('test_query'))"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}
func TestNewTSQueryWithConfig(t *testing.T) {
	columnInfo := ColumnInfo{
		tableInfo: TableInfo{
			Name:  "test_table",
			Alias: "t",
		},
		columnName: "search_column",
	}
	query := "test_query"

	expr := NewTSQueryWithConfig(columnInfo, "dutch", query)

	tsquery, ok := expr.(tsqueryReader)
	if !ok {
		t.Fatal()
	}
	buf := new(bytes.Buffer)
	wtx := NewWriteContext(buf)
	tsquery.SQLOn(wtx)
	if got, want := buf.String(), "(t.search_column @@ to_tsquery('dutch','test_query'))"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}
