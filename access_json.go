package pgtalk

type jsonAccess struct {
	unimplementedBooleanExpression
	ColumnInfo
	valueFieldWriter fieldAccessFunc
	valueToInsert    any
}

func NewJSONAccess(info ColumnInfo,
	valueWriter func(dest any) any) jsonAccess {
	return jsonAccess{ColumnInfo: info, valueFieldWriter: valueWriter}
}

func (a jsonAccess) Set(jsonObject any) jsonAccess {
	a.valueToInsert = jsonObject
	return a
}

func (a jsonAccess) ValueToInsert() any {
	return a.valueToInsert
}

func (a jsonAccess) FieldValueToScan(entity any) any {
	return a.valueFieldWriter(entity)
}

func (a jsonAccess) Column() ColumnInfo { return a.ColumnInfo }

// Extract returns an expresion to get the path, extracted from the JSONB data, as a column
func (a jsonAccess) Extract(path string) SQLExpression {
	return makeBinaryOperator(a, "->", newLiteralString(path))
}

// TableAlias changes the table alias for this column accessor.
func (a jsonAccess) TableAlias(alias string) jsonAccess {
	a.ColumnInfo = a.ColumnInfo.TableAlias(alias)
	return a
}

// AppendScannable is part of ColumnAccessor
func (a jsonAccess) AppendScannable(list []any) []any {
	return append(list, &a.valueToInsert)
}

// Get returns the value for its columnName from a map (row).
func (a jsonAccess) Get(values map[string]any) any {
	v, ok := values[a.columnName]
	if !ok {
		return nil
	}
	return v
}
