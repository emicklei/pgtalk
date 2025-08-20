package pgtalk

import (
	"bytes"
	"io"
	"strings"
)

type WriteContext interface {
	Write(p []byte) (n int, err error)
	WithAlias(tableName, alias string) WriteContext
	TableAlias(tableName, defaultAlias string) string
}

type wc struct {
	writer  io.Writer
	aliases map[string]string
}

// NewWriteContext returns a new WriteContext to produce SQL
func NewWriteContext(w io.Writer) WriteContext {
	return wc{
		writer:  w,
		aliases: map[string]string{},
	}
}

// WithAlias returns a new context that knows about the table alias
func (w wc) WithAlias(tableName, alias string) WriteContext {
	cp := wc{
		writer:  w.writer,
		aliases: map[string]string{},
	}
	cp.aliases[tableName] = alias
	return cp
}

// Write is part of io.Writer
func (w wc) Write(p []byte) (n int, err error) {
	return w.writer.Write(p)
}

// TableAlias returns the alias for a tableName if present.
func (w wc) TableAlias(tableName, defaultAlias string) string {
	a, ok := w.aliases[tableName]
	if ok {
		return a
	}
	return defaultAlias
}

// onelineWriter is a writer that replaces newlines and tabs with spaces,
// and collapses multiple spaces into a single space.
type onelineWriter struct {
	b                *strings.Builder
	lastCharWasSpace bool
}

func newOnelineWriter(b *strings.Builder) *onelineWriter {
	return &onelineWriter{b: b}
}

func (w *onelineWriter) Write(p []byte) (n int, err error) {
	for _, b := range p {
		if b == '\n' || b == '\t' || b == ' ' {
			if !w.lastCharWasSpace {
				w.b.WriteByte(' ')
				w.lastCharWasSpace = true
			}
		} else {
			w.b.WriteByte(b)
			w.lastCharWasSpace = false
		}
	}
	return len(p), nil
}

// IndentedSQL returns source with tabs and lines trying to have a formatted view.
func IndentedSQL(some SQLWriter) string {
	buf := new(bytes.Buffer)
	some.SQLOn(NewWriteContext(buf))
	return buf.String()
}

// SQL returns source as a oneliner without tabs or line ends.
func SQL(some SQLWriter) string {
	var b strings.Builder
	w := newOnelineWriter(&b)
	some.SQLOn(NewWriteContext(w))
	return strings.TrimSpace(b.String())
}
