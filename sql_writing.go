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

// IndentedSQL returns source with tabs and lines trying to have a formatted view.
func IndentedSQL(some SQLWriter) string {
	buf := new(bytes.Buffer)
	some.SQLOn(NewWriteContext(buf))
	return buf.String()
}

// SQL returns source as a oneliner without tabs or line ends.
func SQL(some SQLWriter) string {
	src := IndentedSQL(some)
	return strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(src, "\t", " "), "\n", " "), "  ", " ")
}

type lazySQL struct {
	w SQLWriter
}

func (l lazySQL) String() string {
	return SQL(l.w)
}

// LazySQL is a helper to produce an SQL string when asked for String().
// It is used to avoid the SQL string being computed for logs when not needed.
func LazySQL(w SQLWriter) lazySQL { return lazySQL{w: w} }
