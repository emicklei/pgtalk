package pgtalk

import "fmt"

type TextAccess struct {
	tableInfo   TableInfo
	columnName  string
	fieldWriter func(dest interface{}, i *string)
	insertValue string
}

func NewTextAccess(info TableInfo, columnName string, writer func(dest interface{}, i *string)) TextAccess {
	return TextAccess{tableInfo: info, columnName: columnName, fieldWriter: writer}
}

func (a TextAccess) Value(v string) TextAccess {
	return TextAccess{tableInfo: a.tableInfo, columnName: a.columnName, fieldWriter: a.fieldWriter, insertValue: v}
}

func (a TextAccess) Equals(s string) BinaryOperator {
	return MakeBinaryOperator(a, "=", LiteralString(s))
}

func (a TextAccess) WriteInto(entity interface{}, fieldValue interface{}) {
	var i string = fieldValue.(string)
	a.fieldWriter(entity, &i)
}

func (a TextAccess) Name() string { return a.columnName }

func (a TextAccess) SQL() string {
	return fmt.Sprintf("%s.%s", a.tableInfo.Alias, a.columnName)
}

func (a TextAccess) ValueAsSQL() string {
	return fmt.Sprintf("'%s'", a.insertValue)
}
