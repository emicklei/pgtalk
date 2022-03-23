# pgtalk

[![Build Status](https://api.travis-ci.com/emicklei/pgtalk.svg?branch=main)](https://travis-ci.com/github/emicklei/pgtalk)
[![GoDoc](https://pkg.go.dev/badge/github.com/emicklei/pgtalk)](https://pkg.go.dev/github.com/emicklei/pgtalk)

More type safe SQL query building and execution using Go code generated (pgtalk-gen) from PostgreSQL table definitions.

## generics

This package requires Go SDK version 1.18+ because it uses type parameterization.

## examples

These examples are from the test package in which a few database tables files (categories,products,things) are generated.

### Insert

	m := products.Insert(
		products.ID.Set(10),
		products.Code.Set("testit"),
		products.CategoryID.Set(1))

	it := m.Exec(aConnection)
	if err := it.Err(); err != nil {
		....
	}

or by example:

	p := new(products.Product)
	p.SetID(10).SetCode("testit").SetCategoryID(1)		
	m := products.Insert(p.Setters()...)

### Update

	m := products.Update(
			products.Code.Set("testme"),
			products.CategoryID.Set(1)).
		Where(products.ID.Equals(10)).
		Returning(products.Code)

	it := m.Exec(aConnection)	
	for it.HasNext() {
		p, err := products.Next(p) // p is a *product.Product
		t.Logf("%s,%s", *p.Code)
	}		

or by example

	p := new(products.Product)
	p.SetID(10).SetCode("testme").SetCategoryID(1)		
	m := products.Update(p.Setters()...).
			Where(products.ID.Equals(p.ID)).
			Returning(products.Code)

### Delete

	m := products.Delete().Where(products.ID.Equals(10))

	_ = m.Exec(aConnection)

### Select

	q := products.Select(products.Code).Where(products.Code.Equals("F42"))

	products, err := q.Exec(aConnection) // products is a []*product.Product

## Abritrary SQL expressions

	q := products.Select(products.ID, pgtalk.FieldSQL("UPPER(p1.Code)", "upper"))
	// SELECT p1.id,UPPER(p1.Code) AS upper FROM public.products p1
	list, _ := q.Exec(context.Background(),aConnection)
	for _, each := range list {
		upper := each.GetExpressionResult("upper")
		...
	}

## Joins

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

### Multi Join

	q := offer.Select(offer.Columns()...).
		Where(offer.ID.In("01F674G8MDRAPBWA6SB1HWE2VC").And(offer.IsActive.Equals(true))).
		LeftOuterJoin(offer_buyer_permission_rel.Select(offer_buyer_permission_rel.BuyerPermission)).
		On(offer_buyer_permission_rel.OfferId.Equals(offer.ID)).
		LeftOuterJoin(offer_market_rel.Select(offer_market_rel.MarketId)).
		On(offer_market_rel.OfferId.Equals(offer.ID)).
		LeftOuterJoin(buyer_market_rel.Select().Where(buyer_market_rel.BuyerId.Equals("X1010_0100002"))).
		On(offer_market_rel.MarketId.Equals(offer_market_rel.MarketId))

	t.Log(pgtalk.PrettySQL(q))

	it, err := q.Exec(context.Background(), testConnect)
	if err != nil {
		t.Fatal(err)
	}
	for it.HasNext() {
		offer := new(offer.Offer)
		permission := new(offer_buyer_permission_rel.OfferBuyerPermissionRel)
		market := new(offer_market_rel.OfferMarketRel)
		// The order and types of the entities must match the order of the non-empty Select functions used in the query
		err := it.Next(offer, permission, market)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(offer)
		t.Log(permission)
		t.Log(market)
	}

## supported Column Types

- text, character, uuid
- bigint, integer, number
- date, timestamp
- jsonb, bytes
- daterange
- point
- bytea, text[], interval

Send me a PR for a missing type available from https://www.postgresql.org/docs/9.5/datatype.html .

## how to run the generator

	PGTALK_CONN=postgresql://usr:pwd@host:5432/database pgtalk-gen -o yourpackage
	go fmt ./...

If you want to include and/or exclude table names, use additional flags such as:

	pgtalk-gen -o yourpackage -include "address.*,employee.*" -exclude "org.*"

(c) 2022, http://ernestmicklei.com. MIT License.
