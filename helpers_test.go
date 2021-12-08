package pgtalk

import (
	"bytes"
	"fmt"
	"io"
	"time"

	"github.com/jackc/pgtype"
)

var (
	polyTable = TableInfo{Name: "polies", Schema: "public", Alias: "p1"}
	polyFTime = NewTimeAccess(MakeColumnInfo(polyTable, "ftime", NotPrimary, Nullable, 1),
		func(dest interface{}, v *time.Time) { dest.(*poly).FTime = v })
	polyFFloat = NewFloat64Access(MakeColumnInfo(polyTable, "ffloat", NotPrimary, Nullable, 1),
		func(dest interface{}, v *float64) { dest.(*poly).FFloat = v })
	polyFUUID = NewFieldAccess[pgtype.UUID](MakeColumnInfo(polyTable, "fuuid", NotPrimary, Nullable, 1),
		func(dest interface{}, v *pgtype.UUID) { dest.(*poly).FUUID = v })
	polyColumns = append([]ColumnAccessor{}, polyFTime, polyFFloat)
	polyAccess  = TableAccessor{TableInfo: polyTable, AllColumns: polyColumns}
)

type poly struct {
	FTime  *time.Time
	FFloat *float64
	FBool  *bool
	// pgtypes
	FUUID *pgtype.UUID
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
