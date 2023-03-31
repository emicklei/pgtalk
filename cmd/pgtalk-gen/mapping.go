package main

type mapping struct {
	goFieldType string // non-nullable type
	newFuncCall string // to create accessor for non-nullable type
	// null
	nullableGoFieldType    string // full name of the nullable type
	nullableValueFieldName string // to access the go field value of a nullable type
	convertFuncName        string // to convert from a go field value to a nullable type
	newAccessFuncCall      string // to create the accessor
}

// https://www.postgresql.org/docs/9.1/datatype-numeric.html
var pgMappings = map[string]mapping{
	"timestamp with time zone": {
		goFieldType: "time.Time",
		newFuncCall: "NewTimeAccess",

		nullableValueFieldName: "Time",
		convertFuncName:        "TimeToTimestamptz",
		nullableGoFieldType:    "pgtype.Timestamptz",
		newAccessFuncCall:      "NewFieldAccess[pgtype.Timestamptz]",
	},
	"timestamp without time zone": {
		goFieldType: "time.Time",
		newFuncCall: "NewTimeAccess",

		nullableValueFieldName: "Time",
		convertFuncName:        "TimeToTimestamp",
		nullableGoFieldType:    "pgtype.Timestamp",
		newAccessFuncCall:      "NewFieldAccess[pgtype.Timestamp]",
	},
	"date": {
		goFieldType: "time.Time",
		newFuncCall: "NewTimeAccess",

		nullableValueFieldName: "Time",
		convertFuncName:        "TimeToDate",
		nullableGoFieldType:    "pgtype.Date",
		newAccessFuncCall:      "NewFieldAccess[pgtype.Date]",
	},
	"citext": {
		goFieldType: "string",
		newFuncCall: "NewTextAccess",

		nullableValueFieldName: "String",
		convertFuncName:        "StringToText",
		nullableGoFieldType:    "pgtype.Text",
		newAccessFuncCall:      "NewFieldAccess[pgtype.Text]",
	},
	"text": {
		goFieldType: "string",
		newFuncCall: "NewTextAccess",

		nullableValueFieldName: "String",
		convertFuncName:        "StringToText",
		nullableGoFieldType:    "pgtype.Text",
		newAccessFuncCall:      "NewFieldAccess[pgtype.Text]",
	},
	"character varying": {
		goFieldType: "string",
		newFuncCall: "NewTextAccess",

		nullableValueFieldName: "String",
		convertFuncName:        "StringToText",
		nullableGoFieldType:    "pgtype.Text",
		newAccessFuncCall:      "NewFieldAccess[pgtype.Text]",
	},
	"bigint": {
		goFieldType: "int64",
		newFuncCall: "NewInt64Access",

		nullableValueFieldName: "Int",
		convertFuncName:        "Int64ToInt8",
		nullableGoFieldType:    "pgtype.Int8",
		newAccessFuncCall:      "NewFieldAccess[pgtype.Int8]",
	},
	"integer": {
		goFieldType: "int32",
		newFuncCall: "NewInt32Access",

		nullableValueFieldName: "Int",
		convertFuncName:        "Int32ToInt4",
		nullableGoFieldType:    "pgtype.Int4",
		newAccessFuncCall:      "NewFieldAccess[pgtype.Int4]",
	},
	"jsonb": {
		nullableValueFieldName: "Map",
		nullableGoFieldType:    "p.JSON",
		goFieldType:            "p.JSON",
		newAccessFuncCall:      "NewJSONAccess",
	},
	"json": {
		nullableValueFieldName: "Map",
		nullableGoFieldType:    "p.JSON",
		goFieldType:            "p.JSON",
		newAccessFuncCall:      "NewJSONAccess",
	},
	"uuid": {
		goFieldType:         "pgtype.UUID",
		nullableGoFieldType: "pgtype.UUID",
		newAccessFuncCall:   "NewFieldAccess[pgtype.UUID]",
	},
	// https://github.com/jackc/pgx/wiki/Numeric-and-decimal-support
	"numeric": {
		goFieldType:         "decimal.NullDecimal",
		nullableGoFieldType: "decimal.NullDecimal",
		newAccessFuncCall:   "NewFieldAccess[decimal.NullDecimal]",
	},
	"decimal": {
		goFieldType:         "decimal.NullDecimal",
		nullableGoFieldType: "decimal.NullDecimal",
		newAccessFuncCall:   "NewFieldAccess[decimal.NullDecimal]",
	},
	"point": {
		nullableGoFieldType: "pgtype.Point",
		newAccessFuncCall:   "NewFieldAccess[pgtype.Point]",
	},
	"boolean": {
		goFieldType: "bool",
		newFuncCall: "NewBooleanAccess",

		convertFuncName:        "Bool",
		nullableValueFieldName: "Bool",
		nullableGoFieldType:    "pgtype.Bool",
		newAccessFuncCall:      "NewFieldAccess[pgtype.Bool]",
	},
	"daterange": {
		nullableGoFieldType: "pgtype.Range[pgtype.Date]",
		goFieldType:         "pgtype.Range[pgtype.Date]",
		newAccessFuncCall:   "NewFieldAccess[pgtype.Range[pgtype.Date]]",
	},
	"bytea": {
		nullableGoFieldType: "pgtype.Bytea",
		newAccessFuncCall:   "NewFieldAccess[pgtype.Bytea]",
	},
	"text[]": {
		nullableGoFieldType: "pgtype.TextArray",
		newAccessFuncCall:   "NewFieldAccess[pgtype.TextArray]",
	},
	"interval": {
		nullableGoFieldType: "pgtype.Interval",
		newAccessFuncCall:   "NewFieldAccess[pgtype.Interval]",
	},
}
