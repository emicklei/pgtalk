package test

import (
	"context"
	"testing"
	"time"

	"github.com/emicklei/pgtalk"
	"github.com/emicklei/pgtalk/convert"
	"github.com/emicklei/pgtalk/test/things"
	"github.com/google/uuid"
	"github.com/jackc/pgtype"
)

func TestTableInfoColumnsOfThingsNotEmpty(t *testing.T) {
	if got, want := len(things.AllColumns()), 4; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestManageThing(t *testing.T) {
	// CREATE
	ctx := context.Background()
	id := uuid.New()
	t.Log(id.String())
	create := things.Insert(
		things.ID.Set(convert.UUID(id)),
		things.Tdate.Set(convert.TimeToDate(time.Now())),
		things.Ttimestamp.Set(convert.TimeToTimestamp(time.Now())),
		things.Tjson.Set([]byte(`{"key":"value"}`)))
	t.Log(pgtalk.PrettySQL(create))
	tx, err := testConnect.Begin(ctx)
	if err != nil {
		t.Fatal(err)
	}
	_ = create.Exec(ctx, testConnect)
	err = tx.Commit(ctx)
	if err != nil {
		t.Fatal(err)
	}
	// READ
	{
		read := things.Select(things.ID, things.Tdate, things.Ttimestamp, things.Tjson).Where(things.ID.Equals(convert.UUID(id)))
		t.Log(pgtalk.PrettySQL(read))
		list, err := read.Exec(ctx, testConnect)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(list)
	}
	// UPDATE
	{
		update := things.Update(things.Tdate.Set(pgtype.Date{Status: pgtype.Null})).Where(things.ID.Equals(convert.UUID(id)))
		t.Log(pgtalk.PrettySQL(update))
		tx, err = testConnect.Begin(ctx)
		if err != nil {
			t.Fatal(err)
		}
		update.Exec(ctx, testConnect)
		err = tx.Commit(ctx)
		if err != nil {
			t.Fatal(err)
		}
	}
	// READ
	{
		read := things.Select(things.ID, things.Tdate, things.Ttimestamp, things.Tjson).Where(things.ID.Equals(convert.UUID(id)))
		t.Log(pgtalk.PrettySQL(read))
		list, err := read.Exec(ctx, testConnect)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(list)
		if got, want := list[0].Tdate.Status, pgtype.Null; got != want {
			t.Errorf("got [%v]:%T want [%v]:%T", got, got, want, want)
		}
	}
}

func TestJSONB(t *testing.T) {
	ctx := context.Background()
	if testConnect == nil {
		return
	}

	// delete 2
	tx, _ := testConnect.Begin(ctx)
	del := things.Delete()
	it := del.Exec(ctx, testConnect)
	t.Log(pgtalk.SQL(del))
	if it.Err() != nil {
		t.Fatal(it.Err())
	}
	tx.Commit(ctx)

	// insert 2
	m := things.Insert(
		things.ID.Set(convert.UUID(uuid.New())),
		//things.Tdate.Set(time.Now()),
		//things.Ttimestamp.Set(convert.TimeToTimestamp(time.Now())),
		things.Tjson.Set([]byte(`{"key":"value"}`)))

	// insert 3
	// {
	// 	obj := new(things.Thing)
	// 	obj.SetID(convert.StringToUUID(uuid.NewString())).SetTdate(time.Now())
	// 	things.Insert(obj.Setters()...)
	// }

	tx, err := testConnect.Begin(ctx)
	if err != nil {
		t.Fatal(err)
	}
	r := m.Exec(ctx, testConnect)
	if err := r.Err(); err != nil {
		for i, each := range m.ValuesToInsert() {
			t.Logf("%d:%v", i, each)
		}
		t.Log(pgtalk.SQL(m))
		t.Fatal(err)
	}
	// TODO: mutation knows whether rows are returned ;if not it can close the connection itself
	if r.HasNext() {
		t.Fatal("no data expected")
	}

	err = tx.Commit(ctx)
	if err != nil {
		t.Fatal(err)
	}
}

func TestJSONB_3(t *testing.T) {
	t.Skip()
	ctx := context.Background()

	// delete 3
	tx, _ := testConnect.Begin(ctx)
	// reverse columns
	s := things.AllColumns()
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	del := things.Delete().Where(things.ID.Equals(3)).Returning(s...)
	it := del.Exec(ctx, testConnect)
	t.Log(pgtalk.SQL(del))
	t.Log(it.Err())
	for it.HasNext() {
		var th *things.Thing
		th, err := it.Next()
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("deleted:%#v\n", th)
	}
	tx.Commit(ctx)

	// insert 3
	m := things.Insert(
		things.ID.Set(convert.UUID(uuid.New())),
		//things.Tdate.Set(convert.TimeToDate(time.Now())),
		//things.Ttimestamp.Set(convert.TimeToTimestamp(time.Now())),
		things.Tjson.Set([]byte(`{"key":"value"}`)))

	tx, err := testConnect.Begin(ctx)
	if err != nil {
		t.Fatal(err)
	}
	r := m.Exec(ctx, testConnect)
	if err := r.Err(); err != nil {
		t.Log(pgtalk.SQL(m))
		t.Fatal(err)
	}
	// TODO: mutation knows whether rows are returned ;if not it can close the connection itself
	if r.HasNext() {
		t.Fatal("no data expected")
	}

	err = tx.Commit(ctx)
	if err != nil {
		t.Fatal(err)
	}
}

func TestExtraJSONBField(t *testing.T) {
	a := things.Tjson.Extract("title")
	t.Log(pgtalk.PrettySQL(a))
}
