package pgtalk

import (
	"fmt"
)

const (
	validComparisonOperators = "= > < >= <= <>"
)

type binaryExpression struct {
	Operator string
	Left     SQLExpression
	Right    SQLExpression
}

func (o binaryExpression) SQLOn(w WriteContext) {
	fmt.Fprint(w, "(")
	o.Left.SQLOn(w)
	fmt.Fprintf(w, " %s ", o.Operator)
	o.Right.SQLOn(w)
	fmt.Fprint(w, ")")
}

func MakeBinaryOperator(left SQLExpression, operator string, right SQLExpression) binaryExpression {
	return binaryExpression{
		Left:     left,
		Operator: operator,
		Right:    right,
	}
}

func (o binaryExpression) And(right SQLExpression) binaryExpression {
	return binaryExpression{
		Left:     o,
		Operator: "AND",
		Right:    right,
	}
}

func (o binaryExpression) Or(right SQLExpression) binaryExpression {
	return binaryExpression{
		Left:     o,
		Operator: "OR",
		Right:    right,
	}
}

func (o binaryExpression) Like(pattern string) binaryExpression {
	return binaryExpression{
		Left:     o,
		Operator: "LIKE",
		Right:    valuePrinter{pattern},
	}
}

type BetweenAnd struct {
}

func MakeBetweenAnd(reader ColumnAccessor, begin, end SQLExpression) BetweenAnd { return BetweenAnd{} }

type unaryExpression struct {
	Operator string
	Operand  SQLExpression
}

func MakeUnaryOperator(operator string, operand SQLExpression) unaryExpression {
	return unaryExpression{Operator: operator, Operand: operand}
}

func (u unaryExpression) SQLOn(w WriteContext) {
	fmt.Fprintf(w, "%s (", u.Operator)
	u.Operand.SQLOn(w)
	fmt.Fprint(w, ")")
}

func (u unaryExpression) And(right SQLExpression) binaryExpression {
	return binaryExpression{
		Left:     u,
		Operator: "AND",
		Right:    right,
	}
}
func (u unaryExpression) Or(right SQLExpression) binaryExpression {
	return binaryExpression{
		Left:     u,
		Operator: "OR",
		Right:    right,
	}
}

type NullCheck struct {
	Operand SQLExpression
	// IsNot == true -> IS NOT NULL
	IsNot bool
}

func (n NullCheck) SQLOn(w WriteContext) {
	fmt.Fprint(w, "(")
	n.Operand.SQLOn(w)
	if n.IsNot {
		fmt.Fprint(w, " IS NOT NULL)")
		return
	}
	fmt.Fprint(w, " IS NULL)")
}

// IsNotNull returns an expression with the IS NOT NULL condition
func IsNotNull(e SQLExpression) NullCheck {
	return NullCheck{Operand: e, IsNot: true}
}

// IsNull returns an expression with the IS NULL condition
func IsNull(e SQLExpression) NullCheck {
	return NullCheck{Operand: e, IsNot: false}
}
