package pgtalk

import (
	"time"

	"github.com/jackc/pgtype"
)

func NewTimestamptz(t *time.Time) pgtype.Timestamptz {
	if t == nil {
		return pgtype.Timestamptz{Status: pgtype.Null}
	}
	return pgtype.Timestamptz{Time: *t, Status: pgtype.Present}
}
