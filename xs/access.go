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
	ValueAsSQL() string
}

var EmptyReadWrite = []ReadWrite{}

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
