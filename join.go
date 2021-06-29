package pgtalk

import (
	"context"

	"github.com/jackc/pgx/v4"
)

type JoinType int

const (
	InnerJoinType JoinType = iota
	LeftOuterJoinType
	RightOuterJoinType
	FullOuterJoinType
)

type Join struct {
	LeftSet  QuerySet
	RightSet QuerySet
	OnLeft   ColumnAccessor
	OnRight  ColumnAccessor
	Type     JoinType
}

func (i Join) SQL() string {
	//TODO InnerJoinType only
	return "SELECT " + i.LeftSet.SelectSection() + "," + i.RightSet.SelectSection() +
		" FROM " + i.LeftSet.FromSection() +
		" INNER JOIN " + i.RightSet.FromSection() +
		" ON (" + i.OnLeft.SQL() +
		" = " + i.OnRight.SQL() +
		") WHERE " + i.LeftSet.WhereSection()
}

func (i Join) On(onLeft, onRight ColumnAccessor) Join {
	return Join{
		LeftSet:  i.LeftSet,
		RightSet: i.RightSet,
		OnLeft:   onLeft,
		OnRight:  onRight,
	}
}

func (i Join) Exec(conn Connection) (it JoinResultIterator, err error) {
	rows, err := conn.Query(context.Background(), i.SQL())
	if err != nil {
		return
	}
	return JoinResultIterator{leftSet: i.LeftSet, rightSet: i.RightSet, rows: rows}, nil
}

type JoinResultIterator struct {
	leftSet  QuerySet
	rightSet QuerySet
	rows     pgx.Rows
}

func (i *JoinResultIterator) HasNext() bool {
	if i.rows.Next() {
		return true
	} else {
		i.rows.Close()
	}
	return false
}

func (i *JoinResultIterator) Err() error {
	return i.rows.Err()
}

func (i *JoinResultIterator) Next(left interface{}, right interface{}) error {
	sw := []interface{}{}
	// left
	for _, each := range i.leftSet.selectors {
		rw := ScanToWrite{
			RW:     each,
			Entity: left,
		}
		sw = append(sw, rw)
	}
	// right
	for _, each := range i.rightSet.selectors {
		rw := ScanToWrite{
			RW:     each,
			Entity: right,
		}
		sw = append(sw, rw)
	}
	return i.rows.Scan(sw...)
}
