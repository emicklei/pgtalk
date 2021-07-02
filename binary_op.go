package pgtalk

import (
	"fmt"
	"io"
)

const (
	validComparisonOperators = "= > < >= <= <>"
)

type BinaryOperator struct {
	Operator string
	Left     SQLWriter
	Right    SQLWriter
}

func (o BinaryOperator) SQLOn(b io.Writer) {
	fmt.Fprint(b, "(")
	o.Left.SQLOn(b)
	fmt.Fprintf(b, " %s ", o.Operator)
	o.Right.SQLOn(b)
	fmt.Fprint(b, ")")
}

func MakeBinaryOperator(left SQLWriter, operator string, right SQLWriter) BinaryOperator {
	return BinaryOperator{
		Left:     left,
		Operator: operator,
		Right:    right,
	}
}

func (o BinaryOperator) And(right SQLWriter) BinaryOperator {
	return BinaryOperator{
		Left:     o,
		Operator: "AND",
		Right:    right,
	}
}

func (o BinaryOperator) Or(right SQLWriter) BinaryOperator {
	return BinaryOperator{
		Left:     o,
		Operator: "OR",
		Right:    right,
	}
}

func (o BinaryOperator) Like(pattern string) BinaryOperator {
	return BinaryOperator{
		Left:     o,
		Operator: "LIKE",
		Right:    ValuePrinter{pattern},
	}
}

type BetweenAnd struct {
}

func MakeBetweenAnd(reader ColumnAccessor, begin, end SQLWriter) BetweenAnd { return BetweenAnd{} }
