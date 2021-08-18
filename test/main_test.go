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
	pgtalk.EnableAssert()
	code := m.Run()
	fmt.Println("... db close")
	conn.Close(context.Background())
	os.Exit(code)
}
