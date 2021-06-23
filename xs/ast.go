package xs

import "fmt"

type SQLContext struct {
}

func (s SQLContext) TableAlias(tableName string) string {
	return "t1"
}

type SQLWriter interface {
	SQL(ctx SQLContext) string
}

type BinaryOperator struct {
	Operator string
	Left     SQLWriter
	Right    SQLWriter
}

func (o BinaryOperator) SQL(ctx SQLContext) string {
	return fmt.Sprintf("(%s %s %s)", o.Left.SQL(ctx), o.Operator, o.Right.SQL(ctx))
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
