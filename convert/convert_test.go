package convert

import (
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
)

func TestUUIDToString(t *testing.T) {
	i := pgtype.UUID{
		Bytes: [16]byte{1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4},
		Valid: true,
	}
	if got, want := UUIDToString(i), "01020304-0102-0304-0102-030401020304"; got != want {
		t.Errorf("got [%v]:%T want [%v]:%T", got, got, want, want)
	}
	i.Valid = false
	if got, want := UUIDToString(i), ""; got != want {
		t.Errorf("got [%v]:%T want [%v]:%T", got, got, want, want)
	}
}
