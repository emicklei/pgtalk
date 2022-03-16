package pgtalk

import (
	"fmt"
	"io"
)

type TableInfo struct {
	Name   string
	Schema string
	Alias  string
	// Columns are all known for this table ;initialized by the package
	Columns []ColumnAccessor
}

func (t TableInfo) SQLOn(tableInfo TableInfo, w io.Writer) {
	fmt.Fprintf(w, "%s.%s %s", t.Schema, t.Name, t.Alias)
}

func (t TableInfo) Equals(o TableInfo) bool {
	return t.Name == o.Name && t.Schema == o.Schema && t.Alias == o.Alias
}

func (t TableInfo) String() string {
	return fmt.Sprintf("table(%s.%s %s)", t.Schema, t.Name, t.Alias)
}

func (t TableInfo) WithAlias(aliasName string) TableInfo {
	if aliasName == "" {
		panic("alias cannot be empty")
	}
	t.Alias = aliasName
	return t
}
