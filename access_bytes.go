package pgtalk

import (
	"fmt"
	"io"
)

// BytesAccess can Read a column value (jsonb) and Write a column value and Set a struct field ([]byte).
type BytesAccess struct {
	ColumnInfo
	fieldWriter func(dest interface{}, b *[]byte)
	insertValue []byte
}

func NewBytesAccess(info ColumnInfo, writer func(dest interface{}, b *[]byte)) BytesAccess {
	return BytesAccess{ColumnInfo: info, fieldWriter: writer}
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

type JSONBAccess struct {
	ColumnInfo
	fieldWriter func(dest interface{}, b *string)
	insertValue string
}

func NewJSONBAccess(info ColumnInfo, writer func(dest interface{}, b *string)) JSONBAccess {
	return JSONBAccess{ColumnInfo: info, fieldWriter: writer}
}

func (a JSONBAccess) WriteInto(entity interface{}, fieldValue interface{}) {
	if fieldValue == nil {
		return
	}
	var f = fieldValue.([]byte)
	var s = string(f)
	a.fieldWriter(entity, &s)
}

func (a JSONBAccess) Set(s string) JSONBAccess {
	a.insertValue = s
	return a
}

func (a JSONBAccess) InsertValue() interface{} {
	return a.insertValue
}

func (a JSONBAccess) ValueAsSQLOn(w io.Writer) {
	fmt.Fprintf(w, "%v", a.insertValue) // TODO
}
