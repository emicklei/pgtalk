package pgtalk

type BooleanAccess struct {
	ColumnInfo
	valueFieldWriter FieldAccessFunc
	valueToInsert    bool
}

func NewBooleanAccess(info ColumnInfo, writer FieldAccessFunc) BooleanAccess {
	return BooleanAccess{ColumnInfo: info, valueFieldWriter: writer}
}

func (a BooleanAccess) Collect(list []ColumnAccessor) []ColumnAccessor {
	return append(list, a)
}

func (a BooleanAccess) Column() ColumnInfo { return a.ColumnInfo }

func (a BooleanAccess) Set(v bool) BooleanAccess {
	a.valueToInsert = v
	return a
}
func (a BooleanAccess) ValueToInsert() interface{} {
	return a.valueToInsert
}

func (a BooleanAccess) Equals(b bool) SQLExpression {
	return MakeBinaryOperator(a, "=", valuePrinter{b})
}

func (a BooleanAccess) FieldToScan(entity any) any {
	return a.valueFieldWriter(entity)
}
