package convert

import (
	"encoding/hex"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func UUID(v uuid.UUID) pgtype.UUID {
	return pgtype.UUID{
		Bytes: v,
		Valid: true,
	}
}

// StringToUUID converts a string UUID in standard form to a pgtype.UUID.
// Check for Valid before using the result.
func StringToUUID(s string) pgtype.UUID {
	data, err := parseUUID(s)
	if err != nil {
		return pgtype.UUID{
			Bytes: [16]byte{},
			Valid: false,
		}
	}
	return pgtype.UUID{
		Bytes: data,
		Valid: true,
	}
}

// parseUUID converts a string UUID in standard form to a byte array.
func parseUUID(src string) (dst [16]byte, err error) {
	switch len(src) {
	case 36:
		src = src[0:8] + src[9:13] + src[14:18] + src[19:23] + src[24:]
	case 32:
		// dashes already stripped, assume valid
	default:
		// assume invalid.
		return dst, fmt.Errorf("cannot parse UUID %v", src)
	}

	buf, err := hex.DecodeString(src)
	if err != nil {
		return dst, err
	}

	copy(dst[:], buf)
	return dst, err
}

// UUIDToString returns format xxxx-yyyy-zzzz-rrrr-tttt
func UUIDToString(t pgtype.UUID) string {
	if !t.Valid {
		return ""
	}
	src := t.Bytes
	return fmt.Sprintf("%x-%x-%x-%x-%x", src[0:4], src[4:6], src[6:8], src[8:10], src[10:16])
}

func TimeToTimestamptz(t time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{Time: t, Valid: true}
}

func TimeToTimestamp(t time.Time) pgtype.Timestamp {
	return pgtype.Timestamp{Time: t.UTC(), Valid: true}
}

func TimeToDate(t time.Time) pgtype.Date {
	return pgtype.Date{Time: t, Valid: true}
}

func StringToText(s string) pgtype.Text {
	return pgtype.Text{String: s, Valid: true}
}

func Int64ToInt8(i int64) pgtype.Int8 {
	return pgtype.Int8{Int64: i, Valid: true}
}

func Int4ToInt8(i pgtype.Int4) pgtype.Int8 {
	return pgtype.Int8{Int64: int64(i.Int32), Valid: true}
}

func Int32ToInt4(i int32) pgtype.Int4 {
	return pgtype.Int4{Int32: i, Valid: true}
}

func Int8(i int) pgtype.Int8 {
	return Int64ToInt8(int64(i))
}

func Int4(i int) pgtype.Int4 {
	return Int32ToInt4(int32(i))
}

// func ByteSliceToJSONB(d []byte) pgtype.JSONBCodec {
// 	return pgtype.JSONB{Bytes: d, Valid: true}
// }

func Bool(b bool) pgtype.Bool {
	return pgtype.Bool{Bool: b, Valid: true}
}

func DateToTimePtr(d pgtype.Date) *time.Time {
	if d.Valid {
		t := d.Time
		return &t
	}
	return nil
}

func Float64ToFloat8(f float64) pgtype.Float8 {
	return pgtype.Float8{
		Valid:   true,
		Float64: f,
	}
}

func StringsToTextArray(list []string) pgtype.FlatArray[pgtype.Text] {
	a := pgtype.FlatArray[pgtype.Text]{}
	for _, each := range list {
		a = append(a, StringToText(each))
	}
	return a
}

func TextArrayToStrings(array pgtype.FlatArray[pgtype.Text]) []string {
	list := []string{}
	for _, each := range array {
		list = append(list, each.String)
	}
	return list
}
