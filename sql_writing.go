package pgtalk

import (
	"bytes"
	"io"
)

type writeContext struct {
	writer  io.Writer
	aliases map[string]string
}

func newWriteContext(w io.Writer) writeContext {
	return writeContext{
		writer:  w,
		aliases: map[string]string{},
	}
}

func (w writeContext) WithAlias(tableName, alias string) writeContext {
	cp := newWriteContext(w.writer)
	cp.aliases[tableName] = alias
	return cp
}

// Write is part of io.Writer
func (w writeContext) Write(p []byte) (n int, err error) {
	return w.writer.Write(p)
}

func (w writeContext) TableAlias(tableName, defaultAlias string) string {
	a, ok := w.aliases[tableName]
	if ok {
		return a
	}
	return defaultAlias
}

func SQL(some SQLWriter) string {
	buf := new(bytes.Buffer)
	some.SQLOn(newWriteContext(buf))
	return buf.String()
}
