package pgtalk

import (
	"fmt"
	"io"
)

// BytesAccess can Read a column value (jsonb) and Write a column value and Set a struct field ([]byte).
type BytesAccess struct {
	ColumnInfo
	fieldWriter   func(dest interface{}, b *[]byte)
	valueToInsert []byte
}

func NewBytesAccess(info ColumnInfo, writer func(dest interface{}, b *[]byte)) BytesAccess {
	return BytesAccess{ColumnInfo: info, fieldWriter: writer}
}

func (a BytesAccess) Column() ColumnInfo { return a.ColumnInfo }

func (a BytesAccess) SetFieldValue(entity interface{}, fieldValue interface{}) error { return nil }

func (a BytesAccess) Collect(list []ColumnAccessor) []ColumnAccessor {
	return append(list, a)
}

func (a BytesAccess) WriteInto(entity interface{}, fieldValue interface{}) {
	if fieldValue == nil {
		return
	}
	var f = fieldValue.([]byte)
	a.fieldWriter(entity, &f)
}

func (a BytesAccess) ValueAsSQLOn(w io.Writer) {
	fmt.Fprintf(w, "%v", a.valueToInsert) // TODO
}

func (a BytesAccess) ValueToInsert() interface{} {
	return a.valueToInsert
}

func (a BytesAccess) Set(v []byte) BytesAccess {
	a.valueToInsert = v
	return a
}

type JSONBAccess struct {
	ColumnInfo
	fieldWriter   func(dest interface{}, b *string)
	valueToInsert string
}

func NewJSONBAccess(info ColumnInfo, writer func(dest interface{}, b *string)) JSONBAccess {
	return JSONBAccess{ColumnInfo: info, fieldWriter: writer}
}

func (a JSONBAccess) SetFieldValue(entity interface{}, fieldValue interface{}) error {
	if fieldValue == nil {
		return nil
	}
	f, ok := fieldValue.([]byte)
	if !ok {
		// TODO try string?
		return NewValueConversionError(fieldValue, "[]byte")
	}
	var s = string(f)
	a.fieldWriter(entity, &s)
	return nil
}

func (a JSONBAccess) Set(s string) JSONBAccess {
	a.valueToInsert = s
	return a
}

func (a JSONBAccess) ValueToInsert() interface{} {
	return a.valueToInsert
}

func (a JSONBAccess) ValueAsSQLOn(w io.Writer) {
	fmt.Fprintf(w, "%v", a.valueToInsert) // TODO
}

func (a JSONBAccess) Column() ColumnInfo { return a.ColumnInfo }

func (a JSONBAccess) Collect(list []ColumnAccessor) []ColumnAccessor {
	return append(list, a)
}

// Extract returns an expresion to get the path, extracted from the JSONB data, as a column
func (a JSONBAccess) Extract(path string) SQLExpression {
	return MakeBinaryOperator(a, "->", LiteralString(path))
}
