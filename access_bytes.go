package pgtalk

// BytesAccess can Read a column value (jsonb) and Write a column value and Set a struct field ([]byte).
type BytesAccess struct {
	ColumnInfo
	valueFieldWriter func(dest any) *string
	valueToInsert    []byte
}

func NewBytesAccess(info ColumnInfo,
	valueWriter func(dest any) *string) BytesAccess {
	return BytesAccess{ColumnInfo: info,
		valueFieldWriter: valueWriter,
	}
}

func (a BytesAccess) Column() ColumnInfo { return a.ColumnInfo }

func (a BytesAccess) FieldToScan(entity any) any {
	return a.valueFieldWriter(entity)
}

func (a BytesAccess) ValueToInsert() any {
	return a.valueToInsert
}

func (a BytesAccess) Set(v []byte) BytesAccess {
	a.valueToInsert = v
	return a
}

func (a BytesAccess) TableAlias(alias string) BytesAccess {
	a.ColumnInfo = a.ColumnInfo.TableAlias(alias)
	return a
}

type JSONBAccess struct {
	ColumnInfo
	valueFieldWriter FieldAccessFunc
	valueToInsert    []byte
}

func NewJSONBAccess(info ColumnInfo,
	valueWriter func(dest any) any) JSONBAccess {
	return JSONBAccess{ColumnInfo: info, valueFieldWriter: valueWriter}
}

func (a JSONBAccess) Set(s []byte) JSONBAccess {
	a.valueToInsert = s
	return a
}

func (a JSONBAccess) ValueToInsert() any {
	return a.valueToInsert
}

func (a JSONBAccess) FieldToScan(entity any) any {
	return a.valueFieldWriter(entity)
}

func (a JSONBAccess) Column() ColumnInfo { return a.ColumnInfo }

// Extract returns an expresion to get the path, extracted from the JSONB data, as a column
func (a JSONBAccess) Extract(path string) SQLExpression {
	return MakeBinaryOperator(a, "->", LiteralString(path))
}

func (a JSONBAccess) TableAlias(alias string) JSONBAccess {
	a.ColumnInfo = a.ColumnInfo.TableAlias(alias)
	return a
}
