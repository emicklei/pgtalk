package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type configurableMappingEntry struct {
	Use string `json:"use"` // if this is set then other fields are ignored
	// others
	NullableGoFieldType    string `json:"nullableFieldType"`      // full name of the nullable type
	NullableValueFieldName string `json:"nullableValueFieldName"` // to access the go field value of a nullable type
	ConvertGoFuncName      string `json:"convertFuncName"`        // to convert from a go field value to a nullable type
	NewAccessFuncName      string `json:"newAccessFuncName"`      // to create the accessor

}

func applyConfiguredMappings(location string) error {
	if location == "" {
		return nil
	}
	content, err := os.ReadFile(location)
	if err != nil {
		return err
	}
	entries := map[string]configurableMappingEntry{}
	if err := json.Unmarshal(content, &entries); err != nil {
		return err
	}
	for k, v := range entries {
		// are we using a defined mapping?
		if v.Use != "" {
			// fetch the mapping under "use"
			existing, ok := pgMappings[v.Use]
			if !ok {
				return fmt.Errorf("no such defined mapping: %s", v.Use)
			}
			// make sure existing is not replaced
			_, ok = pgMappings[k]
			if ok {
				return fmt.Errorf("cannot replace mapping: %s", k)
			}
			log.Printf("add datatype mapping %s => %s\n", k, v.Use)
			pgMappings[k] = existing
		} else {
			// custom defined accessor
			// make sure existing is not replaced
			_, ok := pgMappings[k]
			if ok {
				return fmt.Errorf("cannot replace mapping: %s", k)
			}
			newMapping := mapping{
				nullableGoFieldType:    v.NullableGoFieldType,
				nullableValueFieldName: v.NullableValueFieldName,
				convertFuncName:        v.ConvertGoFuncName,
				newAccessFuncCall:      v.NewAccessFuncName,
			}
			if err := newMapping.validate(); err != nil {
				return fmt.Errorf("invalid mapping %s: %v", k, err)
			}
			log.Printf("add new datatype %s\n", k)
			pgMappings[k] = newMapping
		}
	}
	return nil
}

type mapping struct {
	goFieldType string // non-nullable type
	newFuncCall string // to create accessor for non-nullable type
	// null
	nullableGoFieldType    string // full name of the nullable type
	nullableValueFieldName string // to access the go field value of a nullable type
	convertFuncName        string // to convert from a go field value to a nullable type
	newAccessFuncCall      string // to create the accessor
	isArray                bool   // true if the type is an array
}

func (m mapping) validate() error {
	if m.nullableGoFieldType == "" {
		return fmt.Errorf("nullableFieldType is required")
	}
	if m.newAccessFuncCall == "" {
		return fmt.Errorf("newAccessFuncName is required")
	}
	return nil
}

