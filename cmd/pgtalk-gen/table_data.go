package main

type TableType struct {
	TableName  string
	TableAlias string
	GoPackage  string
	GoType     string
	Fields     []ColumnField
}

type ColumnField struct {
	Name          string
	GoStructType  string
	GoType        string
	GoName        string
	FactoryMethod string
}
