package main

var tableTemplateSrc = `
package {{.GoPackage}}

import (
	"github.com/emicklei/pgtalk"
)

var tableInfo = pgtalk.TableInfo{Name: "{{.TableName}}", Alias: "{{.TableAlias}}"}

type {{.GoType}} struct {
{{range .Fields}}
	{{.GoName}}	{{.GoType}}
{{end}}
}

var (
{{range .Fields}}	
	{{.GoName}} = pgtalk.{{.FactoryMethod}}(tableInfo, "{{.Name}}", func(dest interface{}, v {{.GoType}}) { dest.(*{{$.GoType}}).{{.GoName}} = v })
{{end}}
)

func AllColumns() (all []pgtalk.ColumnAccessor) {
	return append(all{{range .Fields}},{{.GoName}}{{end}})
}

func Select(cas ...pgtalk.ColumnAccessor) {{.GoType}}sQuerySet {
	return {{.GoType}}sQuerySet{pgtalk.MakeQuerySet(tableInfo, cas, func() interface{} {
		return new({{.GoType}})
	})}
}

type {{.GoType}}sQuerySet struct {
	pgtalk.QuerySet
}

func (s {{.GoType}}sQuerySet) Unwrap() pgtalk.QuerySet { return s.QuerySet }

// Where is
func (s {{.GoType}}sQuerySet) Where(condition pgtalk.SQLWriter) {{.GoType}}sQuerySet {
	return {{.GoType}}sQuerySet{QuerySet: s.QuerySet.Where(condition)}
}

// Limit is
func (s {{.GoType}}sQuerySet) Limit(limit int) {{.GoType}}sQuerySet {
	return {{.GoType}}sQuerySet{QuerySet: s.QuerySet.Limit(limit)}
}

// GroupBy is
func (s {{.GoType}}sQuerySet) GroupBy(cas ...pgtalk.ColumnAccessor) {{.GoType}}sQuerySet {
	return {{.GoType}}sQuerySet{QuerySet: s.QuerySet.GroupBy(cas...)}
}

// GroupBy is
func (s {{.GoType}}sQuerySet) OrderBy(cas ...pgtalk.ColumnAccessor) {{.GoType}}sQuerySet {
	return {{.GoType}}sQuerySet{QuerySet: s.QuerySet.OrderBy(cas...)}
}

// Exec is
func (s {{.GoType}}sQuerySet) Exec(conn pgtalk.Connection) (list []*{{.GoType}}, err error) {
	err = s.QuerySet.ExecWithAppender(conn, func(each interface{}) {
		list = append(list, each.(*{{.GoType}}))
	})
	return
}

func Insert(cas ...pgtalk.ColumnAccessor) pgtalk.MutationSet {
	return pgtalk.MakeMutationSet(tableInfo, cas, pgtalk.MutationInsert)
}

func Delete() pgtalk.MutationSet {
	return pgtalk.MakeMutationSet(tableInfo, pgtalk.EmptyColumnAccessor, pgtalk.MutationDelete)
}

func Update(cas ...pgtalk.ColumnAccessor) pgtalk.MutationSet {
	return pgtalk.MakeMutationSet(tableInfo, cas, pgtalk.MutationUpdate)
}
`