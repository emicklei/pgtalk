package pgtalk

import (
	"fmt"
	"io"
	"strings"
)

type TextAccess struct {
	columnInfo
	fieldWriter func(dest interface{}, i *string)
	insertValue string
}

func NewTextAccess(info TableInfo, columnName string, writer func(dest interface{}, i *string)) TextAccess {
	return TextAccess{columnInfo: columnInfo{tableInfo: info, columnName: columnName}, fieldWriter: writer}
}

func (a TextAccess) Set(v string) TextAccess {
	a.insertValue = v
	return a
}

func (a TextAccess) Equals(s string) BinaryOperator {
	return MakeBinaryOperator(a, "=", LiteralString(s))
}

func (a TextAccess) Compare(op string, s string) BinaryOperator {
	if !strings.Contains(validComparisonOperators, op) {
		panic("invalid comparison operator:" + op)
	}
	return MakeBinaryOperator(a, op, LiteralString(s))
}

func (a TextAccess) WriteInto(entity interface{}, fieldValue interface{}) {
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
