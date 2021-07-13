package pgtalk

import (
	"fmt"
	"io"
	"time"
)

type TimeAccess struct {
	columnInfo
	fieldWriter func(dest interface{}, i *time.Time)
	insertValue time.Time
}

func NewTimeAccess(info TableInfo, columnName string, writer func(dest interface{}, i *time.Time)) TimeAccess {
	return TimeAccess{columnInfo: columnInfo{tableInfo: info, columnName: columnName}, fieldWriter: writer}
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
