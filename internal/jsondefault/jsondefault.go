// Package jsondefault sets default values for missing or null JSON fields.
package jsondefault

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Defaulter applies default values to JSON log lines for fields that are
// absent or explicitly null.
type Defaulter struct {
	defaults map[string]json.RawMessage
}

// New creates a Defaulter from a slice of "field=value" rules.
// Each value is treated as a raw JSON literal (string, number, bool, etc.).
// Returns an error if any rule is malformed or the value is invalid JSON.
func New(rules []string) (*Defaulter, error) {
	if len(rules) == 0 {
		return nil, ErrNoRules
	}
	defaults := make(map[string]json.RawMessage, len(rules))
	for _, rule := range rules {
		parts := strings.SplitN(rule, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("%w: %q", ErrBadRule, rule)
		}
		field, value := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
		if field == "" {
			return nil, fmt.Errorf("%w: %q", ErrBlankField, rule)
		}
		if !json.Valid([]byte(value)) {
			return nil, fmt.Errorf("%w: field %q value %q", ErrInvalidJSON, field, value)
		}
		defaults[field] = json.RawMessage(value)
	}
	return &Defaulter{defaults: defaults}, nil
}

// Apply sets default values on the given JSON line for any field that is
// missing or null. Non-JSON lines are returned unchanged.
func (d *Defaulter) Apply(line string) string {
	var obj map[string]json.RawMessage
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line
	}
	changed := false
	for field, defVal := range d.defaults {
		if existing, ok := obj[field]; !ok || string(existing) == "null" {
			obj[field] = defVal
			changed = true
		}
	}
	if !changed {
		return line
	}
	out, err := json.Marshal(obj)
	if err != nil {
		return line
	}
	return string(out)
}
