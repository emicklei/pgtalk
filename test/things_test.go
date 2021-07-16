package test

import (
	"context"
	"testing"
	"time"

	"github.com/emicklei/pgtalk"
	"github.com/emicklei/pgtalk/test/things"
)

func TestJSONB(t *testing.T) {
	ctx := context.Background()

	// delete any
	tx, _ := testConnect.Begin(ctx)
	del := things.Delete().Where(things.ID.Equals(2))
	it := del.Exec(ctx, testConnect)
	t.Log(pgtalk.SQL(del))
	t.Log(it.Err())
	it.HasNext() // TODO
	tx.Commit(ctx)

	// insert again
	m := things.Insert(
		things.ID.Set(2),
		things.TDate.Set(time.Now()),
		things.TTimestamp.Set(time.Now()),
		things.TJSON.Set([]byte(`{"key":"value"}`)))

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
