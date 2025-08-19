package main

import (
	"context"
	"encoding/json"
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
	oExcludePatterns = flag.String("exclude", "", "comma separated list of regexp for tables to exclude")
	oMapping         = flag.String("mapping", "", "mapping file for undefined pg data types")
	oCache           = flag.String("cache", "", "use cache file for table metadata if present")
)

func main() {
	flag.Parse()
	if err := applyConfiguredMappings(*oMapping); err != nil {
		fmt.Fprintf(os.Stderr, "Unable to process custom mappings: %v\n", err)
		os.Exit(1)
	}
	var tables []PgTable
	if *oCache != "" {
		if result, err := loadPgTablesFromCache(*oCache); err != nil {
			fmt.Fprintf(os.Stderr, "unable to load cache file: %v\ncontinue using database connection", err)
		} else {
			tables = result
		}
	}
	if len(tables) == 0 {
		tables = fetchPgTables()
		if *oCache != "" {
			if err := savePgTablesToCache(*oCache, tables); err != nil {
				fmt.Fprintf(os.Stderr, "unable to save cache file: %v\n", err)
			}
		}
	}
	filter := NewTableFilter(*oIncludePatterns, *oExcludePatterns)
	for _, each := range tables {
		if filter.Includes(each.Name) {
			if *oDryrun {
				log.Println("[-dry] would generate", each.Name)
			} else {
				generateFromTable(each, *oViews)
			}
		}
	}
}

func fetchPgTables() []PgTable {
	connectionString := os.Getenv("PGTALK_CONN")
	if len(connectionString) == 0 {
		fmt.Fprintf(os.Stderr, "Missing value of environment variable PGTALK_CONN\n")
		os.Exit(1)
	}
	log.Println("fetching tables from database")
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
	return all
}

func loadPgTablesFromCache(cacheFile string) ([]PgTable, error) {
	log.Printf("loading tables from cache: %s", cacheFile)
	data, err := os.ReadFile(cacheFile)
	if err != nil {
		return nil, fmt.Errorf("unable to read cache file: %w", err)
	}
	var tables []PgTable
	if err := json.Unmarshal(data, &tables); err != nil {
		return nil, fmt.Errorf("unable to unmarshal cache file: %w", err)
	}
	return tables, nil
}

func savePgTablesToCache(cacheFile string, tables []PgTable) error {
	log.Printf("saving tables to cache: %s", cacheFile)
	data, err := json.MarshalIndent(tables, "", "  ")
	if err != nil {
		return fmt.Errorf("unable to marshal tables to JSON: %w", err)
	}
	if err := os.WriteFile(cacheFile, data, 0644); err != nil {
		return fmt.Errorf("unable to write cache file: %w", err)
	}
	return nil
}
