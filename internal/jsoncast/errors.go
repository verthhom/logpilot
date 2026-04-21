package jsoncast

import "fmt"

// ErrNoRules is returned when no cast rules are provided.
var ErrNoRules = fmt.Errorf("jsoncast: at least one rule is required")

// ErrBadRule is returned when a rule string is malformed.
type ErrBadRule string

func (e ErrBadRule) Error() string {
	return fmt.Sprintf("jsoncast: bad rule %q: expected \"field:type\"", string(e))
}

// ErrUnknownType is returned when the target type is not supported.
type ErrUnknownType string

func (e ErrUnknownType) Error() string {
	return fmt.Sprintf("jsoncast: unknown type %q: must be string, int, float, or bool", string(e))
}
