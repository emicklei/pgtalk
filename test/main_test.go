package test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/emicklei/pgtalk"
	"github.com/jackc/pgx/v4"
)

var testConnect *pgx.Conn

func TestMain(m *testing.M) {
	// connectionString := "postgres://postgres:pgtalk@localhost:5432/pgtalk"
	connectionString := os.Getenv("PGTALK_CONN")
	if len(connectionString) == 0 {
		println("no database env set")
		os.Exit(m.Run())
		return
	}
	fmt.Println("db open ...", connectionString)
	conn, err := pgx.Connect(context.Background(), connectionString)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	testConnect = conn
	if err := ensureTables(conn); err != nil {
		fmt.Println("DB WARN:", err)
	}
	pgtalk.EnableAssert()
	code := m.Run()
	fmt.Println("... db close")
	conn.Close(context.Background())
	os.Exit(code)
}

func ensureTables(conn *pgx.Conn) error {
	ctx := context.Background()
	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}
	_, err = conn.Exec(ctx, `
	drop table IF EXISTS things;
	create table things(
		id serial primary key,
		tDate date,
		tTimestamp timestamp without time zone,
		TJSON jsonb
	);
	drop table IF EXISTS products;
	create table products(
		id serial primary key,
		created_at  timestamp with time zone,
		updated_at  timestamp with time zone,
		deleted_at  timestamp with time zone,
		code       text,
		price      bigint,
		category_id bigint			
	);
	drop table IF EXISTS categories;
	create table categories(
		id serial primary key,
		title text
	);`)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}
