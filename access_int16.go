package pgtalk

import (
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
)

// int16Access can Read a column value (int2) and Write a column value and Set a struct field (int16).
type int16Access struct {
	unimplementedBooleanExpression
	ColumnInfo
	fieldWriter   fieldAccessFunc
	valueToInsert int16
}

func NewInt16Access(
	info ColumnInfo,
	valueWriter func(dest any) any) int16Access {
	return int16Access{
		ColumnInfo:  info,
		fieldWriter: valueWriter}
}

func (a int16Access) BetweenAnd(begin int16, end int16) betweenAnd {
	return makeBetweenAnd(a, valuePrinter{v: begin}, valuePrinter{v: end})
}

func (a int16Access) FieldValueToScan(entity any) any {
	return a.fieldWriter(entity)
}

func (a int16Access) ValueToInsert() any {
	return a.valueToInsert
}

func (a int16Access) Set(v int16) int16Access {
	a.valueToInsert = v
	return a
}

func (a int16Access) In(values ...int16) SQLExpression {
	if len(values) == 0 {
		return makeConstantExpression(false)
	}
	vs := make([]any, len(values))
	for i := 0; i < len(values); i++ {
		vs[i] = values[i]
	}
	return makeBinaryOperator(a, "IN", valuesPrinter{vs: vs})
}

func (a int16Access) Equals(intLike any) SQLExpression {
	if i, ok := intLike.(int); ok {
		return makeBinaryOperator(a, "=", valuePrinter{v: i})
	}
	if i, ok := intLike.(int16); ok {
		return makeBinaryOperator(a, "=", valuePrinter{v: i})
	}
	if ia, ok := intLike.(int16Access); ok {
		return makeBinaryOperator(a, "=", ia)
	}
	if p, ok := intLike.(*QueryParameter); ok {
		return makeBinaryOperator(a, "=", p)
	}
	if p, ok := intLike.(FieldAccess[pgtype.Int2]); ok {
		return makeBinaryOperator(a, "=", p)
	}
	panic(fmt.Sprintf("int, int16, Int16Access or *QueryParameter expected, got %T", intLike))
}

func (a int16Access) Compare(op string, i int) binaryExpression {
	if !strings.Contains(validComparisonOperators, op) {
		panic("invalid comparison operator:" + op)
	}
	return makeBinaryOperator(a, op, valuePrinter{v: i})
}

func (a int16Access) Column() ColumnInfo { return a.ColumnInfo }

func (a int16Access) TableAlias(alias string) int16Access {
	a.ColumnInfo = a.ColumnInfo.TableAlias(alias)
	return a
}

// AppendScannable is part of ColumnAccessor
func (a int16Access) AppendScannable(list []any) []any {
	return append(list, &a.valueToInsert)
}

// Get returns the value for its columnName from a map (row).
func (a int16Access) Get(values map[string]any) any {
	v, ok := values[a.columnName]
	if !ok {
		return int16(0)
	}
	return v
}
