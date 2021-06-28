# pgtalk

## example

With generated Go code from a table definition (products), you can write and execute

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