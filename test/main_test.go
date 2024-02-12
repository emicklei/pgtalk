package test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

var testConnect *pgx.Conn

func TestMain(m *testing.M) {
	connectionString := os.Getenv("PGTALK_CONN") // "postgres://postgres:pgtalk@localhost:7432/postgres"
	if len(connectionString) == 0 {
		println("no database env set")
		os.Exit(m.Run())
		return
	}
	fmt.Println("db open ...", connectionString)
	conn, err := pgx.Connect(context.Background(), connectionString)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		println("no database available so tests in this package are skipped")
		os.Exit(0)
	}
	testConnect = conn
	if err := ensureTables(conn); err != nil {
		fmt.Println("DB WARN:", err)
	}
	uuid.EnableRandPool()
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
		id uuid,
		tDate date,
		tTimestamp timestamp without time zone,
		TJSONB jsonb,
		TJSON json,
		tText text,
		tNumeric numeric,
		tDecimal decimal
	);
	drop table IF EXISTS products;
	create table products(
		id serial primary key,
		created_at  timestamp with time zone,
		updated_at  timestamp with time zone,
		deleted_at  timestamp with time zone,
		code       text,
		price      bigint,
		category_id integer			
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

func diff(left, right string) string {
	//assume one line
	b := new(bytes.Buffer)
	io.WriteString(b, "\n")
	io.WriteString(b, left)
	io.WriteString(b, "\n")
	leftRunes := []rune(left)
	rightRunes := []rune(right)
	size := len(leftRunes)
	if l := len(rightRunes); l < size {
		size = l
	}
	for c := 0; c < size; c++ {
		l := leftRunes[c]
		r := rightRunes[c]
		if l == r {
			b.WriteRune(l)
		} else {
			fmt.Fprintf(b, "^(%s)...", string(r))
			break
		}
	}
	return b.String()
}

func oneliner(s string) string {
	return strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(s, "\t", " "), "\n", " "), "  ", " ")
}
