package pgtalk

type BooleanAccess struct {
	ColumnInfo
	valueFieldWriter FieldAccessFunc
	valueToInsert    bool
}

func NewBooleanAccess(info ColumnInfo, writer FieldAccessFunc) BooleanAccess {
	return BooleanAccess{ColumnInfo: info, valueFieldWriter: writer}
}

func (a BooleanAccess) Column() ColumnInfo { return a.ColumnInfo }

func (a BooleanAccess) Set(v bool) BooleanAccess {
	a.valueToInsert = v
	return a
}
func (a BooleanAccess) ValueToInsert() any {
	return a.valueToInsert
}

func (a BooleanAccess) Equals(b bool) SQLExpression {
	return MakeBinaryOperator(a, "=", valuePrinter{b})
}

func (a BooleanAccess) FieldToScan(entity any) any {
	return a.valueFieldWriter(entity)
}

func (a BooleanAccess) TableAlias(alias string) BooleanAccess {
	a.ColumnInfo = a.ColumnInfo.TableAlias(alias)
	return a
}
