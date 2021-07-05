# pgtalk

More type safe SQL query building using Go code generated from PostgreSQL table definitions.

## examples

### Insert

	m := products.Insert(
		products.ID.Set(10),
		products.Code.Set("test"),
		products.CategoryID.Set(1))

	err := m.Exec(aConnection)		

### Update

	m := products.Update(
		products.ID.Set(10),
		products.Code.Set("test"),
		products.CategoryID.Set(1)).
		Where(products.ID.Equals(10))

	err := m.Exec(aConnection)		

### Delete

	m := products.Delete().Where(products.ID.Equals(10))

	err := m.Exec(aConnection)

### Select

	q := products.Select(products.Code).Where(products.Code.Equals("F42"))

	products, err := q.Exec(aConnection) // products is a []*product.Product

### Left Outer Join

    q :=products.Select(products.Code).Where(products.Code.Equals("F42")).
        LeftJoin(categories.Select(categories.Title)).
        On(products.ID, categories.ID)

	it, _ := q.Exec(aConnection)
	for it.HasNext() {
		p := new(products.Product)
		c := new(categories.Category)
		_ = it.Next(p, c)
		t.Logf("%s,%s", *p.Code, *c.Title)
	}

## supported Column Types

- text
- bigint

(c) 2021, http://ernestmicklei.com. MIT License.