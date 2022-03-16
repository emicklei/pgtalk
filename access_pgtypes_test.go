package pgtalk

import (
	"testing"

	"github.com/jackc/pgtype"
)

func TestUUID(t *testing.T) {
	u := new(pgtype.UUID)
	e := polyFUUID.Equals(*u)
	if got, want := testSQL(e), "(p1.fuuid = '00000000-0000-0000-0000-000000000000'::uuid)"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func testSQL(ex SQLExpression) string {
	return SQL(ex)
}
