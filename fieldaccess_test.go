package pgtalk

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func TestPgTypeUUIDEqualsGoogleUUID(t *testing.T) {
	val := uuid.New()
	id := NewFieldAccess[pgtype.UUID](ColumnInfo{columnName: "id", tableInfo: TableInfo{Name: "table"}}, nil)
	eq := id.Equals(val)
	if got, want := SQL(eq), fmt.Sprintf("(.id = '%s'::uuid)", val.String()); got != want {
		t.Errorf("got [%v]:%T want [%v]:%T", got, got, want, want)
	}
}

func TestPgTypeUUIDInEmpty(t *testing.T) {
	id := NewFieldAccess[pgtype.UUID](ColumnInfo{columnName: "id", tableInfo: TableInfo{Name: "table"}}, nil)
	eq := id.In()
	if got, want := SQL(eq), "false"; got != want {
		t.Errorf("got [%v]:%T want [%v]:%T", got, got, want, want)
	}
}
