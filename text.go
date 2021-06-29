package pgtalk

import (
	"fmt"
	"strings"
)

type TextAccess struct {
	columnInfo
	fieldWriter func(dest interface{}, i *string)
	insertValue string
}

func NewTextAccess(info TableInfo, columnName string, writer func(dest interface{}, i *string)) TextAccess {
	return TextAccess{columnInfo: columnInfo{tableInfo: info, columnName: columnName}, fieldWriter: writer}
}

func (a TextAccess) Value(v string) TextAccess {
	a.insertValue = v
	return a
}

func (a TextAccess) Equals(s string) BinaryOperator {
	return MakeBinaryOperator(a, "=", LiteralString(s))
}

func (a TextAccess) Compare(op string, s string) BinaryOperator {
	if !strings.Contains(validComparisonOperators, op) {
		panic("invalid comparison operator:" + op)
	}
	return MakeBinaryOperator(a, op, LiteralString(s))
}

func (a TextAccess) WriteInto(entity interface{}, fieldValue interface{}) {
	var i string = fieldValue.(string)
	a.fieldWriter(entity, &i)
}

func (a TextAccess) ValueAsSQL() string {
	return fmt.Sprintf("'%s'", a.insertValue)
}
