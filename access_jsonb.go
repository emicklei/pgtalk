package pgtalk

type jsonBAccess struct {
	unimplementedBooleanExpression
	ColumnInfo
	valueFieldWriter fieldAccessFunc
	valueToInsert    []byte
}

func NewJSONBAccess(info ColumnInfo,
	valueWriter func(dest any) any) jsonBAccess {
	return jsonBAccess{ColumnInfo: info, valueFieldWriter: valueWriter}
}

func (a jsonBAccess) Set(s []byte) jsonBAccess {
	a.valueToInsert = s
	return a
}

func (a jsonBAccess) ValueToInsert() any {
	return a.valueToInsert
}

func (a jsonBAccess) FieldValueToScan(entity any) any {
	return a.valueFieldWriter(entity)
}

func (a jsonBAccess) Column() ColumnInfo { return a.ColumnInfo }

// Extract returns an expresion to get the path, extracted from the JSONB data, as a column
func (a jsonBAccess) Extract(path string) SQLExpression {
	return makeBinaryOperator(a, "->", newLiteralString(path))
}

func (a jsonBAccess) TableAlias(alias string) jsonBAccess {
	a.ColumnInfo = a.ColumnInfo.TableAlias(alias)
	return a
}

// AppendScannable is part of ColumnAccessor
func (a jsonBAccess) AppendScannable(list []any) []any {
	return append(list, &a.valueToInsert)
}

// Get returns the value for its columnName from a map (row).
func (a jsonBAccess) Get(values map[string]any) any {
	v, ok := values[a.columnName]
	if !ok {
		return []byte{}
	}
	return v
}
