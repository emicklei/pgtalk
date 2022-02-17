package main

var tableTemplateSrc = `package {{.GoPackage}}
// generated by github.com/emicklei/pgtalk/cmd/pgtalk-gen from version: {{ .BuildVersion}} 
// DO NOT EDIT

import (
	p "github.com/emicklei/pgtalk"
	c "github.com/emicklei/pgtalk/convert"
	"time"
	"github.com/jackc/pgtype"
)

// {{.GoType}} is generated from the {{.Schema}}.{{.TableName}} table.
type {{.GoType}} struct {
{{- range .Fields}}
	{{.GoName}}	{{.GoType}} // {{.Name}} : {{.DataType}}
{{- end}}
}

var (
{{- range .Fields}}	
	// {{.GoName}} represents the column "{{.Name}}" of with type "{{.DataType}}", nullable:{{not .IsNotNull}}, primary:{{.IsPrimary}}
	{{.GoName}} = p.{{.FactoryMethod}}(p.MakeColumnInfo(tableInfo, "{{.Name}}", {{.IsPrimarySrc}}, {{.IsNotNullSrc}}, {{.TableAttributeNumber}}),
		{{- if .IsNotNull }}
			func(dest interface{}, v {{.GoType}}) { dest.(*{{$.GoType}}).{{.GoName}} = v }, nil
		{{- else }}
			nil, func(dest interface{}, v {{.GoType}}) { dest.(*{{$.GoType}}).{{.GoName}} = v }
		{{- end }})
{{- end}}
	// package private
	_ = time.Now 
	_ = pgtype.Empty // for the occasional unused import from pgtype
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
	if e.{{.GoName}}.Status == pgtype.Present {
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

// AllColumns returns the list of all column accessors for usage in e.g. Select.
func AllColumns() []p.ColumnAccessor {
	return tableInfo.Columns
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

// Filter returns a new QuerySet[{{.GoType}}] for fetching all column data for which the condition is true.
func Filter(condition p.SQLExpression) p.QuerySet[{{.GoType}}] {
	return p.MakeQuerySet[{{.GoType}}](tableInfo, tableInfo.Columns).Where(condition)
}
`
