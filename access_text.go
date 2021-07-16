package pgtalk

import (
	"fmt"
	"io"
	"strings"
)

type TextAccess struct {
	ColumnInfo
	fieldWriter func(dest interface{}, i *string)
	insertValue string
}

func NewTextAccess(info ColumnInfo, writer func(dest interface{}, i *string)) TextAccess {
	return TextAccess{ColumnInfo: info, fieldWriter: writer}
}

func (a TextAccess) Set(v string) TextAccess {
	a.insertValue = v
	return a
}

func (a TextAccess) InsertValue() interface{} {
	return a.insertValue
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

func (a TextAccess) WriteInto(entity interface{}, fieldValue interface{}) {
	if fieldValue == nil {
		return
	}
	var i string = fieldValue.(string)
	a.fieldWriter(entity, &i)
}

func (a TextAccess) ValueAsSQLOn(w io.Writer) {
	fmt.Fprintf(w, "'%s'", a.insertValue)
}

func (a TextAccess) NotNull() NullCheck {
	return NullCheck{Operand: a, IsNot: true}
}

func (a TextAccess) Like(pattern string) BinaryOperator {
	return MakeBinaryOperator(a, "LIKE", ValuePrinter{pattern})
}

func (a TextAccess) In(values ...string) BinaryOperator {
	vs := make([]interface{}, len(values))
	for i := 0; i < len(values); i++ {
		vs[i] = values[i]
	}
	return MakeBinaryOperator(a, "IN", ValuesPrinter{vs})
}
