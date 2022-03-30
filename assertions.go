package pgtalk

import "fmt"

func assertEachAccessorIn(source []ColumnAccessor, target []ColumnAccessor) {
	for _, each := range source {
		for _, other := range target {
			if each.Column().tableInfo.Equals(other.Column().tableInfo) {
				if each.Column().columnName == other.Column().columnName {
					return
				}
			}
		}
	}
	panic("pgtalk.ASSERT: invalid")
}

func assertEachAccessorHasTableInfo(list []ColumnAccessor, tableInfo TableInfo) {
	for _, each := range list {
		if !each.Column().tableInfo.Equals(tableInfo) {
			// TODO
			panic(fmt.Sprintf("pgtalk.ASSERT: [%v] is not a column from [%v]", each, tableInfo))
		}
	}
}
