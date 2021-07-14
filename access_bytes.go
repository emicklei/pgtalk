package pgtalk

import (
	"fmt"
	"io"
)

// BytesAccess can Read a column value (jsonb) and Write a column value and Set a struct field ([]byte).
type BytesAccess struct {
	columnInfo
	fieldWriter func(dest interface{}, b *[]byte)
	insertValue []byte
}

func NewBytesAccess(info TableInfo, columnName string, writer func(dest interface{}, b *[]byte)) BytesAccess {
	return BytesAccess{columnInfo: makeColumnInfo(info, columnName), fieldWriter: writer}
}

func (a BytesAccess) WriteInto(entity interface{}, fieldValue interface{}) {
	if fieldValue == nil {
		return
	}
	var f = fieldValue.([]byte)
	a.fieldWriter(entity, &f)
}

func (a BytesAccess) ValueAsSQLOn(w io.Writer) {
	fmt.Fprintf(w, "%v", a.insertValue) // TODO
}

func (a BytesAccess) InsertValue() interface{} {
	return a.insertValue
}

func (a BytesAccess) Set(v []byte) BytesAccess {
	a.insertValue = v
	return a
}
