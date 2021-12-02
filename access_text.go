package pgtalk

import (
	"fmt"
	"io"
	"strings"
)

type TextAccess struct {
	ColumnInfo
	fieldWriter   func(dest interface{}, i *string)
	valueToInsert string
}

func NewTextAccess(info ColumnInfo, writer func(dest interface{}, i *string)) TextAccess {
	return TextAccess{ColumnInfo: info, fieldWriter: writer}
}

func (a TextAccess) Set(v string) TextAccess {
	a.valueToInsert = v
	return a
}

func (a TextAccess) ValueToInsert() interface{} {
	return a.valueToInsert
}

func (a TextAccess) Equals(stringOrTextAccess interface{}) BinaryOperator {
	return a.Compare("=", stringOrTextAccess)
}

func (a TextAccess) Compare(op string, stringOrTextAccess interface{}) BinaryOperator {
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
	a.fieldWriter(entity, &s)
	return nil
}

func (a TextAccess) ValueAsSQLOn(w io.Writer) {
	fmt.Fprintf(w, "'%s'", a.valueToInsert)
}

func (a TextAccess) NotNull() NullCheck {
	return NullCheck{Operand: a, IsNot: true}
}

func (a TextAccess) Like(pattern string) BinaryOperator {
	return MakeBinaryOperator(a, "LIKE", LiteralString(pattern))
}

func (a TextAccess) In(values ...string) BinaryOperator {
	vs := make([]interface{}, len(values))
	for i := 0; i < len(values); i++ {
		vs[i] = values[i]
	}
	return MakeBinaryOperator(a, "IN", ValuesPrinter{vs})
}

func (a TextAccess) Column() ColumnInfo { return a.ColumnInfo }

func (a TextAccess) String() string {
	return fmt.Sprintf("text(%v)", a.ColumnInfo)
}
