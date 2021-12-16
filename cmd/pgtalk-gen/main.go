package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v4"
)

var (
	oTarget          = flag.String("o", ".", "target directory")
	oSchema          = flag.String("s", "public", "source database schema")
	oVerbose         = flag.Bool("v", true, "use verbose logging")
	oIncludePatterns = flag.String("include", "*", "comma separated list of regexp for tables to include")
	oExludePatterns  = flag.String("exclude", "", "comma separated list of regexp for tables to exclude")
)

func main() {
	flag.Parse()
	connectionString := os.Getenv("PGTALK_CONN")
	conn, err := pgx.Connect(context.Background(), connectionString)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	all, err := LoadTables(context.Background(), conn, *oSchema)
	if err != nil {
		log.Fatal(err)
	}
	for _, each := range all {
		generateFromTable(each)
	}
}
