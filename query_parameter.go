package pgtalk

import "fmt"

type QueryParameter struct {
	unimplementedBooleanExpression
	index      int
	setIndex   int
	queryIndex int
	value      any
}

func (a QueryParameter) SQLOn(w writeContext) {
	if a.index == 0 {
		fmt.Fprint(w, "?")
		return
	}
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

type ParameterSet struct {
	parameters []QueryParameter
}

func NewParameterSet() *ParameterSet { return &ParameterSet{} }

func (s *ParameterSet) NewParameter(value any) QueryParameter {
	p := QueryParameter{value: value, setIndex: len(s.parameters) + 1} // start at 1
	s.parameters = append(s.parameters, p)
	return p
}

func (s *ParameterSet) Parameters() []QueryParameter { return s.parameters }
