package pgtalk

import (
	"strings"
)

// int64Access can Read a column value (int8) and Write a column value and Set a struct field (int64).
type int64Access struct {
	unimplementedBooleanExpression
	ColumnInfo
	fieldWriter   fieldAccessFunc
	valueToInsert int64
}

func NewInt64Access(
	info ColumnInfo,
	valueWriter func(dest any) any) int64Access {
	return int64Access{
		ColumnInfo:  info,
		fieldWriter: valueWriter}
}

func (a int64Access) BetweenAnd(begin int64, end int64) betweenAnd {
	return makeBetweenAnd(a, valuePrinter{v: begin}, valuePrinter{v: end})
}

func (a int64Access) FieldValueToScan(entity any) any {
	return a.fieldWriter(entity)
}

func (a int64Access) ValueToInsert() any {
	return a.valueToInsert
}

func (a int64Access) Set(v int64) int64Access {
	a.valueToInsert = v
	return a
}

func (a int64Access) Equals(intOrInt64Access any) binaryExpression {
	if i, ok := intOrInt64Access.(int); ok {
		return makeBinaryOperator(a, "=", valuePrinter{v: i})
	}
	if ia, ok := intOrInt64Access.(int64Access); ok {
		return makeBinaryOperator(a, "=", ia)
	}
	panic("int or Int64Access expected")
}

func (a int64Access) Compare(op string, i int) binaryExpression {
	if !strings.Contains(validComparisonOperators, op) {
		panic("invalid comparison operator:" + op)
	}
	return makeBinaryOperator(a, op, valuePrinter{v: i})
}

func (a int64Access) Column() ColumnInfo { return a.ColumnInfo }

func (a int64Access) TableAlias(alias string) int64Access {
	a.ColumnInfo = a.ColumnInfo.TableAlias(alias)
	return a
}

// AppendScannable is part of ColumnAccessor
func (a int64Access) AppendScannable(list []any) []any {
	return append(list, &a.valueToInsert)
}

// Get returns the value for its columnName from a map (row).
func (a int64Access) Get(values map[string]any) any {
	v, ok := values[a.columnName]
	if !ok {
		return int64(0)
	}
	return v
}
