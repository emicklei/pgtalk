package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
)

var (
	oDryrun          = flag.Bool("dry", false, "do not generate, report only")
	oTarget          = flag.String("o", ".", "target directory")
	oSchema          = flag.String("s", "public", "source database schema")
	oViews           = flag.Bool("views", false, "generated from views, default is false = use tables")
	oVerbose         = flag.Bool("v", false, "use verbose logging")
	oIncludePatterns = flag.String("include", ".*", "comma separated list of regexp for tables to include")
	oExludePatterns  = flag.String("exclude", "", "comma separated list of regexp for tables to exclude")
	oMapping         = flag.String("mapping", "", "mapping file for undefined pg data types")
)

func main() {
	flag.Parse()
	connectionString := os.Getenv("PGTALK_CONN")
	if len(connectionString) == 0 {
		fmt.Fprintf(os.Stderr, "Missing value of environment variable PGTALK_CONN\n")
		os.Exit(1)
	}

	if err := applyConfiguredMappings(*oMapping); err != nil {
		fmt.Fprintf(os.Stderr, "Unable to process custom mappings: %v\n", err)
		os.Exit(1)
	}

	conn, err := pgx.Connect(context.Background(), connectionString)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	var all []PgTable
	if *oViews {
		all, err = LoadViews(context.Background(), conn, *oSchema)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		all, err = LoadTables(context.Background(), conn, *oSchema)
		if err != nil {
			log.Fatal(err)
		}
	}
	filter := NewTableFilter(*oIncludePatterns, *oExludePatterns)
	for _, each := range all {
		if filter.Includes(each.Name) {
			if *oDryrun {
				log.Println("[DRYRUN] would generate", each.Name)
			} else {
				generateFromTable(each, *oViews)
			}
		}
	}
}
