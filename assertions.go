package pgtalk

import "fmt"

func assertEachAccessorIn(source []ColumnAccessor, target []ColumnAccessor) {
	for _, each := range source {
		for _, other := range target {
			if each.Column().tableInfo == other.Column().tableInfo {
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

type ValueConversionError struct {
	got, want string
}

func NewValueConversionError(got interface{}, want string) error {
	return ValueConversionError{fmt.Sprintf("%T", got), want}
}

func (e ValueConversionError) Error() string {
	return fmt.Sprintf("field value conversion error, got %s expected %s", e.got, e.want)
}