// https://www.postgresql.org/docs/9.1/datatype-numeric.html
var pgMappings = map[string]mapping{
	"timestamp with time zone": {
		goFieldType: "time.Time",
		newFuncCall: "p.NewTimeAccess",

		nullableValueFieldName: "Time",
		convertFuncName:        "c.TimeToTimestamptz",
		nullableGoFieldType:    "pgtype.Timestamptz",
		newAccessFuncCall:      "p.NewFieldAccess[pgtype.Timestamptz]",
	},
	"timestamp without time zone": {
		goFieldType: "time.Time",
		newFuncCall: "p.NewTimeAccess",

		nullableValueFieldName: "Time",
		convertFuncName:        "c.TimeToTimestamp",
		nullableGoFieldType:    "pgtype.Timestamp",
		newAccessFuncCall:      "p.NewFieldAccess[pgtype.Timestamp]",
	},
	"date": {
		goFieldType: "time.Time",
		newFuncCall: "p.NewTimeAccess",

		nullableValueFieldName: "Time",
		convertFuncName:        "c.TimeToDate",
		nullableGoFieldType:    "pgtype.Date",
		newAccessFuncCall:      "p.NewFieldAccess[pgtype.Date]",
	},
	"citext": {
		goFieldType: "string",
		newFuncCall: "p.NewTextAccess",

		nullableValueFieldName: "String",
		convertFuncName:        "c.StringToText",
		nullableGoFieldType:    "pgtype.Text",
		newAccessFuncCall:      "p.NewFieldAccess[pgtype.Text]",
	},
	"text": {
		goFieldType: "string",
		newFuncCall: "p.NewTextAccess",

		nullableValueFieldName: "String",
		convertFuncName:        "c.StringToText",
		nullableGoFieldType:    "pgtype.Text",
		newAccessFuncCall:      "p.NewFieldAccess[pgtype.Text]",
	},
	"character varying": {
		goFieldType: "string",
		newFuncCall: "p.NewTextAccess",

		nullableValueFieldName: "String",
		convertFuncName:        "c.StringToText",
		nullableGoFieldType:    "pgtype.Text",
		newAccessFuncCall:      "p.NewFieldAccess[pgtype.Text]",
	},
	"bigint": {
		goFieldType: "int64",
		newFuncCall: "p.NewInt64Access",

		nullableValueFieldName: "Int",
		convertFuncName:        "c.Int64ToInt8",
		nullableGoFieldType:    "pgtype.Int8",
		newAccessFuncCall:      "p.NewFieldAccess[pgtype.Int8]",
	},
	"integer": {
		goFieldType: "int32",
		newFuncCall: "p.NewInt32Access",

		nullableValueFieldName: "Int",
		convertFuncName:        "c.Int32ToInt4",
		nullableGoFieldType:    "pgtype.Int4",
		newAccessFuncCall:      "p.NewFieldAccess[pgtype.Int4]",
	},
	"jsonb": {
		nullableGoFieldType: "p.NullJSON",
		goFieldType:         "p.NullJSON",
		newAccessFuncCall:   "p.NewJSONAccess",
	},
	"json": {
		nullableGoFieldType: "p.NullJSON",
		goFieldType:         "p.NullJSON",
		newAccessFuncCall:   "p.NewJSONAccess",
	},
	"uuid": {
		goFieldType:         "pgtype.UUID",
		nullableGoFieldType: "pgtype.UUID",
		newAccessFuncCall:   "p.NewFieldAccess[pgtype.UUID]",
	},
	// https://github.com/jackc/pgx/wiki/Numeric-and-decimal-support
	"numeric": {
		goFieldType:         "decimal.NullDecimal",
		nullableGoFieldType: "decimal.NullDecimal",
		newAccessFuncCall:   "p.NewFieldAccess[decimal.NullDecimal]",
	},
	"double precision": {
		goFieldType:         "float64",
		nullableGoFieldType: "pgtype.Float8",
		newAccessFuncCall:   "p.NewFieldAccess[pgtype.Float8]",
	},
	"decimal": {
		goFieldType:         "decimal.NullDecimal",
		nullableGoFieldType: "decimal.NullDecimal",
		newAccessFuncCall:   "p.NewFieldAccess[decimal.NullDecimal]",
	},
	"point": {
		nullableGoFieldType: "pgtype.Point",
		newAccessFuncCall:   "p.NewFieldAccess[pgtype.Point]",
	},
	"boolean": {
		goFieldType: "bool",
		newFuncCall: "p.NewBooleanAccess",

		convertFuncName:        "Bool",
		nullableValueFieldName: "Bool",
		nullableGoFieldType:    "pgtype.Bool",
		newAccessFuncCall:      "p.NewFieldAccess[pgtype.Bool]",
	},
	"daterange": {
		nullableGoFieldType: "pgtype.Range[pgtype.Date]",
		goFieldType:         "pgtype.Range[pgtype.Date]",
		newAccessFuncCall:   "p.NewFieldAccess[pgtype.Range[pgtype.Date]]",
	},
	"bytea": {
		nullableGoFieldType: "pgtype.Bytea",
		newAccessFuncCall:   "p.NewFieldAccess[pgtype.Bytea]",
	},
	"text[]": {
		goFieldType:         "pgtype.FlatArray[pgtype.Text]",
		nullableGoFieldType: "pgtype.FlatArray[pgtype.Text]",
		newAccessFuncCall:   "p.NewFieldAccess[pgtype.FlatArray[pgtype.Text]]",
		isArray:             true,
	},
	"interval": {
		nullableGoFieldType: "pgtype.Interval",
		newAccessFuncCall:   "p.NewFieldAccess[pgtype.Interval]",
	},
}
