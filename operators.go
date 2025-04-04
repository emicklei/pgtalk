package pgtalk

import (
	"fmt"
	"strings"
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
	// inline comparison operators
	if strings.Contains(validComparisonOperators, o.Operator) {
		fmt.Fprintf(w, " ")
	} else {
		fmt.Fprintf(w, "\n\t")
	}
	fmt.Fprintf(w, "%s ", o.Operator)
	o.Right.SQLOn(w)
	fmt.Fprint(w, ")")
}

func makeBinaryOperator(left SQLExpression, operator string, right SQLExpression) binaryExpression {
	return binaryExpression{
		Left:     left,
		Operator: operator,
		Right:    right,
	}
}

func (o binaryExpression) And(right SQLExpression) SQLExpression {
	return binaryExpression{
		Left:     o,
		Operator: "AND",
		Right:    right,
	}
}

func (o binaryExpression) Or(right SQLExpression) SQLExpression {
	return binaryExpression{
		Left:     o,
		Operator: "OR",
		Right:    right,
	}
}

func (o binaryExpression) Like(pattern string) SQLExpression {
	return binaryExpression{
		Left:     o,
		Operator: "LIKE",
		Right:    valuePrinter{v: pattern},
	}
}

type betweenAnd struct {
}

func makeBetweenAnd(_ ColumnAccessor, _, _ SQLExpression) betweenAnd { return betweenAnd{} }

type unaryExpression struct {
	Operator string
	Operand  SQLExpression
}

func (u unaryExpression) SQLOn(w WriteContext) {
	fmt.Fprintf(w, "%s (", u.Operator)
	u.Operand.SQLOn(w)
	fmt.Fprint(w, ")")
}

func (u unaryExpression) And(right SQLExpression) SQLExpression {
	return binaryExpression{
		Left:     u,
		Operator: "AND",
		Right:    right,
	}
}
func (u unaryExpression) Or(right SQLExpression) SQLExpression {
	return binaryExpression{
		Left:     u,
		Operator: "OR",
		Right:    right,
	}
}

type nullCheck struct {
	Operand SQLExpression
	// IsNot == true -> IS NOT NULL
	IsNot bool
}

func (n nullCheck) SQLOn(w WriteContext) {
	fmt.Fprint(w, "(")
	n.Operand.SQLOn(w)
	if n.IsNot {
		fmt.Fprint(w, " IS NOT NULL)")
		return
	}
	fmt.Fprint(w, " IS NULL)")
}

func (n nullCheck) And(right SQLExpression) SQLExpression {
	return makeBinaryOperator(n, "AND", right)
}

func (n nullCheck) Or(right SQLExpression) SQLExpression {
	return makeBinaryOperator(n, "OR", right)
}

// IsNotNull returns an expression with the IS NOT NULL condition
func IsNotNull(e SQLExpression) nullCheck {
	return nullCheck{Operand: e, IsNot: true}
}

// IsNull returns an expression with the IS NULL condition
func IsNull(e SQLExpression) nullCheck {
	return nullCheck{Operand: e, IsNot: false}
}

type constantExpression struct {
	v any
}

func makeConstantExpression(v any) constantExpression {
	return constantExpression{v: v}
}
func (c constantExpression) SQLOn(w WriteContext) {
	makeValuePrinter(c.v).SQLOn(w)
}
func (c constantExpression) Or(right SQLExpression) SQLExpression {
	return makeBinaryOperator(c, "OR", right)
}
func (c constantExpression) And(right SQLExpression) SQLExpression {
	return makeBinaryOperator(c, "AND", right)
}

type sqlExpression struct {
	sql string
}

func NewSQLSource(sql string) SQLExpression {
	return sqlExpression{sql: sql}
}
func (s sqlExpression) SQLOn(w WriteContext) {
	fmt.Fprint(w, s.sql)
}
func (s sqlExpression) Or(right SQLExpression) SQLExpression {
	return makeBinaryOperator(s, "OR", right)
}
func (s sqlExpression) And(right SQLExpression) SQLExpression {
	return makeBinaryOperator(s, "AND", right)
}
