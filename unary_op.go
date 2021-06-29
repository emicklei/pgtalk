package pgtalk

import "fmt"

type UnaryOperator struct {
	Operator string
	Operand  SQLWriter
}

func MakeUnaryOperator(operator string, operand SQLWriter) UnaryOperator {
	return UnaryOperator{Operator: operator, Operand: operand}
}

func (u UnaryOperator) SQL() string {
	return fmt.Sprintf("%s (%s)", u.Operator, u.Operand.SQL())
}

func (u UnaryOperator) And(right SQLWriter) BinaryOperator {
	return BinaryOperator{
		Left:     u,
		Operator: "AND",
		Right:    right,
	}
}
func (u UnaryOperator) Or(right SQLWriter) BinaryOperator {
	return BinaryOperator{
		Left:     u,
		Operator: "OR",
		Right:    right,
	}
}

type NullCheck struct {
	Operand SQLWriter
	// IsNot == true -> IS NOT NULL
	IsNot bool
}

func (n NullCheck) SQL() string {
	if n.IsNot {
		return fmt.Sprintf("(%s IS NOT NULL)", n.Operand.SQL())
	}
	return fmt.Sprintf("(%s IS NULL)", n.Operand.SQL())
}
