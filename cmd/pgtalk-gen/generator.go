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
		m, ok := pgMappings[each.DataType]
		if !ok {
			log.Println("missing map entry for", each.DataType)
			continue
		}
		goType := m.goFieldType
		if !each.NotNull {
			goType = m.nullableGoFieldType
		}
		f := ColumnField{
			Name:                 each.Name,
			GoName:               fieldName(each.Name),
			GoType:               goType,
			DataType:             each.DataType,
			FactoryMethod:        m.newAccessFuncCall,
			IsPrimarySrc:         isPrimarySource(each.IsPrimaryKey),
			IsNotNullSrc:         isNotNullSource(each.NotNull),
			IsPrimary:            each.IsPrimaryKey,
			IsNotNull:            each.NotNull,
			TableAttributeNumber: each.FieldOrdinal,
			ValueFieldName:       m.nullableValueFieldName,
			IsGenericFieldAccess: strings.HasPrefix(m.newAccessFuncCall, "NewField"),
			NonConvertedGoType:   m.goFieldType,
			ConvertFuncName:      m.convertFuncName,
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

type mapping struct {
	goFieldType            string // non-nullable type
	nullableGoFieldType    string // full name of the nullable type
	nullableValueFieldName string // to access the go field value of a nullable type
	convertFuncName        string // to convert from a go field value to a nullable type
	newAccessFuncCall      string // to create the accessor
}

var pgMappings = map[string]mapping{
	"timestamp with time zone": {
		nullableValueFieldName: "Time",
		goFieldType:            "time.Time",
		convertFuncName:        "TimeToTimestampz",
		nullableGoFieldType:    "pgtype.Timestamptz",
		newAccessFuncCall:      "NewFieldAccess[pgtype.Timestamptz]",
	},
	"timestamp without time zone": {
		nullableValueFieldName: "Time",
		goFieldType:            "time.Time",
		convertFuncName:        "TimeToTimestamp",
		nullableGoFieldType:    "pgtype.Timestamp",
		newAccessFuncCall:      "NewFieldAccess[pgtype.Timestamp]",
	},
	"date": {
		nullableValueFieldName: "Time",
		goFieldType:            "time.Time",
		convertFuncName:        "TimeToDate",
		nullableGoFieldType:    "pgtype.Date",
		newAccessFuncCall:      "NewFieldAccess[pgtype.Date]",
	},
	"text": {
		nullableValueFieldName: "String",
		goFieldType:            "string",
		convertFuncName:        "StringToText",
		nullableGoFieldType:    "pgtype.Text",
		newAccessFuncCall:      "NewTextAccess",
	},
	"bigint": {
		nullableValueFieldName: "Int",
		goFieldType:            "int64",
		convertFuncName:        "Int64ToInt8",
		nullableGoFieldType:    "pgtype.Int8",
		newAccessFuncCall:      "NewInt64Access",
	},
	"integer": {
		nullableValueFieldName: "Int",
		goFieldType:            "int64",
		convertFuncName:        "Int64ToInt8",
		nullableGoFieldType:    "pgtype.Int8",
		newAccessFuncCall:      "NewInt64Access",
	},
	"jsonb": {
		nullableValueFieldName: "Bytes",
		goFieldType:            "[]byte",
		convertFuncName:        "ByteSliceToJSONB",
		nullableGoFieldType:    "pgtype.JSONB",
		newAccessFuncCall:      "NewJSONBAccess",
	},
	"uuid": {
		nullableValueFieldName: "-",
		goFieldType:            "string",
		convertFuncName:        "StringToUUID",
		nullableGoFieldType:    "pgtype.UUID",
		newAccessFuncCall:      "NewFieldAccess[pgtype.UUID]",
	},
	"numeric": {
		nullableValueFieldName: "Float",
		goFieldType:            "float64",
		convertFuncName:        "Float64ToFloat8",
		nullableGoFieldType:    "pgtype.Float8",
		newAccessFuncCall:      "NewFieldAccess[pgtype.Float8]",
	},
	"point": {
		nullableValueFieldName: "-",
		goFieldType:            "-",
		convertFuncName:        "-",
		nullableGoFieldType:    "pgtype.Point",
		newAccessFuncCall:      "NewFieldAccess[pgtype.Point]",
	},
	"boolean": {
		nullableValueFieldName: "Bool",
		goFieldType:            "bool",
		convertFuncName:        "Bool",
		nullableGoFieldType:    "pgtype.Bool",
		newAccessFuncCall:      "NewFieldAccess[pgtype.Bool]",
	},
	"daterange": {
		nullableValueFieldName: "-",
		goFieldType:            "-",
		convertFuncName:        "-",
		nullableGoFieldType:    "pgtype.DateRange",
		newAccessFuncCall:      "NewFieldAccess[pgtype.DateRange]",
	},
	// bytea
	// interval
	// text[]
}

/**
	case "interval":
		return "pgtype.Interval", "NewFieldAccess[pgtype.Interval]"
	case "bytea":
		return "pgtype.Bytea", "NewFieldAccess[pgtype.Bytea]"
	case "text[]":
		return "pgtype.TextArray", "NewFieldAccess[pgtype.TextArray]"
**/
