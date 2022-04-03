package pgtalk

import "fmt"

type QueryArgument struct {
	unimplementedBooleanExpression
	index int
	value any
}

func (a QueryArgument) SQLOn(w writeContext) {
	fmt.Fprintf(w, "$%d", a.index)
}

func argumentValues(list []QueryArgument) (values []any) {
	if len(list) == 0 {
		return
	}
	for _, each := range list {
		values = append(values, each.value)
	}
	return values
}
