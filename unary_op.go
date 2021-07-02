package pgtalk

import (
	"fmt"
	"io"
)

type UnaryOperator struct {
	Operator string
	Operand  SQLWriter
}

func MakeUnaryOperator(operator string, operand SQLWriter) UnaryOperator {
	return UnaryOperator{Operator: operator, Operand: operand}
}

func (u UnaryOperator) SQLOn(w io.Writer) {
	fmt.Fprintf(w, "%s (", u.Operator)
	u.Operand.SQLOn(w)
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

func (n NullCheck) SQLOn(w io.Writer) {
	fmt.Fprint(w, "(")
	n.Operand.SQLOn(w)
	if n.IsNot {
		fmt.Fprint(w, " IS NOT NULL)")
		return
	}
	fmt.Fprint(w, " IS NULL)")
}
