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
	Left     SQLExpression
	Right    SQLExpression
}

func (o BinaryOperator) SQLOn(b io.Writer) {
	fmt.Fprint(b, "(")
	o.Left.SQLOn(b)
	fmt.Fprintf(b, " %s ", o.Operator)
	o.Right.SQLOn(b)
	fmt.Fprint(b, ")")
}

func MakeBinaryOperator(left SQLExpression, operator string, right SQLExpression) BinaryOperator {
	return BinaryOperator{
		Left:     left,
		Operator: operator,
		Right:    right,
	}
}

func (o BinaryOperator) And(right SQLExpression) BinaryOperator {
	return BinaryOperator{
		Left:     o,
		Operator: "AND",
		Right:    right,
	}
}

func (o BinaryOperator) Or(right SQLExpression) BinaryOperator {
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

func (o BinaryOperator) Collect(list []ColumnAccessor) []ColumnAccessor {
	return o.Left.Collect(o.Right.Collect(list))
}

type BetweenAnd struct {
}

func MakeBetweenAnd(reader ColumnAccessor, begin, end SQLExpression) BetweenAnd { return BetweenAnd{} }

type UnaryOperator struct {
	Operator string
	Operand  SQLExpression
}

func MakeUnaryOperator(operator string, operand SQLExpression) UnaryOperator {
	return UnaryOperator{Operator: operator, Operand: operand}
}

func (u UnaryOperator) SQLOn(w io.Writer) {
	fmt.Fprintf(w, "%s (", u.Operator)
	u.Operand.SQLOn(w)
	fmt.Fprint(w, ")")
}

func (u UnaryOperator) And(right SQLExpression) BinaryOperator {
	return BinaryOperator{
		Left:     u,
		Operator: "AND",
		Right:    right,
	}
}
func (u UnaryOperator) Or(right SQLExpression) BinaryOperator {
	return BinaryOperator{
		Left:     u,
		Operator: "OR",
		Right:    right,
	}
}

func (u UnaryOperator) Collect(list []ColumnAccessor) []ColumnAccessor {
	return u.Operand.Collect(list)
}

type NullCheck struct {
	Operand SQLExpression
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

func (n NullCheck) Collect(list []ColumnAccessor) []ColumnAccessor {
	return n.Operand.Collect(list)
}
