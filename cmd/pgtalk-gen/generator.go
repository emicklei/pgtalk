package main

import (
	"bytes"
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
	if *oVerbose {
		log.Printf("generating from %s.%s\n", table.Schema, table.Name)
	}
	tt := TableType{
		Created:    time.Now(),
		Schema:     *oSchema,
		TableName:  table.Name,
		TableAlias: alias(table.Name),
		GoPackage:  table.Name,
		GoType:     asSingular(strcase.ToCamel(table.Name)),
	}
	for _, each := range table.Columns {
		goType, method := goFieldTypeAndAccess(each.DataType, each.NotNull)
		f := ColumnField{
			Name:                 each.Name,
			GoName:               fieldName(each.Name),
			GoType:               goType,
			NonPointerGoType:     goType[1:],
			DataType:             each.DataType,
			FactoryMethod:        method,
			IsPrimarySrc:         isPrimarySource(each.IsPrimaryKey),
			IsNotNullSrc:         isNotNullSource(each.NotNull),
			IsPrimary:            each.IsPrimaryKey,
			IsNotNull:            each.NotNull,
			TableAttributeNumber: each.FieldOrdinal,
			ValueFieldName:       nullableValueFieldName(each.DataType),
			IsGenericFieldAccess: strings.HasPrefix(method, "NewField"),
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
	abb := abbreviate(s)
	index := 1
	if known, ok := knownAliases[abb]; ok {
		index = known + 1
	}
	a := fmt.Sprintf("%s%d", abb, index)
	knownAliases[abb] = index
	return a
}

// happy_world => hw
func abbreviate(s string) string {
	b := new(bytes.Buffer)
	start := true
	for _, each := range s {
		if start {
			b.WriteRune(each)
			start = false
		}
		if each == '_' || each == '.' {
			start = true
		}
	}
	return b.String()
}

func goFieldTypeAndAccess(datatype string, notNull bool) (string, string) {
	switch datatype {
	case "date":
		if notNull {
			return "time.Time", "NewTimeAccess"
		}
		return "pgtype.Date", "NewTimeAccess"
	case "timestamp with time zone":
		if notNull {
			return "time.Time", "NewTimeAccess"
		}
		return "pgtype.Timestamptz", "NewFieldAccess[pgtype.Timestamptz]"
	case "timestamp", "timestamp without time zone":
		if notNull {
			return "time.Time", "NewTimeAccess"
		}
		return "pgtype.Timestamp", "NewFieldAccess[pgtype.Timestamp]"
	case "text":
		if notNull {
			return "string", "NewTextAccess"
		}
		return "pgtype.Text", "NewTextAccess"
	case "bigint", "integer":
		if notNull {
			return "int64", "NewInt64Access"
		}
		return "pgtype.Int8", "NewInt64Access"
	case "jsonb":
		if notNull {
			return "string", "NewJSONBAccess"
		}
		return "pgtype.JSONB", "NewJSONBAccess"
	case "point":
		return "pgtype.Point", "NewFieldAccess[pgtype.Point]"
	case "boolean":
		if notNull {
			return "bool", "NewBooleanAccess"
		}
		return "pgtype.Bool", "NewBooleanAccess"
	case "daterange":
		return "pgtype.Daterange", "NewFieldAccess[pgtype.Daterange]"
	case "interval":
		return "pgtype.Interval", "NewFieldAccess[pgtype.Interval]"
	case "bytea":
		return "pgtype.Bytea", "NewFieldAccess[pgtype.Bytea]"
	case "text[]":
		return "pgtype.TextArray", "NewFieldAccess[pgtype.TextArray]"
	case "uuid":
		return "pgtype.UUID", "NewFieldAccess[pgtype.UUID]"
	}
	if strings.HasPrefix(datatype, "character") {
		if notNull {
			return "string", "NewTextAccess"
		}
		return "pgtype.Text", "NewTextAccess"
	}
	if strings.HasPrefix(datatype, "numeric") {
		if notNull {
			return "float64", "NewFieldAccess[float64]"
		}
		return "pgtype.Float8", "NewFieldAccess[pgtype.Float8]"
	}
	if *oVerbose {
		log.Println("[WARN] unknown datatype, using fallback for:", datatype)
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

func isPrimarySource(isPrimary bool) string {
	// import package is aliased to "p"
	if isPrimary {
		return "p.IsPrimary"
	}
	return "p.NotPrimary"
}

func isNotNullSource(isNotNull bool) string {
	// import package is aliased to "p"
	if isNotNull {
		return "p.NotNull"
	}
	return "p.Nullable"
}

func nullableValueFieldName(dataType string) string {
	switch dataType {
	case "text":
		return "String"
	case "jsonb":
		return "Bytes"
	case "timestamp with time zone", "date", "timestamp without time zone":
		return "Time"
	case "bigint":
		return "Int"
	}
	return "UNKOWN:" + dataType
}
