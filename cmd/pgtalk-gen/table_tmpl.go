package main

var tableTemplateSrc = `package {{.GoPackage}}
// generated by pgtalk-gen on {{.Created}} 
// DO NOT EDIT

import (
	"github.com/emicklei/pgtalk"
	"time"
	"fmt"
	"bytes"
	"github.com/jackc/pgx/v4"
	"context"
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
	tableAccess = pgtalk.TableAccessor{TableInfo: tableInfo, AllColumns: []pgtalk.ColumnAccessor{
{{- range .Fields}}{{.GoName}},{{- end}}
}}
)

// ColumnUpdatesFrom returns the list of changes to a {{.GoType}} for which updates need to be processed.
// Cannot be used to set null values for columns.
func ColumnUpdatesFrom(e {{.GoType}}) (list []pgtalk.SQLWriter) {
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

// Select returns a new {{.GoType}}sQuerySet for fetching column data.
func Select(cas ...pgtalk.ColumnAccessor) {{.GoType}}sQuerySet {
	return {{.GoType}}sQuerySet{pgtalk.MakeQuerySet(tableAccess, cas, func() interface{} {
		return new({{.GoType}})
	})}
}

// {{.GoType}}sQuerySet can query for *{{.GoType}} values.
type {{.GoType}}sQuerySet struct {
	pgtalk.QuerySet
}

func (s {{.GoType}}sQuerySet) Unwrap() pgtalk.QuerySet { return s.QuerySet }

// Where returns a new QuerySet with WHERE clause.
func (s {{.GoType}}sQuerySet) Where(condition pgtalk.SQLWriter) {{.GoType}}sQuerySet {
	return {{.GoType}}sQuerySet{QuerySet: s.QuerySet.Where(condition)}
}

// Limit returns a new QuerySet with the maximum number of results set.
func (s {{.GoType}}sQuerySet) Limit(limit int) {{.GoType}}sQuerySet {
	return {{.GoType}}sQuerySet{QuerySet: s.QuerySet.Limit(limit)}
}

// GroupBy returns a new QuerySet with the GROUP BY clause.
func (s {{.GoType}}sQuerySet) GroupBy(cas ...pgtalk.ColumnAccessor) {{.GoType}}sQuerySet {
	return {{.GoType}}sQuerySet{QuerySet: s.QuerySet.GroupBy(cas...)}
}

// OrderBy returns a new QuerySet with the ORDER BY clause.
func (s {{.GoType}}sQuerySet) OrderBy(cas ...pgtalk.ColumnAccessor) {{.GoType}}sQuerySet {
	return {{.GoType}}sQuerySet{QuerySet: s.QuerySet.OrderBy(cas...)}
}

// Exec runs the query and returns the list of *{{.GoType}}.
func (s {{.GoType}}sQuerySet) Exec(ctx context.Context,conn *pgx.Conn) (list []*{{.GoType}}, err error) {
	err = s.QuerySet.ExecWithAppender(ctx, conn, func(each interface{}) {
		list = append(list, each.(*{{.GoType}}))
	})
	return
}

// Insert creates a MutationSet for inserting data with zero or more columns.
func Insert(cas ...pgtalk.ColumnAccessor) pgtalk.MutationSet {
	return pgtalk.MakeMutationSet(tableAccess, cas, pgtalk.MutationInsert)
}

// Delete creates a MutationSet for deleting data.
func Delete() pgtalk.MutationSet {
	return pgtalk.MakeMutationSet(tableAccess, pgtalk.EmptyColumnAccessor, pgtalk.MutationDelete)
}

// Update creates a MutationSet to update zero or more columns.
func Update(cas ...pgtalk.ColumnAccessor) pgtalk.MutationSet {
	return pgtalk.MakeMutationSet(tableAccess, cas, pgtalk.MutationUpdate)
}
`
