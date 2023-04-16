package pgtalk

type booleanAccess struct {
	unimplementedBooleanExpression
	ColumnInfo
	valueFieldWriter fieldAccessFunc
	valueToInsert    bool
}

func NewBooleanAccess(info ColumnInfo, writer fieldAccessFunc) booleanAccess {
	return booleanAccess{ColumnInfo: info, valueFieldWriter: writer}
}

func (a booleanAccess) Column() ColumnInfo { return a.ColumnInfo }

func (a booleanAccess) Set(v bool) booleanAccess {
	a.valueToInsert = v
	return a
}
func (a booleanAccess) ValueToInsert() any {
	return a.valueToInsert
}

func (a booleanAccess) And(e SQLExpression) SQLExpression {
	return makeBinaryOperator(a, "AND", e)
}

func (a booleanAccess) Equals(b bool) SQLExpression {
	return makeBinaryOperator(a, "=", valuePrinter{v: b})
}

func (a booleanAccess) FieldValueToScan(entity any) any {
	return a.valueFieldWriter(entity)
}

// TableAlias changes the table alias for this column accessor.
func (a booleanAccess) TableAlias(alias string) booleanAccess {
	a.ColumnInfo = a.ColumnInfo.TableAlias(alias)
	return a
}

// AppendScannable is part of ColumnAccessor
func (a booleanAccess) AppendScannable(list []any) []any {
	return append(list, &a.valueToInsert)
}

// Get returns the value for its columnName from a map (row).
func (a booleanAccess) Get(values map[string]any) any {
	v, ok := values[a.columnName]
	if !ok {
		return false
	}
	return v
}
