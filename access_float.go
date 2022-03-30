package pgtalk

import (
	"strings"

	"github.com/jackc/pgtype"
)

// Float64Access can Read a column value (float) and Write a column value and Set a struct field (float64).
type Float64Access struct {
	unimplementedBooleanExpression
	ColumnInfo
	fieldWriter         fieldAccessFunc
	nullableFieldWriter func(dest any, f pgtype.Float8)
	valueToInsert       float64
}

func NewFloat64Access(info ColumnInfo, writer fieldAccessFunc) Float64Access {
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
		return makeBinaryOperator(a, op, valuePrinter{v: f})
	}
	if ta, ok := float64OrFloat64Access.(Float64Access); ok {
		return makeBinaryOperator(a, op, ta)
	}
	panic("float64 or Float64Access expected")
}

func (a Float64Access) FieldValueToScan(entity any) any {
	return a.fieldWriter(entity)
}

func (a Float64Access) TableAlias(alias string) Float64Access {
	a.ColumnInfo = a.ColumnInfo.TableAlias(alias)
	return a
}

// AppendScannable is part of ColumnAccessor
func (a Float64Access) AppendScannable(list []any) []any {
	return append(list, &a.valueToInsert)
}

// Get returns the value for its columnName from a map (row).
func (a Float64Access) Get(values map[string]any) any {
	v, ok := values[a.columnName]
	if !ok {
		return float64(0.0)
	}
	return v
}
