package pgtalk

import (
	"fmt"
	"io"
)

// TableInfo describes a table in the database.
type TableInfo struct {
	Name   string
	Schema string // e.g. public
	Alias  string
	// Columns are all known columns for this table ;initialized by the generated package
	Columns []ColumnAccessor
}

// FullName returns the fully qualified name of the table.
func (t TableInfo) FullName() string {
	return fmt.Sprintf("%s.%s", t.Schema, t.Name)
}

func (t TableInfo) SQLOn(w io.Writer) {
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
