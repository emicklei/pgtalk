package pgtalk

import "fmt"

type QueryParameter struct {
	unimplementedBooleanExpression
	index int
	value any
}

func (a QueryParameter) SQLOn(w writeContext) {
	fmt.Fprintf(w, "$%d", a.index)
}

func argumentValues(list []QueryParameter) (values []any) {
	if len(list) == 0 {
		return
	}
	for _, each := range list {
		values = append(values, each.value)
	}
	return values
}
