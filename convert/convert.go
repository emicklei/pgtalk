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

func StringToUUID(s string) (pgtype.UUID, bool) {
	data, err := parseUUID(s)
	if err != nil {
		return pgtype.UUID{}, false
	}
	return pgtype.UUID{
		Bytes:  data,
		Status: pgtype.Present,
	}, true
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

func UUIDToString(t pgtype.UUID) (string, bool) {
	if t.Status != pgtype.Present {
		return "", false
	}
	data, err := t.MarshalJSON()
	if err != nil {
		return "", false
	}
	return string(data), true
}

func TimeToTimestamptz(t *time.Time) pgtype.Timestamptz {
	if t == nil {
		return pgtype.Timestamptz{Status: pgtype.Null}
	}
	return pgtype.Timestamptz{Time: *t, Status: pgtype.Present}
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
