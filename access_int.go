package pgtalk

import (
	"strings"
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

func (a Int64Access) FieldValueToScan(entity any) any {
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

// AppendScannable is part of ColumnAccessor
func (a Int64Access) AppendScannable(list []any) []any {
	return append(list, &a.valueToInsert)
}

// Get returns the value for its columnName from a map (row).
func (a Int64Access) Get(values map[string]any) (any, bool) {
	v, ok := values[a.columnName]
	return v, ok
}
