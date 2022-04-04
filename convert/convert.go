package convert

import (
	"encoding/hex"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgtype"
)

func UUID(v uuid.UUID) pgtype.UUID {
	return pgtype.UUID{
		Bytes:  v,
		Status: pgtype.Present,
	}
}

func StringToUUID(s string) pgtype.UUID {
	data, err := parseUUID(s)
	if err != nil {
		panic(err)
	}
	return pgtype.UUID{
		Bytes:  data,
		Status: pgtype.Present,
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
	if t.Status != pgtype.Present {
		return ""
	}
	src := t.Bytes
	return fmt.Sprintf("%x-%x-%x-%x-%x", src[0:4], src[4:6], src[6:8], src[8:10], src[10:16])
}

func TimeToTimestamptz(t time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{Time: t, Status: pgtype.Present}
}

func TimeToTimestamp(t time.Time) pgtype.Timestamp {
	return pgtype.Timestamp{Time: t.UTC(), Status: pgtype.Present}
}

func TimeToDate(t time.Time) pgtype.Date {
	return pgtype.Date{Time: t, Status: pgtype.Present}
}

func StringToText(s string) pgtype.Text {
	return pgtype.Text{String: s, Status: pgtype.Present}
}

func Int64ToInt8(i int64) pgtype.Int8 {
	return pgtype.Int8{Int: i, Status: pgtype.Present}
}

func Int8(i int) pgtype.Int8 {
	return Int64ToInt8(int64(i))
}

func ByteSliceToJSONB(d []byte) pgtype.JSONB {
	return pgtype.JSONB{Bytes: d, Status: pgtype.Present}
}

func Bool(b bool) pgtype.Bool {
	return pgtype.Bool{Bool: b}
}
