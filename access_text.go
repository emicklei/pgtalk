package pgtalk

import (
	"strings"

	"github.com/jackc/pgtype"
)

type TextAccess struct {
	unimplementedBooleanExpression
	ColumnInfo
	valueFieldWriter fieldAccessFunc
	valueToInsert    string
}

func NewTextAccess(info ColumnInfo, writer fieldAccessFunc) TextAccess {
	return TextAccess{ColumnInfo: info, valueFieldWriter: writer}
}

func (a TextAccess) Set(v string) TextAccess {
	a.valueToInsert = v
	return a
}

func (a TextAccess) ValueToInsert() any {
	return a.valueToInsert
}

func (a TextAccess) Equals(stringOrTextAccess any) binaryExpression {
	return a.Compare("=", stringOrTextAccess)
}

func (a TextAccess) Compare(op string, stringOrTextAccess any) binaryExpression {
	if !strings.Contains(validComparisonOperators, op) {
		panic("invalid comparison operator:" + op)
	}
	if s, ok := stringOrTextAccess.(string); ok {
		return makeBinaryOperator(a, op, newLiteralString(s))
	}
	if ta, ok := stringOrTextAccess.(TextAccess); ok {
		return makeBinaryOperator(a, op, ta)
	}
	panic("string or TextAcces expected")
}

func (a TextAccess) FieldValueToScan(entity any) any {
	return a.valueFieldWriter(entity)
}

func (a TextAccess) Like(pattern string) binaryExpression {
	return makeBinaryOperator(a, "LIKE", newLiteralString(pattern))
}

func (a TextAccess) In(values ...string) binaryExpression {
	vs := make([]any, len(values))
	for i := 0; i < len(values); i++ {
		vs[i] = values[i]
	}
	return makeBinaryOperator(a, "IN", valuesPrinter{vs: vs})
}

func (a TextAccess) Column() ColumnInfo { return a.ColumnInfo }

func (a TextAccess) TableAlias(alias string) TextAccess {
	a.ColumnInfo = a.ColumnInfo.TableAlias(alias)
	return a
}

// AppendScannable is part of ColumnAccessor
func (a TextAccess) AppendScannable(list []any) []any {
	return append(list, &a.valueToInsert)
}

// Get returns the value for its columnName from a map (row).
func (a TextAccess) Get(values map[string]any) any {
	v, ok := values[a.columnName]
	if !ok {
		return pgtype.Text{}
	}
	return v
}
