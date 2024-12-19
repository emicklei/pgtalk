package pgtalk

import (
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
)

type textAccess struct {
	unimplementedBooleanExpression
	ColumnInfo
	valueFieldWriter fieldAccessFunc
	valueToInsert    string
}

func NewTextAccess(info ColumnInfo, writer fieldAccessFunc) textAccess {
	return textAccess{ColumnInfo: info, valueFieldWriter: writer}
}

func (a textAccess) Set(v string) textAccess {
	a.valueToInsert = v
	return a
}

func (a textAccess) ValueToInsert() any {
	return a.valueToInsert
}

func (a textAccess) Equals(textLike any) binaryExpression {
	return a.Compare("=", textLike)
}

func (a textAccess) Compare(op string, stringOrTextAccess any) binaryExpression {
	if !strings.Contains(validComparisonOperators, op) {
		panic("invalid comparison operator:" + op)
	}
	if s, ok := stringOrTextAccess.(string); ok {
		return makeBinaryOperator(a, op, newLiteralString(s))
	}
	if ta, ok := stringOrTextAccess.(textAccess); ok {
		return makeBinaryOperator(a, op, ta)
	}
	if p, ok := stringOrTextAccess.(*QueryParameter); ok {
		return makeBinaryOperator(a, op, p)
	}
	panic("string or TextAccess expected")
}

func (a textAccess) FieldValueToScan(entity any) any {
	return a.valueFieldWriter(entity)
}

func (a textAccess) Like(pattern string) binaryExpression {
	return makeBinaryOperator(a, "LIKE", newLiteralString(pattern))
}

// ILIKE can be used instead of LIKE to make the match case-insensitive according to the active locale.
func (a textAccess) ILike(pattern string) binaryExpression {
	return makeBinaryOperator(a, "ILIKE", newLiteralString(pattern))
}

func (a textAccess) In(values ...string) SQLExpression {
	if len(values) == 0 {
		return makeConstantExpression(false)
	}
	vs := make([]any, len(values))
	for i := 0; i < len(values); i++ {
		vs[i] = values[i]
	}
	return makeBinaryOperator(a, "IN", valuesPrinter{vs: vs})
}

func (a textAccess) Column() ColumnInfo { return a.ColumnInfo }

// TableAlias changes the table alias for this column accessor.
func (a textAccess) TableAlias(alias string) textAccess {
	a.ColumnInfo = a.ColumnInfo.TableAlias(alias)
	return a
}

// AppendScannable is part of ColumnAccessor
func (a textAccess) AppendScannable(list []any) []any {
	return append(list, &a.valueToInsert)
}

// Get returns the value for its columnName from a map (row).
func (a textAccess) Get(values map[string]any) any {
	v, ok := values[a.columnName]
	if !ok {
		return pgtype.Text{}
	}
	return v
}
