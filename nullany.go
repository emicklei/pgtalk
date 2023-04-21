package pgtalk

// NullAny is a value that can scan a Nullable value to an empty interface (any)
type NullAny struct {
	Any   any
	Valid bool
}

// Scan implements the database/sql Scanner interface.
func (na *NullAny) Scan(src any) error {
	if src == nil {
		*na = NullAny{}
		return nil
	}
	*na = NullAny{Any: src, Valid: true}
	return nil
}
