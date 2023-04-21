package pgtalk

import "encoding/json"

// NullJSON is a value that can scan a Nullable value to an empty interface (any)
type NullJSON struct {
	Any   any
	Valid bool
}

// Scan implements the database/sql Scanner interface.
func (na *NullJSON) Scan(src any) error {
	if src == nil {
		*na = NullJSON{}
		return nil
	}
	var val any
	switch src := src.(type) {
	case string:
		err := json.Unmarshal([]byte(src), &val)
		if err != nil {
			return err
		}
	case []byte:
		err := json.Unmarshal(src, &val)
		if err != nil {
			return err
		}
	}
	*na = NullJSON{Any: val, Valid: true}
	return nil
}
