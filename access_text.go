package pgtalk

import (
	"database/sql"
	"strings"
)

type TextAccess struct {
	ColumnInfo
	valueFieldWriter    func(dest interface{}, i string)
	nullableFieldWriter func(dest interface{}, i sql.NullString)
	valueToInsert       string
}

func NewTextAccess(info ColumnInfo, writer func(dest interface{}, i string), nullableWriter func(dest interface{}, i sql.NullString)) TextAccess {
	return TextAccess{ColumnInfo: info, nullableFieldWriter: nullableWriter, valueFieldWriter: writer}
}

func (a TextAccess) Set(v string) TextAccess {
	a.valueToInsert = v
	return a
}

func (a TextAccess) ValueToInsert() interface{} {
	return a.valueToInsert
}

func (a TextAccess) Equals(stringOrTextAccess interface{}) binaryExpression {
	return a.Compare("=", stringOrTextAccess)
}

func (a TextAccess) Compare(op string, stringOrTextAccess interface{}) binaryExpression {
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

// Collect is part of SQLExpression
func (a TextAccess) Collect(list []ColumnAccessor) []ColumnAccessor {
	return append(list, a)
}

func (a TextAccess) SetFieldValue(entity interface{}, fieldValue interface{}) error {
	if fieldValue == nil {
		return nil
	}
	s, ok := fieldValue.(string)
	if !ok {
		return NewValueConversionError(fieldValue, "string")
	}
	if a.notNull {
		a.valueFieldWriter(entity, s)
	} else {
		a.nullableFieldWriter(entity, &s)
	}
	return nil
}

func (a TextAccess) Like(pattern string) binaryExpression {
	return MakeBinaryOperator(a, "LIKE", LiteralString(pattern))
}

func (a TextAccess) In(values ...string) binaryExpression {
	vs := make([]interface{}, len(values))
	for i := 0; i < len(values); i++ {
		vs[i] = values[i]
	}
	return MakeBinaryOperator(a, "IN", valuesPrinter{vs})
}

func (a TextAccess) Column() ColumnInfo { return a.ColumnInfo }
