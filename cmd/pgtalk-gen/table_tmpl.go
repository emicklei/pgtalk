package main

var tableTemplateSrc = `package {{.GoPackage}}
// generated by pgtalk-gen on {{.Created}} 
// DO NOT EDIT

import (
	"github.com/emicklei/pgtalk"
	"time"
	"fmt"
	"bytes"
)

var (
	_ = time.Now()
	tableInfo = pgtalk.TableInfo{Schema: "{{.Schema}}", Name: "{{.TableName}}", Alias: "{{.TableAlias}}"}
)

type {{.GoType}} struct {
{{- range .Fields}}
	{{.GoName}}	{{.GoType}} // {{.DataType}}
{{- end}}
}

var (
{{- range .Fields}}	
	{{.GoName}} = pgtalk.{{.FactoryMethod}}(pgtalk.MakeColumnInfo(tableInfo, "{{.Name}}", {{.IsPrimary}}, {{.IsNotNull}}, {{.TableAttributeNumber}}),
		func(dest interface{}, v {{.GoType}}) { dest.(*{{$.GoType}}).{{.GoName}} = v })
{{- end}}
	tableAccess = pgtalk.TableAccessor{TableInfo: tableInfo, 
		Factory: func() interface{}{return new({{.GoType}})}, AllColumns: []pgtalk.ColumnAccessor{
{{- range .Fields}}{{.GoName}},{{- end}}
}}
)

// ColumnUpdatesFrom returns the list of changes to a {{.GoType}} for which updates need to be processed.
// Cannot be used to set null values for columns.
func ColumnUpdatesFrom(e *{{.GoType}}) (list []pgtalk.SQLExpression) {
{{- range .Fields}}
	if e.{{.GoName}} != nil {
		list = append(list, {{.GoName}}.Set(*e.{{.GoName}}))
	}
{{- end}}	
	return
}

// Next returns the next *{{.GoType}} from the iterator data.
// Use err to check for failure.
func Next(it *pgtalk.ResultIterator) (e *{{.GoType}}, err error) {
	var each = new({{.GoType}})
	// first check for query error in case caller forgot
	if err = it.Err(); err != nil {
		return nil, err
	}
	err = it.Next(each)
	return each, err
}

// String returns the debug string for *{{.GoType}} with all non-nil field values.
func (e *{{.GoType}}) String() string {
	b := new(bytes.Buffer)
	fmt.Fprint(b, "{{.TableName}}.{{.GoType}}{")
{{- range .Fields}}
	if e.{{.GoName}} != nil {
		fmt.Fprintf(b, "{{.GoName}}:%v ", *e.{{.GoName}})
	}
{{- end}}
	fmt.Fprint(b, "}")
	return b.String()
}

// AllColumns returns the list of all column accessors for usage in e.g. Select.
func AllColumns() []pgtalk.ColumnAccessor {
	return tableAccess.AllColumns
}

// Select returns a new QuerySet[{{.GoType}}] for fetching column data.
func Select(cas ...pgtalk.ColumnAccessor) pgtalk.QuerySet[{{.GoType}}] {
	return pgtalk.MakeQuerySet[{{.GoType}}](tableAccess, cas)
}

// Insert creates a MutationSet for inserting data with zero or more columns.
func Insert(cas ...pgtalk.ColumnAccessor) pgtalk.MutationSet[{{.GoType}}] {
	return pgtalk.MakeMutationSet[{{.GoType}}](tableAccess, cas, pgtalk.MutationInsert)
}

// Delete creates a MutationSet for deleting data.
func Delete() pgtalk.MutationSet[{{.GoType}}] {
	return pgtalk.MakeMutationSet[{{.GoType}}](tableAccess, pgtalk.EmptyColumnAccessor, pgtalk.MutationDelete)
}

// Update creates a MutationSet to update zero or more columns.
func Update(cas ...pgtalk.ColumnAccessor) pgtalk.MutationSet[{{.GoType}}] {
	return pgtalk.MakeMutationSet[{{.GoType}}](tableAccess, cas, pgtalk.MutationUpdate)
}
`
