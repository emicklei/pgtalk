package things

import "github.com/emicklei/pgtalk"

func Collect(cas ...pgtalk.ColumnAccessor) pgtalk.UntypedQuerySet {
	return pgtalk.NewUntypedQuerySet(tableInfo, cas)
}
