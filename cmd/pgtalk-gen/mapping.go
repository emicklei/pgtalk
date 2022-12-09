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
		goFieldType: "int64",
		newFuncCall: "NewInt64Access",

		nullableValueFieldName: "Int",
		convertFuncName:        "Int64ToInt8",
		nullableGoFieldType:    "pgtype.Int8",
		newAccessFuncCall:      "NewFieldAccess[pgtype.Int8]",
	},
	"jsonb": {
		nullableValueFieldName: "Bytes",
		goFieldType:            "[]byte",
		convertFuncName:        "ByteSliceToJSONB",
		nullableGoFieldType:    "pgtype.JSONB",
		newAccessFuncCall:      "NewJSONBAccess",
	},
	"uuid": {
		goFieldType:         "pgtype.UUID",
		nullableGoFieldType: "pgtype.UUID",
		newAccessFuncCall:   "NewFieldAccess[pgtype.UUID]",
	},
	"numeric": {
		goFieldType:         "numeric.Numeric",
		nullableGoFieldType: "numeric.Numeric",
		newAccessFuncCall:   "NewFieldAccess[numeric.Numeric]",
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
		nullableGoFieldType: "pgtype.Daterange",
		goFieldType:         "pgtype.Daterange",
		newAccessFuncCall:   "NewFieldAccess[pgtype.Daterange]",
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
