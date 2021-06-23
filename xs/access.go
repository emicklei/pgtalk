package xs

import (
	"bytes"
	"fmt"
	"io"
)

type ReadWrite interface {
	Name() string
	Value(entity interface{}, fieldValue interface{})
}
type Int8Access struct {
	tableName string
	name      string
	writer    func(dest interface{}, i *int64)
}

func NewInt8Access(tableName string, columnName string, writer func(dest interface{}, i *int64)) Int8Access {
	return Int8Access{tableName: tableName, name: columnName, writer: writer}
}

func (a Int8Access) Value(entity interface{}, fieldValue interface{}) {
	var i int64 = fieldValue.(int64)
	a.writer(entity, &i)
}

func (a Int8Access) Equals(i int) BinaryOperator {
	return MakeBinaryOperator(a, "=", Printer{i})
}

func (a Int8Access) SQL(ctx SQLContext) string {
	alias := ctx.TableAlias(a.tableName)
	if alias != "" {
		return fmt.Sprintf("%s.%s", alias, a.name)
	}
	return a.name
}

type Printer struct {
	v interface{}
}

func (p Printer) SQL(ctx SQLContext) string { return fmt.Sprintf("%v", p.v) }

type ScanToWrite struct {
	RW     ReadWrite
	Entity interface{}
}

func (s ScanToWrite) Scan(fieldValue interface{}) error {
	s.RW.Value(s.Entity, fieldValue)
	return nil
}

func (a Int8Access) Name() string { return a.name }

type TextAccess struct {
	tableName string
	name      string
	writer    func(dest interface{}, i *string)
}

func (a TextAccess) Equals(s string) BinaryOperator {
	return MakeBinaryOperator(a, "=", LiteralString(s))
}

func NewTextAccess(tableName string, columnName string, writer func(dest interface{}, i *string)) TextAccess {
	return TextAccess{tableName: tableName, name: columnName, writer: writer}
}

func (a TextAccess) Value(entity interface{}, fieldValue interface{}) {
	var i string = fieldValue.(string)
	a.writer(entity, &i)
}

func (a TextAccess) Name() string { return a.name }

func (a TextAccess) SQL(ctx SQLContext) string {
	alias := ctx.TableAlias(a.tableName)
	if alias != "" {
		return fmt.Sprintf("%s.%s", alias, a.name)
	}
	return a.name
}

type LiteralString string

func (l LiteralString) SQL(ctx SQLContext) string {
	b := new(bytes.Buffer)
	io.WriteString(b, "'")
	io.WriteString(b, string(l))
	io.WriteString(b, "'")
	return b.String()
}

type NoCondition struct{}

var EmptyCondition = NoCondition{}

func (n NoCondition) SQL(ctx SQLContext) string { return "" }
