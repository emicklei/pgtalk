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

func (a TimeAccess) SetFieldValue(entity interface{}, fieldValue interface{}) error {
	if fieldValue == nil {
		return nil
	}
	v, ok := fieldValue.(time.Time)
	if !ok {
		return NewValueConversionError(fieldValue, "time.Time")
	}
	a.fieldWriter(entity, &v)
	return nil
}

func (a TimeAccess) ValueAsSQLOn(w io.Writer) {
	fmt.Fprintf(w, "%v", a.valueToInsert)
}

func (a TimeAccess) ValueToInsert() interface{} {
	return a.valueToInsert
}

func (a TimeAccess) Set(v time.Time) TimeAccess {
	a.valueToInsert = v
	return a
}

func (a TimeAccess) Column() ColumnInfo { return a.ColumnInfo }
