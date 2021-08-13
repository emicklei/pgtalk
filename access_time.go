package pgtalk

import (
	"fmt"
	"io"
	"time"
)

type TimeAccess struct {
	ColumnInfo
	fieldWriter   func(dest interface{}, i *time.Time)
	valueToInsert time.Time
}

func NewTimeAccess(info ColumnInfo, writer func(dest interface{}, i *time.Time)) TimeAccess {
	return TimeAccess{ColumnInfo: info, fieldWriter: writer}
}

func (a TimeAccess) SetFieldValue(entity interface{}, fieldValue interface{}) {
	if fieldValue == nil {
		return
	}
	var v time.Time = fieldValue.(time.Time)
	a.fieldWriter(entity, &v)
}

func (a TimeAccess) ValueAsSQLOn(w io.Writer) {
	fmt.Fprintf(w, "%v", a.ValueToInsert)
}

func (a TimeAccess) ValueToInsert() interface{} {
	return a.valueToInsert
}

func (a TimeAccess) Set(v time.Time) TimeAccess {
	a.valueToInsert = v
	return a
}

func (a TimeAccess) Column() ColumnInfo { return a.ColumnInfo }
