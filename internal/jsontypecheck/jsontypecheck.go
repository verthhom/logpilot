// Package jsontypecheck validates that JSON fields match expected types,
// passing through lines that conform and dropping (or flagging) those that do not.
package jsontypecheck

import (
	"encoding/json"
	"fmt"
	"strings"
)

// allowedTypes lists the type names accepted in rule definitions.
var allowedTypes = map[string]struct{}{
	"string":  {},
	"number":  {},
	"bool":    {},
	"array":   {},
	"object":  {},
	"null":    {},
}

// Checker validates JSON field types according to a set of rules.
type Checker struct {
	rules map[string]string // field -> expected type
}

// New creates a Checker from a slice of "field:type" rule strings.
// Returns an error if any rule is malformed or references an unknown type.
func New(rules []string) (*Checker, error) {
	if len(rules) == 0 {
		return nil, ErrNoRules
	}
	parsed := make(map[string]string, len(rules))
	for _, r := range rules {
		parts := strings.SplitN(r, ":", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("%w: %q", ErrBadRule, r)
		}
		field, typ := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
		if field == "" || typ == "" {
			return nil, fmt.Errorf("%w: %q", ErrBlankPart, r)
		}
		if _, ok := allowedTypes[typ]; !ok {
			return nil, fmt.Errorf("%w: %q", ErrUnknownType, typ)
		}
		parsed[field] = typ
	}
	return &Checker{rules: parsed}, nil
}

// Check returns the line unchanged if all governed fields match their expected
// types (or are absent). It returns an empty string when a field is present but
// has the wrong type.
func (c *Checker) Check(line string) string {
	var obj map[string]json.RawMessage
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line // non-JSON passes through
	}
	for field, want := range c.rules {
		raw, ok := obj[field]
		if !ok {
			continue // absent fields are not checked
		}
		if jsonTypeName(raw) != want {
			return ""
		}
	}
	return line
}

// jsonTypeName returns the JSON type name for a raw JSON value.
func jsonTypeName(raw json.RawMessage) string {
	if len(raw) == 0 {
		return "null"
	}
	switch raw[0] {
	case '"':
		return "string"
	case '{':
		return "object"
	case '[':
		return "array"
	case 't', 'f':
		return "bool"
	case 'n':
		return "null"
	default:
		return "number"
	}
}
