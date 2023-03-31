package main

var tableTemplateSrc = `package {{.GoPackage}}
// generated by github.com/emicklei/pgtalk/cmd/pgtalk-gen from version: {{ .BuildVersion}} 
// DO NOT EDIT

import (
	p "github.com/emicklei/pgtalk"
	c "github.com/emicklei/pgtalk/convert"
	"time"
	"strings"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shopspring/decimal" 
)

// {{.GoType}} is generated from the {{.Schema}}.{{.TableName}} table.
type {{.GoType}} struct {
{{- range .Fields}}
	{{.GoName}}	{{.GoType}} // {{.Name}} : {{.DataType}}
{{- end}}
	// for storing custom field expression result values
	expressionResults map[string]any
}

var (
{{- range .Fields}}	
	// {{.GoName}} represents the column "{{.Name}}" of with type "{{.DataType}}", nullable:{{not .IsNotNull}}, primary:{{.IsPrimary}}
	{{.GoName}} = p.{{.FactoryMethod}}(p.MakeColumnInfo(tableInfo, "{{.Name}}", {{.IsPrimarySrc}}, {{.IsNotNullSrc}}, {{.TableAttributeNumber}}),
		func(dest any) any { return &dest.(*{{$.GoType}}).{{.GoName}} })
{{- end}}
	// package private
	_ = c.UUID // for the occasional unused import from convert
	_ = time.Now 
	_ = pgtype.Empty // for the occasional unused import from pgtype
	_ = decimal.Decimal{}
	tableInfo = p.TableInfo{Schema: "{{.Schema}}", Name: "{{.TableName}}", Alias: "{{.TableAlias}}" }
)

func init() {
	// after var initialization (to prevent cycle) we need to update the tableInfo to set all columns
	tableInfo.Columns = []p.ColumnAccessor{ {{- range .Fields}}{{.GoName}},{{- end}} }
}

{{- range .Fields}}
{{- if .IsNotNull }}

// Set{{.GoName}} sets the value to the field value and returns the receiver.
func (e *{{$.GoType}}) Set{{.GoName}}(v {{.GoType}}) *{{$.GoType}} { e.{{.GoName}} = v ; return e }
{{- else }}

{{ if eq .ConvertFuncName "" }}
// Set{{.GoName}} sets the value to the field value and returns the receiver.
func (e *{{$.GoType}}) Set{{.GoName}}(v {{.GoType}}) *{{$.GoType}} { e.{{.GoName}} = v ; return e }
{{- else}}
// Set{{.GoName}} sets the value to the field value and returns the receiver.
func (e *{{$.GoType}}) Set{{.GoName}}(v {{.NonConvertedGoType}}) *{{$.GoType}} { e.{{.GoName}} = c.{{.ConvertFuncName}}(v) ; return e }
{{- end }}

{{- end }}
{{- end}}

// Setters returns the list of changes to a {{.GoType}} for which updates/inserts need to be processed.
// Can be used in Insert,Update,Select. Cannot be used to set null values for columns.
func (e *{{.GoType}}) Setters() (list []p.ColumnAccessor) {
{{- range .Fields}}
	{{- if .IsNotNull }}
	list = append(list, {{.GoName}}.Set(e.{{.GoName}}))
	{{- else }}
	if e.{{.GoName}}{{.IsValidSrc}} {
		{{- if .IsGenericFieldAccess }}
		list = append(list, {{.GoName}}.Set(e.{{.GoName}}))
		{{- else }}
		list = append(list, {{.GoName}}.Set(e.{{.GoName}}.{{.ValueFieldName}}))
		{{- end }}
	}
	{{- end }}	
{{- end}}	
	return
}

// String returns the debug string for *{{.GoType}} with all non-nil field values.
func (e *{{.GoType}}) String() string {
	return p.StringWithFields(e, p.HideNilValues)
}

// Columns returns the ColumnAccessor list for the given column names.
// If the names is empty then return all columns.
func Columns(names ...string) (list []p.ColumnAccessor) {
	if len(names) == 0 {
		return tableInfo.Columns
	}
	for _, each := range names {
		for _, other := range tableInfo.Columns {
			n := other.Column().Name()
			if strings.HasPrefix(n,"'") { // mixed case names are quoted
				n = strings.Trim(n,"'")
			} 
			if n == each {
				list = append(list, other)
			}
		}
	}
	return
}

// AddExpressionResult puts a value into the custom expression results
func (e *{{.GoType}}) AddExpressionResult(key string, value any) {
	if e.expressionResults == nil {
		// lazy initialize
		e.expressionResults = map[string]any{}
	}
	e.expressionResults[key]=value
}

// GetExpressionResult gets a value from the custom expression results. Returns nil if absent.
func (e *{{.GoType}}) GetExpressionResult(key string) any {
	v, ok := e.expressionResults[key]
	if !ok {
		return nil
	}
	pv := v.(*any)
	return *pv
}

// Select returns a new QuerySet[{{.GoType}}] for fetching column data.
func Select(cas ...p.ColumnAccessor) p.QuerySet[{{.GoType}}] {
	return p.MakeQuerySet[{{.GoType}}](tableInfo, cas)
}

// Insert creates a MutationSet for inserting data with zero or more columns.
func Insert(cas ...p.ColumnAccessor) p.MutationSet[{{.GoType}}] {
	return p.MakeMutationSet[{{.GoType}}](tableInfo, cas, p.MutationInsert)
}

// Delete creates a MutationSet for deleting data.
func Delete() p.MutationSet[{{.GoType}}] {
	return p.MakeMutationSet[{{.GoType}}](tableInfo, p.EmptyColumnAccessor, p.MutationDelete)
}

// Update creates a MutationSet to update zero or more columns.
func Update(cas ...p.ColumnAccessor) p.MutationSet[{{.GoType}}] {
	return p.MakeMutationSet[{{.GoType}}](tableInfo, cas, p.MutationUpdate)
}
`
