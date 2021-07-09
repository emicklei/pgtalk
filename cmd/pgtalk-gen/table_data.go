package main

import "time"

type TableType struct {
	Created    time.Time
	Schema     string
	TableName  string
	TableAlias string
	GoPackage  string
	GoType     string
	Fields     []ColumnField
}

type ColumnField struct {
	Name          string
	DataType      string
	GoStructType  string
	GoType        string
	GoName        string
	FactoryMethod string
}
