package pgtalk

import (
	"fmt"
	"strings"
)

// Int8Access can Read a column value (int8) and Write a column value and Set a struct field (int64).
type Int8Access struct {
	columnInfo
	fieldWriter func(dest interface{}, i *int64)
	insertValue int64
}

func NewInt8Access(info TableInfo, columnName string, writer func(dest interface{}, i *int64)) Int8Access {
	return Int8Access{columnInfo: columnInfo{tableInfo: info, columnName: columnName}, fieldWriter: writer}
}

func (a Int8Access) ValueAsSQL() string {
	return fmt.Sprintf("%d", a.insertValue)
}

func (a Int8Access) BetweenAnd(begin int64, end int64) BetweenAnd {
	return MakeBetweenAnd(a, ValuePrinter{begin}, ValuePrinter{end})
}

func (a Int8Access) WriteInto(entity interface{}, fieldValue interface{}) {
	var i int64 = fieldValue.(int64)
	a.fieldWriter(entity, &i)
}

func (a Int8Access) Value(v int64) Int8Access {
	a.insertValue = v
	return a
}

func (a Int8Access) Equals(i int) BinaryOperator {
	return MakeBinaryOperator(a, "=", ValuePrinter{i})
}

func (a Int8Access) Compare(op string, i int) BinaryOperator {
	if !strings.Contains(validComparisonOperators, op) {
		panic("invalid comparison operator:" + op)
	}
	return MakeBinaryOperator(a, op, ValuePrinter{i})
}

func (a Int8Access) NotNull() NullCheck {
	return NullCheck{Operand: a, IsNot: true}
}
