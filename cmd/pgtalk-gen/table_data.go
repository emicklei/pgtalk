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
	Name                 string
	DataType             string
	GoStructType         string
	GoType               string
	NonPointerGoType     string
	GoName               string
	FactoryMethod        string
	IsPrimary            bool
	IsPrimarySrc         string
	IsNotNull            bool
	IsNotNullSrc         string
	TableAttributeNumber int
	ValueFieldName       string //  can be empty
	IsGenericFieldAccess bool
}
