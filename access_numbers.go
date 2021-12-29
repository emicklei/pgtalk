package pgtalk

import (
	"strings"

	"github.com/jackc/pgtype"
)

// Int64Access can Read a column value (int8) and Write a column value and Set a struct field (int64).
type Int64Access struct {
	ColumnInfo
	fieldWriter         func(dest interface{}, i int64) // either or
	nullableFieldWriter func(dest interface{}, i pgtype.Int8)
	valueToInsert       int64
}

func NewInt64Access(
	info ColumnInfo,
	valueWriter func(dest interface{}, i int64),
	nullableWriter func(dest interface{}, i pgtype.Int8)) Int64Access {
	return Int64Access{
		ColumnInfo:          info,
		fieldWriter:         valueWriter,
		nullableFieldWriter: nullableWriter}
}

func (a Int64Access) Collect(list []ColumnAccessor) []ColumnAccessor {
	return append(list, a)
}

func (a Int64Access) BetweenAnd(begin int64, end int64) BetweenAnd {
	return MakeBetweenAnd(a, valuePrinter{begin}, valuePrinter{end})
}

func (a Int64Access) SetFieldValue(entity interface{}, fieldValue interface{}) error {
	if fieldValue == nil {
		return nil
	}
	i, ok := fieldValue.(int64)
	if !ok {
		return NewValueConversionError(fieldValue, "int64")
	}
	if a.notNull {
		a.fieldWriter(entity, i)
	} else {
		a.nullableFieldWriter(entity, pgtype.Int8{Int: i, Status: pgtype.Present})
	}
	return nil
}

func (a Int64Access) ValueToInsert() interface{} {
	return a.valueToInsert
}

func (a Int64Access) Set(v int64) Int64Access {
	a.valueToInsert = v
	return a
}

func (a Int64Access) Equals(intOrInt64Access interface{}) binaryExpression {
	if i, ok := intOrInt64Access.(int); ok {
		return MakeBinaryOperator(a, "=", valuePrinter{i})
	}
	if ia, ok := intOrInt64Access.(Int64Access); ok {
		return MakeBinaryOperator(a, "=", ia)
	}
	panic("int or Int64Access expected")
}

func (a Int64Access) Compare(op string, i int) binaryExpression {
	if !strings.Contains(validComparisonOperators, op) {
		panic("invalid comparison operator:" + op)
	}
	return MakeBinaryOperator(a, op, valuePrinter{i})
}

func (a Int64Access) Column() ColumnInfo { return a.ColumnInfo }

// Float64Access can Read a column value (float) and Write a column value and Set a struct field (float64).
type Float64Access struct {
	ColumnInfo
	fieldWriter         func(dest interface{}, f float64)
	nullableFieldWriter func(dest interface{}, f pgtype.Float8)
	valueToInsert       float64
}

func NewFloat64Access(info ColumnInfo, writer func(dest interface{}, f float64), nullableWriter func(dest interface{}, f pgtype.Float8)) Float64Access {
	return Float64Access{ColumnInfo: info, fieldWriter: writer, nullableFieldWriter: nullableWriter}
}

func (a Float64Access) SetFieldValue(entity interface{}, fieldValue interface{}) error {
	if fieldValue == nil {
		return nil
	}
	f, ok := fieldValue.(float64)
	if !ok {
		// TODO try string?
		return NewValueConversionError(fieldValue, "float64")
	}
	a.fieldWriter(entity, f)
	return nil
}

func (a Float64Access) ValueToInsert() interface{} {
	return a.ValueToInsert
}

func (a Float64Access) Set(v float64) Float64Access {
	a.valueToInsert = v
	return a
}

func (a Float64Access) Column() ColumnInfo { return a.ColumnInfo }

func (a Float64Access) Collect(list []ColumnAccessor) []ColumnAccessor {
	return append(list, a)
}

func (a Float64Access) Equals(float64OrFloat64Access interface{}) binaryExpression {
	return a.Compare("=", float64OrFloat64Access)
}

func (a Float64Access) Compare(op string, float64OrFloat64Access interface{}) binaryExpression {
	if !strings.Contains(validComparisonOperators, op) {
		panic("invalid comparison operator:" + op)
	}
	if f, ok := float64OrFloat64Access.(float64); ok {
		return MakeBinaryOperator(a, op, valuePrinter{f})
	}
	if ta, ok := float64OrFloat64Access.(Float64Access); ok {
		return MakeBinaryOperator(a, op, ta)
	}
	panic("float64 or Float64Access expected")
}
