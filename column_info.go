package pgtalk

import (
	"fmt"
	"strconv"
	"strings"
)

type ColumnInfo struct {
	tableInfo            TableInfo
	columnName           string
	isPrimary            bool
	notNull              bool
	isMixedCase          bool
	tableAttributeNumber uint16
}

func MakeColumnInfo(tableInfo TableInfo, columnName string, isPrimary bool, isNotNull bool, tableAttributeNumber uint16) ColumnInfo {
	return ColumnInfo{
		tableInfo:            tableInfo,
		columnName:           columnName,
		notNull:              isNotNull,
		isPrimary:            isPrimary,
		isMixedCase:          strings.ToLower(columnName) != columnName,
		tableAttributeNumber: tableAttributeNumber,
	}
}

func (c ColumnInfo) Name() string {
	if c.isMixedCase {
		return strconv.Quote(c.columnName)
	}
	return c.columnName
}

func (c ColumnInfo) SQLOn(w WriteContext) {
	fmt.Fprintf(w, "%s.%s", w.TableAlias(c.tableInfo.Name, c.tableInfo.Alias), c.Name())
}

// TableAlias changes the table alias for this column info.
func (c ColumnInfo) TableAlias(alias string) ColumnInfo {
	c.tableInfo = c.tableInfo.WithAlias(alias)
	return c
}
