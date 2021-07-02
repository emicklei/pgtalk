package pgtalk

import (
	"fmt"
	"io"

	"github.com/jackc/pgtype"
)

type DateAccess struct {
	columnInfo
	fieldWriter func(dest interface{}, i *pgtype.Date)
	insertValue pgtype.Date
}

func NewDateAccess(info TableInfo, columnName string, writer func(dest interface{}, i *pgtype.Date)) DateAccess {
	return DateAccess{columnInfo: columnInfo{tableInfo: info, columnName: columnName}, fieldWriter: writer}
}

func (a DateAccess) WriteInto(entity interface{}, fieldValue interface{}) {
	var v pgtype.Date = fieldValue.(pgtype.Date)
	a.fieldWriter(entity, &v)
}

func (a DateAccess) ValueAsSQLOn(w io.Writer) {
	fmt.Fprintf(w, "%v", a.insertValue)
}

func NewTimestampAccess(info TableInfo, columnName string, writer func(dest interface{}, i *pgtype.Date)) DateAccess {
	return NewDateAccess(info, columnName, writer)
}
