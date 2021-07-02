package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"os"
	"strings"
	"time"

	"github.com/jackc/pgx/v4"
)

func main() {
	connectionString := os.Getenv("PGTALK_CONN")
	conn, err := pgx.Connect(context.Background(), connectionString)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	all, err := LoadTables(context.Background(), conn, "public")
	if err != nil {
		log.Fatal(err)
	}
	for _, each := range all {
		if each.Name == "things" {
			generateFromTable(each)
		}
	}
}

func generateFromTable(table PgTable) {
	tt := TableType{
		Created:    time.Now(),
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
			FactoryMethod: method,
		}
		tt.Fields = append(tt.Fields, f)
	}
	tmpl, err := template.New("tt").Parse(tableTemplateSrc)
	if err != nil {
		log.Fatal(err)
	}
	err = tmpl.Execute(os.Stdout, tt)
	if err != nil {
		log.Fatal(err)
	}
}

func alias(s string) string {
	return strings.ToLower(s[0:1]) + "1" // TODO
}

func goFieldTypeAndAccess(datatype string) (string, string) {
	switch datatype {
	case "date":
		return "*pgtype.Date", "NewDateAccess"
	default:
		return datatype, "New" + datatype
	}
}

func fieldName(s string) string {
	return strings.Title(s)
}

func withoutTrailingS(s string) string {
	if strings.HasSuffix(s, "s") {
		return s[0 : len(s)-1]
	}
	return s
}
