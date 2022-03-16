package pgtalk

import (
	"bytes"
	"io"
)

type WriteContext struct {
	writer  io.Writer
	aliases map[string]string
}

func NewWriteContext(w io.Writer) WriteContext {
	return WriteContext{
		writer:  w,
		aliases: map[string]string{},
	}
}

func (w WriteContext) WithAlias(tableName, alias string) WriteContext {
	w.aliases[tableName] = alias
	return w
}

// Write is part of io.Writer
func (w WriteContext) Write(p []byte) (n int, err error) {
	return w.writer.Write(p)
}

func (w WriteContext) TableAlias(tableName, defaultAlias string) string {
	a, ok := w.aliases[tableName]
	if ok {
		return a
	}
	return defaultAlias
}

func SQL(some SQLWriter) string {
	buf := new(bytes.Buffer)
	some.SQLOn(NewWriteContext(buf)) // info is ignored
	return buf.String()
}
