package main

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/iancoleman/strcase"
)

func generateFromTable(table PgTable) {
	log.Printf("generating from %s.%s\n", table.Schema, table.Name)
	tt := TableType{
		Created:    time.Now(),
		Schema:     *oSchema,
		TableName:  table.Name,
		TableAlias: alias(table.Name),
		GoPackage:  table.Name,
		GoType:     asSingular(strings.Title(table.Name)),
	}
	for _, each := range table.Columns {
		goType, method := goFieldTypeAndAccess(each.DataType)
		f := ColumnField{
			Name:                 each.Name,
			GoName:               fieldName(each.Name),
			GoType:               goType,
			NonPointerGoType:     goType[1:],
			DataType:             each.DataType,
			FactoryMethod:        method,
			IsPrimary:            each.IsPrimaryKey,
			IsNotNull:            each.NotNull,
			TableAttributeNumber: each.FieldOrdinal,
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

var knownAliases = map[string]int{}

func alias(s string) string {
	first := strings.ToLower(s[0:1])
	index := 1
	if known, ok := knownAliases[first]; ok {
		index = known + 1
	}
	a := fmt.Sprintf("%s%d", first, index)
	knownAliases[first] = index
	return a
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
		return "*string", "NewJSONBAccess"
	case "point":
		return "*pgtalk.Point", "NewPointAccess"
	case "boolean":
		return "*bool", "NewBooleanAccess"
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
	return strcase.ToCamel(s)
}

func asSingular(s string) string {
	if strings.HasSuffix(s, "ies") {
		return s[0:len(s)-3] + "y"
	}
	if strings.HasSuffix(s, "ses") {
		return s[0 : len(s)-2]
	}
	if strings.HasSuffix(s, "s") {
		return s[0 : len(s)-1]
	}
	return s
}
