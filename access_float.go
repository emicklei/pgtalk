package pgtalk

import (
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
)

// float64Access can Read a column value (float) and Write a column value and Set a struct field (float64).
type float64Access struct {
	unimplementedBooleanExpression
	ColumnInfo
	fieldWriter         fieldAccessFunc
	nullableFieldWriter func(dest any, f pgtype.Float8)
	valueToInsert       float64
}

func NewFloat64Access(info ColumnInfo, writer fieldAccessFunc) float64Access {
	return float64Access{ColumnInfo: info, fieldWriter: writer}
}

func (a float64Access) ValueToInsert() any {
	return a.ValueToInsert
}

func (a float64Access) Set(v float64) float64Access {
	a.valueToInsert = v
	return a
}

func (a float64Access) Column() ColumnInfo { return a.ColumnInfo }

func (a float64Access) Equals(isFloatLike any) binaryExpression {
	return a.Compare("=", isFloatLike)
}

func (a float64Access) Compare(op string, isFloatLike any) binaryExpression {
	if !strings.Contains(validComparisonOperators, op) {
		panic("invalid comparison operator:" + op)
	}
	if f, ok := isFloatLike.(float64); ok {
		return makeBinaryOperator(a, op, valuePrinter{v: f})
	}
	if ta, ok := isFloatLike.(float64Access); ok {
		return makeBinaryOperator(a, op, ta)
	}
	if qa, ok := isFloatLike.(*QueryParameter); ok {
		if _, ok := qa.value.(float64); ok {
			return makeBinaryOperator(a, op, qa)
		}
	}
	panic("float64, Float64Access or *QueryArgument[float64] expected")
}

func (a float64Access) FieldValueToScan(entity any) any {
	return a.fieldWriter(entity)
}

func (a float64Access) TableAlias(alias string) float64Access {
	a.ColumnInfo = a.ColumnInfo.TableAlias(alias)
	return a
}

// AppendScannable is part of ColumnAccessor
func (a float64Access) AppendScannable(list []any) []any {
	return append(list, &a.valueToInsert)
}

// Get returns the value for its columnName from a map (row).
func (a float64Access) Get(values map[string]any) any {
	v, ok := values[a.columnName]
	if !ok {
		return float64(0.0)
	}
	return v
}
