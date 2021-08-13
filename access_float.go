package pgtalk

import (
	"fmt"
	"io"
)

// Float64Access can Read a column value (float) and Write a column value and Set a struct field (float64).
type Float64Access struct {
	ColumnInfo
	fieldWriter   func(dest interface{}, f *float64)
	valueToInsert float64
}

func NewFloat64Access(info ColumnInfo, writer func(dest interface{}, f *float64)) Float64Access {
	return Float64Access{ColumnInfo: info, fieldWriter: writer}
}

func (a Float64Access) SetFieldValue(entity interface{}, fieldValue interface{}) {
	if fieldValue == nil {
		return
	}
	var f = fieldValue.(float64)
	a.fieldWriter(entity, &f)
}

func (a Float64Access) ValueToInsert() interface{} {
	return a.ValueToInsert
}

func (a Float64Access) ValueAsSQLOn(w io.Writer) {
	fmt.Fprintf(w, "%v", a.valueToInsert)
}

func (a Float64Access) Column() ColumnInfo { return a.ColumnInfo }
