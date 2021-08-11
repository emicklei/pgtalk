package pgtalk

import (
	"fmt"
	"io"
	"time"
)

type TimeAccess struct {
	ColumnInfo
	fieldWriter func(dest interface{}, i *time.Time)
	insertValue time.Time
}

func NewTimeAccess(info ColumnInfo, writer func(dest interface{}, i *time.Time)) TimeAccess {
	return TimeAccess{ColumnInfo: info, fieldWriter: writer}
}

func (a TimeAccess) WriteInto(entity interface{}, fieldValue interface{}) {
	if fieldValue == nil {
		return
	}
	var v time.Time = fieldValue.(time.Time)
	a.fieldWriter(entity, &v)
}

func (a TimeAccess) ValueAsSQLOn(w io.Writer) {
	fmt.Fprintf(w, "%v", a.insertValue)
}

func (a TimeAccess) InsertValue() interface{} {
	return a.insertValue
}

func (a TimeAccess) Set(v time.Time) TimeAccess {
	a.insertValue = v
	return a
}

func (a TimeAccess) Column() ColumnInfo { return a.ColumnInfo }
