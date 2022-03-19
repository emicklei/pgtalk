package pgtalk

import (
	"strings"
)

type TextAccess struct {
	ColumnInfo
	valueFieldWriter FieldAccessFunc
	valueToInsert    string
}

func NewTextAccess(info ColumnInfo, writer FieldAccessFunc) TextAccess {
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
		return MakeBinaryOperator(a, op, LiteralString(s))
	}
	if ta, ok := stringOrTextAccess.(TextAccess); ok {
		return MakeBinaryOperator(a, op, ta)
	}
	panic("string or TextAcces expected")
}

func (a TextAccess) FieldToScan(entity any) any {
	return a.valueFieldWriter(entity)
}

func (a TextAccess) Like(pattern string) binaryExpression {
	return MakeBinaryOperator(a, "LIKE", LiteralString(pattern))
}

func (a TextAccess) In(values ...string) binaryExpression {
	vs := make([]any, len(values))
	for i := 0; i < len(values); i++ {
		vs[i] = values[i]
	}
	return MakeBinaryOperator(a, "IN", valuesPrinter{vs})
}

func (a TextAccess) Column() ColumnInfo { return a.ColumnInfo }

func (a TextAccess) TableAlias(alias string) TextAccess {
	a.ColumnInfo = a.ColumnInfo.TableAlias(alias)
	return a
}
