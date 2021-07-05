package pgtalk

import (
	"fmt"
	"io"
	"strings"
)

// Int64Access can Read a column value (int8) and Write a column value and Set a struct field (int64).
type Int64Access struct {
	columnInfo
	fieldWriter func(dest interface{}, i *int64)
	insertValue int64
}

func NewInt64Access(info TableInfo, columnName string, writer func(dest interface{}, i *int64)) Int64Access {
	return Int64Access{columnInfo: columnInfo{tableInfo: info, columnName: columnName}, fieldWriter: writer}
}

func (a Int64Access) ValueAsSQLOn(w io.Writer) {
	fmt.Fprintf(w, "%d", a.insertValue)
}

func (a Int64Access) BetweenAnd(begin int64, end int64) BetweenAnd {
	return MakeBetweenAnd(a, ValuePrinter{begin}, ValuePrinter{end})
}

func (a Int64Access) WriteInto(entity interface{}, fieldValue interface{}) {
	var i int64 = fieldValue.(int64)
	a.fieldWriter(entity, &i)
}

func (a Int64Access) Set(v int64) Int64Access {
	a.insertValue = v
	return a
}

func (a Int64Access) Equals(i int) BinaryOperator {
	return MakeBinaryOperator(a, "=", ValuePrinter{i})
}

func (a Int64Access) Compare(op string, i int) BinaryOperator {
	if !strings.Contains(validComparisonOperators, op) {
		panic("invalid comparison operator:" + op)
	}
	return MakeBinaryOperator(a, op, ValuePrinter{i})
}

func (a Int64Access) NotNull() NullCheck {
	return NullCheck{Operand: a, IsNot: true}
}
