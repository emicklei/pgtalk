package test

import (
	"testing"

	"github.com/emicklei/pgtalk"
	"github.com/emicklei/pgtalk/test/tables/categories"
	"github.com/emicklei/pgtalk/test/tables/products"
)

func TestRelational(t *testing.T) {
	// cats = products collect:[:e| e category]
	// Collect(products.Select(), Block[products.Product, categories.Category](
	// 	return categories.ID.Equals(products.CategoryID)
	// ))
	t.Log(pgtalk.SQL(categories.ID.Equals(products.CategoryId)))
}
