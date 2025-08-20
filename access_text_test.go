package pgtalk

import (
	"reflect"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
)

var (
	ti = TableInfo{
		Name:   "things",
		Schema: "public",
		Alias:  "t1",
	}
	ci = ColumnInfo{
		tableInfo:  ti,
		columnName: "label",
	}
)

func TestTextIn(t *testing.T) {
	ta := NewTextAccess(ci, nil)
	in := ta.In("a", "b")
	sql := SQL(in)
	if got, want := sql, "(t1.label IN ('a','b'))"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestTextColumn(t *testing.T) {
	ta := NewTextAccess(ci, nil)
	if got, want := ta.Column(), ci; !reflect.DeepEqual(got, want) {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestTextGet(t *testing.T) {
	ta := NewTextAccess(ci, nil)
	t.Run("key exists", func(t *testing.T) {
		row := map[string]any{"label": "some-text"}
		v := ta.Get(row)
		if got, want := v, "some-text"; got != want {
			t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
		}
	})
	t.Run("key does not exist", func(t *testing.T) {
		row := map[string]any{"other": "some-text"}
		v := ta.Get(row)
		if _, ok := v.(pgtype.Text); !ok {
			t.Errorf("got %T want pgtype.Text", v)
		}
	})
}

func TestTextAppendScannable(t *testing.T) {
	ta := NewTextAccess(ci, nil)
	list := []any{}
	list = ta.AppendScannable(list)
	if len(list) != 1 {
		t.Errorf("got %d want 1", len(list))
	}
	if _, ok := list[0].(*string); !ok {
		t.Errorf("got %T want *string", list[0])
	}
}

func TestTextFieldValueToScan(t *testing.T) {
	type myEntity struct {
		Label string
	}
	entity := myEntity{Label: "my-label"}
	varLabel := NewTextAccess(ci, func(entity any) any {
		return &entity.(*myEntity).Label
	})
	scanner := varLabel.FieldValueToScan(&entity)
	*(scanner.(*string)) = "new-label"
	if got, want := entity.Label, "new-label"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestTextTableAlias(t *testing.T) {
	ta := NewTextAccess(ci, nil)
	aliased := ta.TableAlias("other")
	sql := SQL(aliased.Equals("test"))
	if got, want := sql, "(other.label = 'test')"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestTextCompare(t *testing.T) {
	ta1 := NewTextAccess(ci, nil)
	ci2 := ColumnInfo{
		tableInfo:  ti,
		columnName: "extra",
	}
	ta2 := NewTextAccess(ci2, nil)
	t.Run("compare with other", func(t *testing.T) {
		op := ta1.Compare(">", ta2)
		sql := SQL(op)
		if got, want := sql, "(t1.label > t1.extra)"; got != want {
			t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
		}
	})
	t.Run("compare with parameter", func(t *testing.T) {
		op := ta1.Compare("<", NewParameter(42))
		sql := SQL(op)
		if got, want := sql, "(t1.label < ?)"; got != want {
			t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
		}
	})
	t.Run("invalid operator", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("The code did not panic")
			}
		}()
		ta1.Compare("!", "other")
	})
}

func TestTextILike(t *testing.T) {
	ta := NewTextAccess(ci, nil)
	in := ta.ILike("*me")
	sql := SQL(in)
	if got, want := sql, "(t1.label ILIKE '*me')"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestTextSet(t *testing.T) {
	ta := NewTextAccess(ci, nil)
	ta = ta.Set("hello")
	if got, want := ta.ValueToInsert(), "hello"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestTextInEmpty(t *testing.T) {
	ta := NewTextAccess(ci, nil)
	in := ta.In()
	sql := SQL(in)
	if got, want := sql, "false"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestTextLike(t *testing.T) {
	ta := NewTextAccess(ci, nil)
	in := ta.Like("*me")
	sql := SQL(in)
	if got, want := sql, "(t1.label LIKE '*me')"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestTextNotNull(t *testing.T) {
	ta := NewTextAccess(ci, nil)
	op := IsNotNull(ta)
	sql := SQL(op)
	if got, want := sql, "(t1.label IS NOT NULL)"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestTextEquals(t *testing.T) {
	ta := NewTextAccess(ci, nil)
	op := ta.Equals("me")
	sql := SQL(op)
	if got, want := sql, "(t1.label = 'me')"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}
