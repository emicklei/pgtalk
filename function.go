package pgtalk

import "io"

// SQLFunction is for calling any Postgres standard function with zero or more arguments
// as part of your query.
// Typically used together with pgtalk.SQLAs.
type SQLFunction struct {
	Name      string
	Arguments []SQLExpression
}

// SQLOn is part of SQLWriter
func (f SQLFunction) SQLOn(w WriteContext) {
	io.WriteString(w, f.Name)
	io.WriteString(w, "(")
	for i, each := range f.Arguments {
		if i > 0 {
			io.WriteString(w, ",")
		}
		each.SQLOn(w)
	}
	io.WriteString(w, ")")
}

// NewSQLFunction creates a new SQLFunction value.
func NewSQLFunction(name string, arguments ...SQLExpression) SQLFunction {
	return SQLFunction{Name: name, Arguments: arguments}
}

func NewSQLConstant(value any) SQLExpression {
	return valuePrinter{v: value}
}
