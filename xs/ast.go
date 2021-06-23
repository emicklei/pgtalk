package xs

import "fmt"

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

func (o BinaryOperator) And(right SQLWriter) SQLWriter {
	return BinaryOperator{
		Left:     o,
		Operator: "AND",
		Right:    right,
	}
}
