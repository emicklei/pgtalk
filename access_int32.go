package pgtalk

import (
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
)

// int32Access can Read a column value (int4) and Write a column value and Set a struct field (int32).
type int32Access struct {
	unimplementedBooleanExpression
	ColumnInfo
	fieldWriter   fieldAccessFunc
	valueToInsert int32
}

func NewInt32Access(
	info ColumnInfo,
	valueWriter func(dest any) any) int32Access {
	return int32Access{
		ColumnInfo:  info,
		fieldWriter: valueWriter}
}

func (a int32Access) BetweenAnd(begin int32, end int32) betweenAnd {
	return makeBetweenAnd(a, valuePrinter{v: begin}, valuePrinter{v: end})
}

func (a int32Access) FieldValueToScan(entity any) any {
	return a.fieldWriter(entity)
}

func (a int32Access) ValueToInsert() any {
	return a.valueToInsert
}

func (a int32Access) Set(v int32) int32Access {
	a.valueToInsert = v
	return a
}

func (a int32Access) In(values ...int32) SQLExpression {
	if len(values) == 0 {
		return makeConstantExpression(false)
	}
	vs := make([]any, len(values))
	for i := 0; i < len(values); i++ {
		vs[i] = values[i]
	}
	return makeBinaryOperator(a, "IN", valuesPrinter{vs: vs})
}

func (a int32Access) Equals(intLike any) binaryExpression {
	if i, ok := intLike.(int); ok {
		return makeBinaryOperator(a, "=", valuePrinter{v: i})
	}
	if ia, ok := intLike.(int32Access); ok {
		return makeBinaryOperator(a, "=", ia)
	}
	if p, ok := intLike.(*QueryParameter); ok {
		return makeBinaryOperator(a, "=", p)
	}
	if p, ok := intLike.(FieldAccess[pgtype.Int4]); ok {
		return makeBinaryOperator(a, "=", p)
	}
	panic(fmt.Sprintf("int or Int32Access or *QueryParameter expected, got %T", intLike))
}

func (a int32Access) Compare(op string, i int) binaryExpression {
	if !strings.Contains(validComparisonOperators, op) {
		panic("invalid comparison operator:" + op)
	}
	return makeBinaryOperator(a, op, valuePrinter{v: i})
}

func (a int32Access) Column() ColumnInfo { return a.ColumnInfo }

func (a int32Access) TableAlias(alias string) int32Access {
	a.ColumnInfo = a.ColumnInfo.TableAlias(alias)
	return a
}

// AppendScannable is part of ColumnAccessor
func (a int32Access) AppendScannable(list []any) []any {
	return append(list, &a.valueToInsert)
}

// Get returns the value for its columnName from a map (row).
func (a int32Access) Get(values map[string]any) any {
	v, ok := values[a.columnName]
	if !ok {
		return int32(0)
	}
	return v
}
