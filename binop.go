package pgtalk

import "fmt"

const (
	validComparisonOperators = "= > < >= <= <>"
)

type SQLWriter interface {
	SQL() string
}

type BinaryOperator struct {
	Operator string
	Left     SQLWriter
	Right    SQLWriter
}

func (o BinaryOperator) SQL() string {
	return fmt.Sprintf("(%s %s %s)", o.Left.SQL(), o.Operator, o.Right.SQL())
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
