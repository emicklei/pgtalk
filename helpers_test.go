package pgtalk

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgproto3/v2"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
)

var (
	polyTable = TableInfo{Name: "polies", Schema: "public", Alias: "p1"}
	polyFTime = NewTimeAccess(MakeColumnInfo(polyTable, "ftime", NotPrimary, Nullable, 1),
		func(dest any) any { return &dest.(*poly).FTime })
	polyFFloat = NewFloat64Access(MakeColumnInfo(polyTable, "ffloat", NotPrimary, Nullable, 1),
		func(dest any) any { return &dest.(*poly).FFloat })
	polyFUUID = NewFieldAccess[pgtype.UUID](MakeColumnInfo(polyTable, "fuuid", NotPrimary, Nullable, 1),
		func(dest any) any { return &dest.(*poly).FUUID })
	polyFString = NewTextAccess(MakeColumnInfo(polyTable, "fstring", NotPrimary, Nullable, 1),
		func(dest any) any { return &dest.(*poly).FString })
	polyColumns = append([]ColumnAccessor{}, polyFTime, polyFFloat)
)

type poly struct {
	FTime   time.Time
	FFloat  float64
	FBool   bool
	FString string
	// pgtypes
	FUUID pgtype.UUID
	// for storing custom field expression result values
	expressionResults map[string]any
}

func diff(left, right string) string {
	//assume one line
	b := new(bytes.Buffer)
	io.WriteString(b, "\n")
	io.WriteString(b, left)
	io.WriteString(b, "\n")
	leftRunes := []rune(left)
	rightRunes := []rune(right)
	size := len(leftRunes)
	if l := len(rightRunes); l < size {
		size = l
	}
	for c := 0; c < size; c++ {
		l := leftRunes[c]
		r := rightRunes[c]
		if l == r {
			b.WriteRune(l)
		} else {
			fmt.Fprintf(b, "^(%s)...", string(r))
			break
		}
	}
	return b.String()
}

func newMockConnection(t *testing.T) *fakeConnection {
	return &fakeConnection{t: t}
}

type fakeConnection struct {
	t *testing.T
}

type fakeRows struct {
	t *testing.T
}

func (f fakeRows) Close()                        {}
func (f fakeRows) Err() error                    { return nil }
func (f fakeRows) CommandTag() pgconn.CommandTag { return pgconn.CommandTag{} }
func (f fakeRows) FieldDescriptions() []pgproto3.FieldDescription {
	return []pgproto3.FieldDescription{}
}
func (f fakeRows) Next() bool { return false }
func (f fakeRows) Scan(dest ...interface{}) error {
	f.t.Helper()
	f.t.Log("destinations:", dest)
	return nil
}
func (f fakeRows) Values() ([]interface{}, error) { return []any{}, nil }
func (f fakeRows) RawValues() [][]byte            { return [][]byte{} }

func (f *fakeConnection) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	f.t.Helper()
	f.t.Log("sql:", oneliner(sql))
	f.t.Log("parameters:", args)
	return fakeRows{f.t}, nil
}

func (f *fakeConnection) ctx() context.Context { return context.Background() }
