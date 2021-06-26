package xs

import "fmt"

// Int8Access can Read a column value (int8) and Write a column value and Set a struct field (int64).
type Int8Access struct {
	tableInfo   TableInfo
	columnName  string
	fieldWriter func(dest interface{}, i *int64)
	insertValue int64
}

func NewInt8Access(info TableInfo, columnName string, writer func(dest interface{}, i *int64)) Int8Access {
	return Int8Access{tableInfo: info, columnName: columnName, fieldWriter: writer}
}

func (a Int8Access) ValueAsSQL() string {
	return fmt.Sprintf("%d", a.insertValue)
}

func (a Int8Access) BetweenAnd(begin int64, end int64) BetweenAnd {
	return MakeBetweenAnd(a, Printer{begin}, Printer{end})
}

func (a Int8Access) WriteInto(entity interface{}, fieldValue interface{}) {
	var i int64 = fieldValue.(int64)
	a.fieldWriter(entity, &i)
}

func (a Int8Access) Value(v int64) Int8Access {
	return Int8Access{tableInfo: a.tableInfo, columnName: a.columnName, fieldWriter: a.fieldWriter, insertValue: v}
}

func (a Int8Access) Equals(i int) BinaryOperator {
	return MakeBinaryOperator(a, "=", Printer{i})
}

func (a Int8Access) SQL() string {
	return fmt.Sprintf("%s.%s", a.tableInfo.Alias, a.columnName)
}
func (a Int8Access) Name() string { return a.columnName }
