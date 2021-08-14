package pgtalk

import (
	"fmt"
	"io"
	"strings"
)

// Int64Access can Read a column value (int8) and Write a column value and Set a struct field (int64).
type Int64Access struct {
	ColumnInfo
	fieldWriter   func(dest interface{}, i *int64)
	valueToInsert int64
}

func NewInt64Access(
	info ColumnInfo,
	writer func(dest interface{}, i *int64)) Int64Access {
	return Int64Access{
		ColumnInfo:  info,
		fieldWriter: writer}
}

func (a Int64Access) ValueAsSQLOn(w io.Writer) {
	fmt.Fprintf(w, "%d", a.valueToInsert)
}

func (a Int64Access) BetweenAnd(begin int64, end int64) BetweenAnd {
	return MakeBetweenAnd(a, ValuePrinter{begin}, ValuePrinter{end})
}

func (a Int64Access) SetFieldValue(entity interface{}, fieldValue interface{}) error {
	if fieldValue == nil {
		return nil
	}
	i, ok := fieldValue.(int64)
	if !ok {
		return NewValueConversionError(fieldValue, "int64")
	}
	a.fieldWriter(entity, &i)
	return nil
}

func (a Int64Access) ValueToInsert() interface{} {
	return a.valueToInsert
}

func (a Int64Access) Set(v int64) Int64Access {
	a.valueToInsert = v
	return a
}

func (a Int64Access) Equals(intOrInt64Access interface{}) BinaryOperator {
	if i, ok := intOrInt64Access.(int); ok {
		return MakeBinaryOperator(a, "=", ValuePrinter{i})
	}
	if ia, ok := intOrInt64Access.(Int64Access); ok {
		return MakeBinaryOperator(a, "=", ia)
	}
	panic("int or Int64Access expected")
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

func (a Int64Access) Column() ColumnInfo { return a.ColumnInfo }
