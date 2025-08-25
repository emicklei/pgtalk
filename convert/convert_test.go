package convert

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func TestStringToUUID(t *testing.T) {
	s := "123e4567-e89b-12d3-a456-426614174000"
	u := StringToUUID(s)
	if !u.Valid {
		t.Fatal("valid uuid expected")
	}
	if got, want := u.Bytes, [16]byte{0x12, 0x3e, 0x45, 0x67, 0xe8, 0x9b, 0x12, 0xd3, 0xa4, 0x56, 0x42, 0x66, 0x14, 0x17, 0x40, 0x00}; got != want {
		t.Errorf("got %v want %v", got, want)
	}
	if want := StringToUUID("123"); want.Valid {
		t.Error("invalid uuid expected")
	}
	if want := StringToUUID("123e4567-e89b-12d3-a456-42661417400"); want.Valid {
		t.Error("invalid uuid expected")
	}
	s = "123e4567e89b12d3a456426614174000"
	u = StringToUUID(s)
	if !u.Valid {
		t.Fatal("valid uuid expected")
	}
}

func TestParseUUID(t *testing.T) {
	if _, err := parseUUID("123"); err == nil {
		t.Error("error expected")
	}
	if _, err := parseUUID("123e4567-e89b-12d3-a456-42661417400k"); err == nil {
		t.Error("error expected")
	}
}

func TestUUID(t *testing.T) {
	u := uuid.New()
	pu := UUID(u)
	if !pu.Valid {
		t.Error("valid pgtype.UUID expected")
	}
}

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

func TestTime(t *testing.T) {
	now := time.Now()
	if got, want := TimeToTimestamptz(now).Time, now; !got.Equal(want) {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := TimeToTimestamp(now).Time, now.UTC(); !got.Equal(want) {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := TimeToDate(now).Time, now; !got.Equal(want) {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := TimeToTime(now).Microseconds, now.UnixMicro(); got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestInts(t *testing.T) {
	if got, want := Int16ToInt2(1).Int16, int16(1); got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := Int64ToInt8(1).Int64, int64(1); got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := Int4ToInt8(pgtype.Int4{Int32: 1, Valid: true}).Int64, int64(1); got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := Int32ToInt4(1).Int32, int32(1); got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := Int8(1).Int64, int64(1); got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := Int4(1).Int32, int32(1); got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestOthers(t *testing.T) {
	if got, want := StringToText("a").String, "a"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := Bool(true).Bool, true; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	d := pgtype.Date{Time: time.Now(), Valid: true}
	if got := DateToTimePtr(d); got == nil {
		t.Error("got nil want time.Time")
	}
	d.Valid = false
	if got := DateToTimePtr(d); got != nil {
		t.Error("got time.Time want nil")
	}
	if got, want := Float64ToFloat8(1.0).Float64, 1.0; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	sa := StringsToTextArray([]string{"a", "b"})
	if got, want := len(sa), 2; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestTextArrayToStrings(t *testing.T) {
	a := pgtype.FlatArray[pgtype.Text]{}
	a = append(a, pgtype.Text{String: "a", Valid: true})
	a = append(a, pgtype.Text{String: "b", Valid: true})
	s := TextArrayToStrings(a)
	if got, want := len(s), 2; got != want {
		t.Errorf("got [%[1]v:%[1]T] want [%[2]v:%[2]T]", got, want)
	}
	if got, want := s[0], "a"; got != want {
		t.Errorf("got [%[1]v:%[1]T] want [%[2]v:%[2]T]", got, want)
	}
	if got, want := s[1], "b"; got != want {
		t.Errorf("got [%[1]v:%[1]T] want [%[2]v:%[2]T]", got, want)
	}
}
