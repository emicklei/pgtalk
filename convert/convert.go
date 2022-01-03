package convert

import (
	"time"

	"github.com/jackc/pgtype"
)

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
