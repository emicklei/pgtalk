package test

import (
	"context"
	"testing"
	"time"

	"github.com/emicklei/pgtalk"
	"github.com/emicklei/pgtalk/convert"
	"github.com/emicklei/pgtalk/test/tables/things"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func TestSelfJoin(t *testing.T) {
	left := things.Select(things.Tdate)
	right := things.Select(things.Tdate).TableAlias("other")
	join := left.Join(right).On(things.ID.Equals(things.ID.TableAlias("other")))
	sql := oneliner(pgtalk.SQL(join))
	if got, want := sql, "SELECT t1.tdate, other.tdate FROM public.things t1 INNER JOIN public.things other ON (t1.id = other.id)"; got != want {
		t.Errorf("got [%v]:%T want [%v]:%T", got, got, want, want)
	}
}

func TestGetColumn(t *testing.T) {
	ids := things.Columns("id")
	t.Log(ids)
}

func TestCustomExpressionExtension(t *testing.T) {
	createAThing(t)
	q := things.Select(things.ID, pgtalk.SQLAs("12 * 24", "id2"))
	if got, want := oneliner(pgtalk.SQL(q)), "SELECT t1.id, (12 * 24) AS id2 FROM public.things t1"; got != want {
		t.Errorf("got [%v]:%T want [%v]:%T", got, got, want, want)
	}
	list, err := q.Exec(context.Background(), testConnect)
	if err != nil {
		t.Fatal(err)
	}
	for _, each := range list {
		cev := each.GetExpressionResult("id2")
		if got, want := cev, int32(288); got != want {
			t.Errorf("got [%v]:%T want [%v]:%T", got, got, want, want)
		}
	}
}

func TestCustomExpression(t *testing.T) {
	createAThing(t)
	q := things.Select(things.ID, things.Ttext.Concat("cc", things.Ttext))
	if got, want := oneliner(pgtalk.SQL(q)), "SELECT t1.id, (t1.ttext || t1.ttext) AS cc FROM public.things t1"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	list, err := q.ExecIntoMaps(context.Background(), testConnect)
	if err != nil {
		t.Fatal(err)
	}
	for _, each := range list {
		id := things.ID.Get(each).(pgtype.UUID)
		t.Logf("%v (%T)", id.Bytes, id)

		cc := each["cc"].(string)
		t.Logf("%v (%T)", cc, cc)
	}
}

func TestSelectMaps(t *testing.T) {
	createAThing(t)
	q := things.Select(things.ID, things.Ttext, things.Tdate)
	list, err := q.ExecIntoMaps(context.Background(), testConnect)
	if err != nil {
		t.Fatal(err)
	}
	for _, each := range list {
		id := things.ID.Get(each).(pgtype.UUID)
		t.Logf("%v (%T)", id.Bytes, id)

		txt := things.Ttext.Get(each).(pgtype.Text)
		t.Logf("%v (%T)", txt.String, txt)

		dt := things.Tdate.Get(each).(pgtype.Date)
		t.Logf("%v (%T)", dt.Time, dt)
	}
}

func TestTableInfoColumnsOfThingsNotEmpty(t *testing.T) {
	if got, want := len(things.Columns()), 10; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestManageThing(t *testing.T) {
	ctx := context.Background()
	id := createAThing(t)
	// READ
	{
		read := things.Select(things.ID, things.Tdate, things.Ttimestamp, things.Tjson, things.Tjsonb).Where(things.ID.Equals(convert.UUID(id)))
		t.Log(pgtalk.SQL(read))
		list, err := read.Exec(ctx, testConnect)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(list)
		t.Log(list[0].Tjson)
		t.Log(list[0].Tjsonb)
	}
	// UPDATE
	{
		update := things.Update(things.Tdate.Set(pgtype.Date{Valid: false})).Where(things.ID.Equals(convert.UUID(id)))
		t.Log(pgtalk.SQL(update))
		tx, err := testConnect.Begin(ctx)
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
		t.Log(pgtalk.SQL(read))
		list, err := read.Exec(ctx, testConnect)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(list)
		if got, want := list[0].Tdate.Valid, false; got != want {
			t.Errorf("got [%v]:%T want [%v]:%T", got, got, want, want)
		}
	}
}

func createAThing(t *testing.T) uuid.UUID {
	// CREATE
	ctx := context.Background()
	id := uuid.New()
	create := things.Insert(
		things.ID.Set(convert.UUID(id)),
		things.Tdate.Set(convert.TimeToDate(time.Now())),
		things.Ttimestamp.Set(convert.TimeToTimestamp(time.Now())),
		things.Tjson.Set(map[string]any{"key1": "value1"}),
		things.Tjsonb.Set(map[string]any{"key2": "value2"}),
		things.Ttext.Set(convert.StringToText("hello")),
		things.Ttextarray.Set(convert.StringsToTextArray([]string{"a", "b", "c"})),
	)
	tx, err := testConnect.Begin(ctx)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(pgtalk.SQL(create))
	_ = create.Exec(ctx, testConnect)
	err = tx.Commit(ctx)
	if err != nil {
		t.Fatal(err)
	}
	return id
}

func TestManageNullThing(t *testing.T) {
	ctx := context.Background()
	id := createNullThing(t)
	// READ
	{
		read := things.Select(things.ID, things.Tdate, things.Ttimestamp, things.Tjson, things.Tjsonb).Where(things.ID.Equals(convert.UUID(id)))
		t.Log(pgtalk.SQL(read))
		list, err := read.Exec(ctx, testConnect)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(list)
		t.Log(list[0].Tdate)
		t.Log(list[0].Ttimestamp)
		t.Log(list[0].Tjson)
		t.Log(list[0].Tjsonb)
		t.Log(list[0].Ttext)
	}
}

func createNullThing(t *testing.T) uuid.UUID {
	// CREATE
	ctx := context.Background()
	id := uuid.New()
	create := things.Insert(
		things.ID.Set(convert.UUID(id)),
		things.Tdate.Set(pgtype.Date{Valid: false}),
		things.Ttimestamp.Set(pgtype.Timestamp{Valid: false}),
		things.Tjson.Set(nil),
		things.Tjsonb.Set(nil),
		things.Ttext.Set(pgtype.Text{Valid: false}),
	)
	tx, err := testConnect.Begin(ctx)
	if err != nil {
		t.Fatal(err)
	}
	_ = create.Exec(ctx, testConnect)
	err = tx.Commit(ctx)
	if err != nil {
		t.Fatal(err)
	}
	return id
}

func TestGetJSONString(t *testing.T) {
	r := testConnect.QueryRow(context.Background(), `SELECT '5'::json;`)
	var v any
	err := r.Scan(&v)
	t.Log(err, v)
}

// func TestPutJSONString(t *testing.T) {
// 	_, err := testConnect.Exec(context.Background(), `INSERT INTO things(id,tjson) values($1,$2);`, uuid.New(), "something")
// 	t.Log(err)
// }

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
		things.Tjson.Set(map[string]any{"key": "value"}))

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
	s := things.Columns()
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
		things.Tdate.Set(convert.TimeToDate(time.Now())),
		things.Ttimestamp.Set(convert.TimeToTimestamp(time.Now())),
		things.Tjson.Set(map[string]any{"key": "value"}))

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
	t.Log(a) // TODO
}

func TestReadTextArray(t *testing.T) {
	ctx := context.Background()
	id := createAThing(t)
	read := things.Select(things.ID, things.Ttextarray).Where(things.ID.Equals(convert.UUID(id)))
	t.Log(pgtalk.SQL(read))
	list, err := read.Exec(ctx, testConnect)
	if err != nil {
		t.Fatal(err)
	}
	if len(list) == 0 {
		t.Fatal("no data")
	}
	if len(list[0].Ttextarray) != 3 {
		t.Fatal("expected 3 elements")
	}
	if !list[0].Ttextarray[0].Valid {
		t.Fatal("expected first element to be valid")
	}
	if list[0].Ttextarray[0].String != "a" {
		t.Fatal("expected first element to be 'a'")
	}
}
