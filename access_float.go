package pgtalk

import (
	"fmt"
	"io"
)

// Float64Access can Read a column value (float) and Write a column value and Set a struct field (float64).
type Float64Access struct {
	columnInfo
	fieldWriter func(dest interface{}, f *float64)
	insertValue float64
}

func NewFloat64Access(info TableInfo, columnName string, writer func(dest interface{}, f *float64)) Float64Access {
	return Float64Access{columnInfo: makeColumnInfo(info, columnName), fieldWriter: writer}
}

func (a Float64Access) WriteInto(entity interface{}, fieldValue interface{}) {
	if fieldValue == nil {
		return
	}
	var f = fieldValue.(float64)
	a.fieldWriter(entity, &f)
}

func (a Float64Access) InsertValue() interface{} {
	return a.insertValue
}

func (a Float64Access) ValueAsSQLOn(w io.Writer) {
	fmt.Fprintf(w, "%v", a.insertValue)
}
