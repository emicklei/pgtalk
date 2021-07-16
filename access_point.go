package pgtalk

import (
	"fmt"
	"io"
)

type Point struct{}

type PointAccess struct {
	ColumnInfo
	fieldWriter func(dest interface{}, b *Point)
	insertValue string
}

func NewPointAccess(info ColumnInfo, writer func(dest interface{}, b *Point)) PointAccess {
	return PointAccess{ColumnInfo: info, fieldWriter: writer}
}

func (a PointAccess) WriteInto(entity interface{}, fieldValue interface{}) {
	if fieldValue == nil {
		return
	}
	// TODO
	// a.fieldWriter(entity, &Point)
}

func (a PointAccess) InsertValue() interface{} {
	return a.insertValue
}

func (a PointAccess) ValueAsSQLOn(w io.Writer) {
	fmt.Fprintf(w, "%v", a.insertValue) // TODO
}
