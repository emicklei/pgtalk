package pgtalk

// Work in Progress

type CollectSet struct {
	expressions []SQLWriter
	set         SQLWriter
}

func (c CollectSet) SQLOn(w WriteContext) {
	c.set.SQLOn(w)
}

func (d QuerySet[T]) Collect(expressions ...SQLWriter) CollectSet {
	return CollectSet{
		set:         d,
		expressions: expressions,
	}
}

func (a FieldAccess[T]) Concat(ex SQLExpression) binaryExpression {
	return binaryExpression{
		Left:     a,
		Operator: "||",
		Right:    ex,
	}
}
