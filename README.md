# pgtalk

More type safe SQL query building using Go code generated (pgtalk-gen) from PostgreSQL table definitions.

## examples

See also [booking demo](https://github.com/emicklei/pgtalk-demo).

### Insert

	m := products.Insert(
		products.ID.Set(10),
		products.Code.Set("test"),
		products.CategoryID.Set(1))

	it := m.Exec(aConnection)		

### Update

	m := products.Update(
			products.ID.Set(10),
			products.Code.Set("test"),
			products.CategoryID.Set(1)).
		Where(products.ID.Equals(10)).
		Returning(products.Code)

	it := m.Exec(aConnection)	
	for it.HasNext() {
		p := new(products.Product)
		_ = it.Next(p)
		t.Logf("%s,%s", *p.Code)
	}		

### Delete

	m := products.Delete().Where(products.ID.Equals(10))

	_ = m.Exec(aConnection)

### Select

	q := products.Select(products.Code).Where(products.Code.Equals("F42"))

	products, err := q.Exec(aConnection) // products is a []*product.Product

### Left Outer Join

    q :=products.Select(products.Code).Where(products.Code.Equals("F42")).
        LeftOuterJoin(categories.Select(categories.Title)).
        On(products.ID.Equals(categories.ID))

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
- date
- timestamp
- jsonb
- bytes
- number
- character
- integer

https://www.postgresql.org/docs/9.5/datatype.html

## dev notes

The whole implementation might be better once Go has Type parameters (generics) support.

(c) 2021, http://ernestmicklei.com. MIT License.