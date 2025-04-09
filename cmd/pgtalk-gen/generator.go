package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"runtime/debug"
	"slices"
	"strings"
	"text/template"

	"github.com/iancoleman/strcase"
)

var varcharPattern = regexp.MustCompile(`character varying\(\d+\)`)

func generateFromTable(table PgTable, isView bool) {
	if *oVerbose {
		log.Printf("generating from %s.%s\n", table.Schema, table.Name)
	}
	tt := TableType{
		BuildVersion: "(dev)",
		Schema:       *oSchema,
		TableName:    table.Name,
		TableAlias:   alias(table.Name),
		GoPackage:    table.Name,
		GoType:       asSingular(strcase.ToCamel(table.Name)),
		Imports:      []string{"strings"},
	}
	// need version to put in generated files
	bi, ok := debug.ReadBuildInfo()
	if ok && len(bi.Main.Version) > 0 {
		tt.BuildVersion = bi.Main.Version
	}

	// collect imports required for the fields
	imports := []string{}
	for _, each := range table.Columns {
		m, ok := pgMappings[each.DataType]
		if !ok {
			// handle special cases
			if varcharPattern.MatchString(each.DataType) {
				m = pgMappings["character varying"]
			} else {
				log.Println("[warn] missing map entry for", each.DataType, "column '", each.Name, "' is unmapped and only known as column info")
				unmapped := ColumnField{
					Name:         each.Name,
					GoName:       fieldName(each.Name),
					DataType:     each.DataType,
					IsPrimary:    each.IsPrimaryKey,
					IsPrimarySrc: isPrimarySource(each.IsPrimaryKey),
					IsNotNull:    each.NotNull,
					IsNotNullSrc: isNotNullSource(each.NotNull),
				}
				tt.UnmappedFields = append(tt.UnmappedFields, unmapped)
				continue
			}
		}
		goType := m.goFieldType
		if !each.NotNull {
			goType = m.nullableGoFieldType
		}
		factoryMethod := m.newFuncCall
		if !each.NotNull || factoryMethod == "" {
			factoryMethod = m.newAccessFuncCall
		}
		f := ColumnField{
			Name:                 each.Name,
			GoName:               fieldName(each.Name),
			GoType:               goType,
			DataType:             each.DataType,
			FactoryMethod:        factoryMethod,
			IsPrimarySrc:         isPrimarySource(each.IsPrimaryKey),
			IsNotNullSrc:         isNotNullSource(each.NotNull),
			IsPrimary:            each.IsPrimaryKey,
			IsNotNull:            each.NotNull,
			ValueFieldName:       m.nullableValueFieldName,
			IsGenericFieldAccess: isGenericFieldAccess(m.newAccessFuncCall),
			NonConvertedGoType:   m.goFieldType,
			ConvertFuncName:      m.convertFuncName,
			IsValidSrc:           ".Valid",
			IsArray:              m.isArray,
		}
		imports = append(imports, m.imports...)
		tt.Fields = append(tt.Fields, f)
	}
	// sort fields to have stable generated output
	slices.SortFunc(tt.Fields, func(a, b ColumnField) int {
		return strings.Compare(a.Name, b.Name)
	})

	tt.Imports = unique(imports)

	tmpl, err := template.New("tt").Parse(tableTemplateSrc)
	if err != nil {
		log.Fatal(err)
	}
	kind := "tables"
	if isView {
		kind = "views"
	}
	path := filepath.Join(*oTarget, kind, table.Name)
	os.MkdirAll(path, os.ModeDir|os.ModePerm)
	fileName := "table.go"
	if isView {
		fileName = "view.go"
	}
	path = filepath.Join(path, fileName)
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

func isGenericFieldAccess(call string) bool {
	return strings.HasPrefix(call, "p.NewField") || call == "p.NewJSONAccess" // TODO change template i.o this workaround
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

func unique[A comparable](input []A) []A {
	seen := make(map[A]bool)
	var result []A
	for _, v := range input {
		if !seen[v] {
			seen[v] = true
			result = append(result, v)
		}
	}
	return result
}
