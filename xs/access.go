package xs

import (
	"bytes"
	"fmt"
	"io"
)

type TableInfo struct {
	Name  string
	Alias string
}

type ReadWrite interface {
	Name() string
	// temp name
	WriteInto(entity interface{}, fieldValue interface{})
	// temp name
	SetSQL() string
}

var EmptyReadWrite = []ReadWrite{}

type Int8Access struct {
	tableInfo   TableInfo
	name        string
	writer      func(dest interface{}, i *int64)
	insertValue int64
}

func NewInt8Access(info TableInfo, columnName string, writer func(dest interface{}, i *int64)) Int8Access {
	return Int8Access{tableInfo: info, name: columnName, writer: writer}
}

func (a Int8Access) SetSQL() string {
	return fmt.Sprintf("%d", a.insertValue)
}

func (a Int8Access) BetweenAnd(begin int64, end int64) BetweenAnd {
	return MakeBetweenAnd(a, Printer{begin}, Printer{end})
}

func (a Int8Access) WriteInto(entity interface{}, fieldValue interface{}) {
	var i int64 = fieldValue.(int64)
	a.writer(entity, &i)
}

func (a Int8Access) Value(v int64) Int8Access {
	return Int8Access{tableInfo: a.tableInfo, name: a.name, writer: a.writer, insertValue: v}
}

func (a Int8Access) Equals(i int) BinaryOperator {
	return MakeBinaryOperator(a, "=", Printer{i})
}

func (a Int8Access) SQL() string {
	return fmt.Sprintf("%s.%s", a.tableInfo.Alias, a.name)
}

type Printer struct {
	v interface{}
}

func (p Printer) SQL() string { return fmt.Sprintf("%v", p.v) }

type ScanToWrite struct {
	RW     ReadWrite
	Entity interface{}
}

func (s ScanToWrite) Scan(fieldValue interface{}) error {
	s.RW.WriteInto(s.Entity, fieldValue)
	return nil
}

func (a Int8Access) Name() string { return a.name }

type TextAccess struct {
	tableInfo   TableInfo
	name        string
	writer      func(dest interface{}, i *string)
	insertValue string
}

func (a TextAccess) Value(v string) TextAccess {
	return TextAccess{tableInfo: a.tableInfo, name: a.name, writer: a.writer, insertValue: v}
}

func (a TextAccess) Equals(s string) BinaryOperator {
	return MakeBinaryOperator(a, "=", LiteralString(s))
}

func NewTextAccess(info TableInfo, columnName string, writer func(dest interface{}, i *string)) TextAccess {
	return TextAccess{tableInfo: info, name: columnName, writer: writer}
}

func (a TextAccess) WriteInto(entity interface{}, fieldValue interface{}) {
	var i string = fieldValue.(string)
	a.writer(entity, &i)
}

func (a TextAccess) Name() string { return a.name }

func (a TextAccess) SQL() string {
	return fmt.Sprintf("%s.%s", a.tableInfo.Alias, a.name)
}

func (a TextAccess) SetSQL() string {
	return fmt.Sprintf("'%s'", a.insertValue)
}

type LiteralString string

func (l LiteralString) SQL() string {
	b := new(bytes.Buffer)
	io.WriteString(b, "'")
	io.WriteString(b, string(l))
	io.WriteString(b, "'")
	return b.String()
}

type NoCondition struct{}

var EmptyCondition = NoCondition{}

func (n NoCondition) SQL() string { return "" }
