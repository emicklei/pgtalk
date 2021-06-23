package products

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/emicklei/pgtalk/xs"
	"github.com/jackc/pgx/v4"
)

type Product struct {
	ID    *int64
	Code  *string
	Price *int64
}

// func Table(conn *pgx.Conn) *ProductsQuery {
// 	return &ProductsQuery{conn: conn}
// }

var ID = xs.NewInt8Access(
	"id",
	func(dest interface{}, i *int64) {
		e := dest.(*Product)
		e.ID = i
	})

var Code = xs.NewTextAccess(
	"code",
	func(dest interface{}, i *string) {
		e := dest.(*Product)
		e.Code = i
	})

type ProductsQuery struct {
	selectors []xs.ReadWrite
}

func Select(as ...xs.ReadWrite) ProductsQuery {
	return ProductsQuery{selectors: as}
}

func (d ProductsQuery) Exec(conn *pgx.Conn) (list []*Product, err error) {
	// TODO use GORM here
	buf := new(bytes.Buffer)
	for i, each := range d.selectors {
		if i > 0 {
			io.WriteString(buf, ",")
		}
		io.WriteString(buf, each.Name())
	}
	rows, err := conn.Query(context.Background(), fmt.Sprintf("select %s from products", buf))
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		entity := new(Product)
		sw := []interface{}{}
		for _, each := range d.selectors {
			rw := xs.ScanToWrite{
				RW:     each,
				Entity: entity,
			}
			sw = append(sw, rw)
		}
		if err := rows.Scan(sw...); err != nil {
			return list, err
		}
		list = append(list, entity)
	}
	return
}

//////
