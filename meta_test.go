package pgtalk

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/emicklei/pgtalk/products"
	"github.com/jackc/pgx/v4"
)

func TestReadProductTable(t *testing.T) {
	ta, err := LoadTables(context.Background(), testConnect, "public")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%#v", ta)
}

func TestAccessProductTable(t *testing.T) {
	products, err := products.Select(products.ID, products.Code).Exec(testConnect)
	log.Printf("%v,%v,%v", *products[0].ID, *products[0].Code, err)
}

func TestSelectProductsWhere(t *testing.T) {
	q := products.
		Select(products.ID, products.Code).
		Where(products.Code.Equals("F42").
			And(products.ID.Equals(1))).
		Limit(1)
	t.Log(q.SQL())
	products, err := q.Exec(testConnect)
	log.Printf("%v,%v,%v", *products[0].ID, *products[0].Code, err)
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
