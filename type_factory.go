package pgtalk

import (
	"time"

	"github.com/jackc/pgtype"
)

func MakeTimestamptz(t *time.Time) pgtype.Timestamptz {
	if t == nil {
		return pgtype.Timestamptz{Status: pgtype.Null}
	}
	return pgtype.Timestamptz{Time: *t, Status: pgtype.Present}
}

func MakeTimestamp(t time.Time) pgtype.Timestamp {
	return pgtype.Timestamp{Time: t, Status: pgtype.Present}
}

func MakeDate(t time.Time) pgtype.Date {
	return pgtype.Date{Time: t, Status: pgtype.Present}
}
