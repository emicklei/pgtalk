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

func (a BytesAccess) FieldValueToScan(entity any) any {
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

// AppendScannable is part of ColumnAccessor
func (a BytesAccess) AppendScannable(list []any) []any {
	return append(list, &a.valueToInsert)
}

// Get returns the value for its columnName from a map (row).
func (a BytesAccess) Get(values map[string]any) (any, bool) {
	v, ok := values[a.columnName]
	return v, ok
}
