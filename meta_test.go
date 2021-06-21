package pgtalk

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/jackc/pgx/v4"
)

func TestReadProductTable(t *testing.T) {

	ta, err := LoadTables(context.Background(), testConnect, "public")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%#v", ta)
}

var testConnect *pgx.Conn

func TestMain(m *testing.M) {
	connectionString := os.Getenv("PGTALK_CONN")
	if len(connectionString) == 0 {
		os.Exit(1)
	}
	fmt.Println("db open ...")
	conn, err := pgx.Connect(context.Background(), connectionString)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	testConnect = conn
	code := m.Run()
	fmt.Println("... db close")
	conn.Close(context.Background())
	os.Exit(code)
}
