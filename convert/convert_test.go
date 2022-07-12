package convert

import (
	"math/big"
	"testing"

	"github.com/jackc/pgtype"
)

func TestUUIDToString(t *testing.T) {
	i := pgtype.UUID{
		Bytes:  [16]byte{1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4},
		Status: pgtype.Present,
	}
	if got, want := UUIDToString(i), "01020304-0102-0304-0102-030401020304"; got != want {
		t.Errorf("got [%v]:%T want [%v]:%T", got, got, want, want)
	}
	i.Status = pgtype.Undefined
	if got, want := UUIDToString(i), ""; got != want {
		t.Errorf("got [%v]:%T want [%v]:%T", got, got, want, want)
	}
}

func TestBigFloatToNumeric(t *testing.T) {
	bf := big.NewFloat(3.14159265359)
	num := BigFloatToNumeric(*bf)
	t.Log(num) // TODO check
}
