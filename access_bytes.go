package pgtalk

// bytesAccess can Read a column value (jsonb) and Write a column value and Set a struct field ([]byte).
type bytesAccess struct {
	ColumnInfo
	valueFieldWriter func(dest any) *string
	valueToInsert    []byte
}

func NewBytesAccess(info ColumnInfo,
	valueWriter func(dest any) *string) bytesAccess {
	return bytesAccess{ColumnInfo: info,
		valueFieldWriter: valueWriter,
	}
}

func (a bytesAccess) Column() ColumnInfo { return a.ColumnInfo }

func (a bytesAccess) FieldValueToScan(entity any) any {
	return a.valueFieldWriter(entity)
}

func (a bytesAccess) ValueToInsert() any {
	return a.valueToInsert
}

func (a bytesAccess) Set(v []byte) bytesAccess {
	a.valueToInsert = v
	return a
}

// TableAlias changes the table alias for this column accessor.
func (a bytesAccess) TableAlias(alias string) bytesAccess {
	a.ColumnInfo = a.ColumnInfo.TableAlias(alias)
	return a
}

// AppendScannable is part of ColumnAccessor
func (a bytesAccess) AppendScannable(list []any) []any {
	return append(list, &a.valueToInsert)
}

// Get returns the value for its columnName from a map (row).
func (a bytesAccess) Get(values map[string]any) any {
	v, ok := values[a.columnName]
	if !ok {
		return []byte{}
	}
	return v
}
