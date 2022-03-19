package pgtalk

import (
	"strings"

	"github.com/jackc/pgtype"
)

// Int64Access can Read a column value (int8) and Write a column value and Set a struct field (int64).
type Int64Access struct {
	ColumnInfo
	fieldWriter   FieldAccessFunc
	valueToInsert int64
}

func NewInt64Access(
	info ColumnInfo,
	valueWriter func(dest any) any) Int64Access {
	return Int64Access{
		ColumnInfo:  info,
		fieldWriter: valueWriter}
}

func (a Int64Access) BetweenAnd(begin int64, end int64) BetweenAnd {
	return MakeBetweenAnd(a, valuePrinter{begin}, valuePrinter{end})
}

func (a Int64Access) FieldToScan(entity any) any {
	return a.fieldWriter(entity)
}

func (a Int64Access) ValueToInsert() any {
	return a.valueToInsert
}

func (a Int64Access) Set(v int64) Int64Access {
	a.valueToInsert = v
	return a
}

func (a Int64Access) Equals(intOrInt64Access any) binaryExpression {
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

func (a Int64Access) TableAlias(alias string) Int64Access {
	a.ColumnInfo = a.ColumnInfo.TableAlias(alias)
	return a
}

// Float64Access can Read a column value (float) and Write a column value and Set a struct field (float64).
type Float64Access struct {
	ColumnInfo
	fieldWriter         FieldAccessFunc
	nullableFieldWriter func(dest any, f pgtype.Float8)
	valueToInsert       float64
}

func NewFloat64Access(info ColumnInfo, writer FieldAccessFunc) Float64Access {
	return Float64Access{ColumnInfo: info, fieldWriter: writer}
}

func (a Float64Access) ValueToInsert() any {
	return a.ValueToInsert
}

func (a Float64Access) Set(v float64) Float64Access {
	a.valueToInsert = v
	return a
}

func (a Float64Access) Column() ColumnInfo { return a.ColumnInfo }

func (a Float64Access) Equals(float64OrFloat64Access any) binaryExpression {
	return a.Compare("=", float64OrFloat64Access)
}

func (a Float64Access) Compare(op string, float64OrFloat64Access any) binaryExpression {
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

func (a Float64Access) FieldToScan(entity any) any {
	return a.fieldWriter(entity)
}

func (a Float64Access) TableAlias(alias string) Float64Access {
	a.ColumnInfo = a.ColumnInfo.TableAlias(alias)
	return a
}
