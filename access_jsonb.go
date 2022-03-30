package pgtalk

type JSONBAccess struct {
	unimplementedBooleanExpression
	ColumnInfo
	valueFieldWriter fieldAccessFunc
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

func (a JSONBAccess) FieldValueToScan(entity any) any {
	return a.valueFieldWriter(entity)
}

func (a JSONBAccess) Column() ColumnInfo { return a.ColumnInfo }

// Extract returns an expresion to get the path, extracted from the JSONB data, as a column
func (a JSONBAccess) Extract(path string) SQLExpression {
	return makeBinaryOperator(a, "->", newLiteralString(path))
}

func (a JSONBAccess) TableAlias(alias string) JSONBAccess {
	a.ColumnInfo = a.ColumnInfo.TableAlias(alias)
	return a
}

// AppendScannable is part of ColumnAccessor
func (a JSONBAccess) AppendScannable(list []any) []any {
	return append(list, &a.valueToInsert)
}

// Get returns the value for its columnName from a map (row).
func (a JSONBAccess) Get(values map[string]any) any {
	v, ok := values[a.columnName]
	if !ok {
		return []byte{}
	}
	return v
}
