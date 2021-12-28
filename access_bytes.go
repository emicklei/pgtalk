package pgtalk

import "database/sql"

// BytesAccess can Read a column value (jsonb) and Write a column value and Set a struct field ([]byte).
type BytesAccess struct {
	ColumnInfo
	valueFieldWriter    func(dest interface{}, b string)
	nullableFieldWriter func(dest interface{}, b sql.NullString)
	valueToInsert       []byte
}

func NewBytesAccess(info ColumnInfo,
	valueWriter func(dest interface{}, b string),
	nullableValueWriter func(dest interface{}, b sql.NullString)) BytesAccess {
	return BytesAccess{ColumnInfo: info,
		valueFieldWriter:    valueWriter,
		nullableFieldWriter: nullableValueWriter,
	}
}

func (a BytesAccess) Column() ColumnInfo { return a.ColumnInfo }

func (a BytesAccess) SetFieldValue(entity interface{}, fieldValue interface{}) error { return nil }

// Collect is part of SQLExpression
func (a BytesAccess) Collect(list []ColumnAccessor) []ColumnAccessor {
	return append(list, a)
}

func (a BytesAccess) WriteInto(entity interface{}, fieldValue interface{}) {
	if fieldValue == nil {
		return
	}
	var bytesValue = fieldValue.([]byte)
	if a.notNull {
		a.valueFieldWriter(entity, string(bytesValue))
	} else {
		a.nullableFieldWriter(entity, sql.NullString{String: string(bytesValue), Valid: true})
	}
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
	valueFieldWriter    func(dest interface{}, b string)
	nullableFieldWriter func(dest interface{}, b sql.NullString)
	valueToInsert       string
}

func NewJSONBAccess(info ColumnInfo,
	valueWriter func(dest interface{}, b string),
	nullableValueWriter func(dest interface{}, b sql.NullString)) JSONBAccess {
	return JSONBAccess{ColumnInfo: info, valueFieldWriter: valueWriter, nullableFieldWriter: nullableValueWriter}
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
	if a.notNull {
		a.valueFieldWriter(entity, s)
	} else {
		a.nullableFieldWriter(entity, sql.NullString{String: s, Valid: true})
	}
	return nil
}

func (a JSONBAccess) Set(s string) JSONBAccess {
	a.valueToInsert = s
	return a
}

func (a JSONBAccess) ValueToInsert() interface{} {
	return a.valueToInsert
}

func (a JSONBAccess) Column() ColumnInfo { return a.ColumnInfo }

// Collect is part of SQLExpression
func (a JSONBAccess) Collect(list []ColumnAccessor) []ColumnAccessor {
	return append(list, a)
}

// Extract returns an expresion to get the path, extracted from the JSONB data, as a column
func (a JSONBAccess) Extract(path string) SQLExpression {
	return MakeBinaryOperator(a, "->", LiteralString(path))
}
