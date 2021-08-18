package pgtalk

import (
	"fmt"
	"io"
)

type Point struct{}

type PointAccess struct {
	ColumnInfo
	fieldWriter   func(dest interface{}, b *Point)
	valueToInsert string
}

func NewPointAccess(info ColumnInfo, writer func(dest interface{}, b *Point)) PointAccess {
	return PointAccess{ColumnInfo: info, fieldWriter: writer}
}

func (a PointAccess) Column() ColumnInfo { return a.ColumnInfo }

func (a PointAccess) Collect(list []ColumnAccessor) []ColumnAccessor {
	return append(list, a)
}

func (a PointAccess) SetFieldValue(entity interface{}, fieldValue interface{}) error {
	if fieldValue == nil {
		return nil
	}
	// TODO
	// a.fieldWriter(entity, &Point)
	return nil
}

func (a PointAccess) ValueToInsert() interface{} {
	return a.valueToInsert
}

func (a PointAccess) ValueAsSQLOn(w io.Writer) {
	fmt.Fprintf(w, "%v", a.valueToInsert) // TODO
}
