package pgtalk

import (
	"fmt"
	"io"
)

type queryCombination struct {
	leftSet  querySet
	operator string // UNION,INTERSECT,EXCEPT
	rightSet querySet
}

func (u queryCombination) SQLOn(w WriteContext) {
	fmt.Fprint(w, "(")
	writeQuerySet(w, u.leftSet)
	fmt.Fprintf(w, " %s ", u.operator)
	writeQuerySet(w, u.rightSet)
	fmt.Fprint(w, ")")
}

func writeQuerySet(w WriteContext, qs querySet) {
	fmt.Fprint(w, "(SELECT\n")
	left := qs.selectAccessors()
	wl := qs.augmentedContext(w)
	writeAccessOn(left, wl)
	fmt.Fprint(w, "\nFROM ")
	qs.fromSectionOn(wl)
	if _, ok := qs.whereCondition().(noCondition); !ok {
		fmt.Fprint(wl, "\nWHERE ")
		qs.whereCondition().SQLOn(wl)
	}
	io.WriteString(w, ")")
}
