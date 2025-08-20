package pgtalk

import "testing"


func TestInt32InEmpty(t *testing.T) {
	ta := NewInt32Access(ci, nil)
	in := ta.In()
	sql := SQL(in)
	if got, want := sql, "false"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestInt32FieldValueToScan(t *testing.T) {
	var i int32 = 42
	ta := NewInt32Access(ci, func(dest any) any {
		return &i
	})
	if got, want := ta.FieldValueToScan(nil), &i; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestInt32AppendScannable(t *testing.T) {
	ta := NewInt32Access(ci, nil)
	list := []any{}
	list = ta.AppendScannable(list)
	if got, want := len(list), 1; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestInt32TableAlias(t *testing.T) {
	ta := NewInt32Access(ci, nil)
	taa := ta.TableAlias("other")
	sql := SQL(taa.Equals(1))
	if got, want := sql, "(other.label = 1)"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestInt32Get(t *testing.T) {
	ta := NewInt32Access(ci, nil)
	row := map[string]any{"label": int32(42)}
	if got, want := ta.Get(row), int32(42); got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestInt32Set(t *testing.T) {
	ta := NewInt32Access(ci, nil)
	ta = ta.Set(42)
	if got, want := ta.ValueToInsert(), int32(42); got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestInt32BetweenAnd(t *testing.T) {
	ta := NewInt32Access(ci, nil)
	in := ta.BetweenAnd(1, 100)
	sql := SQL(in)
	if got, want := sql, "(t1.label BETWEEN 1 AND 100)"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestInt32Compare(t *testing.T) {
	ta := NewInt32Access(ci, nil)
	in := ta.Compare(">", 1)
	sql := SQL(in)
	if got, want := sql, "(t1.label > 1)"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestInt32In(t *testing.T) {
	ta := NewInt32Access(ci, nil)
	in := ta.In(1, 2)
	sql := SQL(in)
	if got, want := sql, "(t1.label IN (1,2))"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestInt32Euals(t *testing.T) {
	ta := NewInt32Access(ci, nil)
	in := ta.Equals(int32(1))
	sql := SQL(in)
	if got, want := sql, "(t1.label = 1)"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}
func TestIntEuals(t *testing.T) {
	ta := NewInt32Access(ci, nil)
	in := ta.Equals(int(1))
	sql := SQL(in)
	if got, want := sql, "(t1.label = 1)"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}
