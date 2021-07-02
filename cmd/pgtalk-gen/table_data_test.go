package main

import (
	"os"
	"testing"
	"text/template"
	"time"
)

func TestTemplate(t *testing.T) {
	tt := TableType{
		Created:    time.Now(),
		TableName:  "prices",
		TableAlias: "p1",
		GoPackage:  "prices",
		GoType:     "Price",
		Fields: []ColumnField{
			{
				Name:          "id",
				GoName:        "ID",
				GoType:        "*int64",
				FactoryMethod: "NewInt8Access",
			},
			{
				Name:          "currency",
				GoName:        "Currency",
				GoType:        "*string",
				FactoryMethod: "NewTextAccess",
			},
		},
	}
	tmpl, err := template.New("tt").Parse(tableTemplateSrc)
	check(t, err)
	check(t, tmpl.Execute(os.Stdout, tt))
}

func check(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}
