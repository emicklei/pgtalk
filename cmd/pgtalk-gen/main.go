package main

import (
	"context"
	"flag"
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jackc/pgx/v4"
)

var (
	oTarget = flag.String("o", ".", "target directory")
	oSchema = flag.String("s", "public", "source database schema")
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

func generateFromTable(table PgTable) {
	log.Println("generating ", table.Name)
	tt := TableType{
		Created:    time.Now(),
		Schema:     *oSchema,
		TableName:  table.Name,
		TableAlias: alias(table.Name),
		GoPackage:  table.Name,
		GoType:     withoutTrailingS(strings.Title(table.Name)),
	}
	for _, each := range table.Columns {
		goType, method := goFieldTypeAndAccess(each.DataType)
		f := ColumnField{
			Name:          each.Name,
			GoName:        fieldName(each.Name),
			GoType:        goType,
			DataType:      each.DataType,
			FactoryMethod: method,
		}
		tt.Fields = append(tt.Fields, f)
	}
	tmpl, err := template.New("tt").Parse(tableTemplateSrc)
	if err != nil {
		log.Fatal(err)
	}
	path := filepath.Join(*oTarget, table.Name)
	os.MkdirAll(path, os.ModeDir|os.ModePerm)
	path = filepath.Join(path, "table.go")
	fileOut, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	defer fileOut.Close()

	err = tmpl.Execute(fileOut, tt)
	if err != nil {
		log.Fatal(err)
	}
}

func alias(s string) string {
	return strings.ToLower(s[0:1]) + "1" // TODO
}

func goFieldTypeAndAccess(datatype string) (string, string) {
	switch datatype {
	case "date", "timestamp", "timestamp without time zone", "timestamp with time zone":
		return "*time.Time", "NewTimeAccess"
	case "text":
		return "*string", "NewTextAccess"
	case "bigint", "integer":
		return "*int64", "NewInt64Access"
	case "jsonb":
		return "*[]byte", "NewBytesAccess"
	case "point":
		return "*Point", "NewPointAccess"
	}
	if strings.HasPrefix(datatype, "character") {
		return "*string", "NewTextAccess"
	}
	if strings.HasPrefix(datatype, "numeric") {
		return "*float64", "NewFloat64Access"
	}
	return datatype, "New" + datatype
}

func fieldName(s string) string {
	if s == "id" {
		return "ID"
	}
	return strings.Title(s)
}

func withoutTrailingS(s string) string {
	if strings.HasSuffix(s, "s") {
		return s[0 : len(s)-1]
	}
	return s
}
