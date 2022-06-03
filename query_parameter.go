package pgtalk

import "fmt"

// QueryParameter captures any value as a parameter to use in a SQL query or mutation.
type QueryParameter struct {
	unimplementedBooleanExpression
	queryIndex int
	value      any
}

// SQLOn is part of SQLWriter
func (a QueryParameter) SQLOn(w WriteContext) {
	if a.queryIndex == 0 {
		fmt.Fprint(w, "?")
		return
	}
	fmt.Fprintf(w, "$%d", a.queryIndex)
}

// argumentValues returns the values for each parameter.
// it has the intended side-effect to update the query index of each parameter.
func argumentValues(list []*QueryParameter) (values []any) {
	if len(list) == 0 {
		return
	}
	for i, each := range list {
		each.queryIndex = i + 1 // starts at 1
		values = append(values, each.value)
	}
	return values
}

// TODO can we use type parameterization here?
// func NewParameter[T any](value T) *QueryParameter[T] { return &QueryParameter{value: value} }
func NewParameter(value any) *QueryParameter { return &QueryParameter{value: value} }
