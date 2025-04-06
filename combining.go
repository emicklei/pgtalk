package pgtalk

import (
	"fmt"
)

type QueryCombineable interface {
	SQLWriter
	Union(o QueryCombineable, all ...bool) QueryCombineable
	Intersect(o QueryCombineable, all ...bool) QueryCombineable
	Except(o QueryCombineable, all ...bool) QueryCombineable
}

type queryCombination struct {
	left     QueryCombineable
	operator string // UNION,INTERSECT,EXCEPT
	right    QueryCombineable
}

func combineOperator(combiner string, all ...bool) string {
	if len(all) > 0 && all[0] {
		return combiner + " ALL"
	}
	return combiner
}

func (q queryCombination) SQLOn(w WriteContext) {
	fmt.Fprint(w, "((")
	q.left.SQLOn(w)
	fmt.Fprintf(w, ") %s (", q.operator)
	q.right.SQLOn(w)
	fmt.Fprint(w, "))")
}

func (q queryCombination) Union(o QueryCombineable, all ...bool) QueryCombineable {
	return queryCombination{
		left:     q,
		operator: "UNION",
		right:    o,
	}
}
func (q queryCombination) Intersect(o QueryCombineable, all ...bool) QueryCombineable {
	return queryCombination{
		left:     q,
		operator: "INTERSECT",
		right:    o,
	}
}
func (q queryCombination) Except(o QueryCombineable, all ...bool) QueryCombineable {
	return queryCombination{
		left:     q,
		operator: "EXCEPT",
		right:    o,
	}
}
