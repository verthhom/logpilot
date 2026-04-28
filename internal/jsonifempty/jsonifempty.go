// Package jsonifempty sets a fallback value on a JSON field when the field
// is missing, null, or an empty string.
package jsonifempty

import (
	"encoding/json"
	"strings"

	"github.com/logpilot/logpilot/internal/jsonifempty/internal/errors"
)

// Rule pairs a field name with the raw JSON fallback value to inject.
type Rule struct {
	Field    string
	Fallback json.RawMessage
}

// Filler applies if-empty fallback rules to JSON log lines.
type Filler struct {
	rules []Rule
}

// New constructs a Filler from a slice of "field=value" rule strings.
// value must be valid JSON (e.g. `"unknown"`, `0`, `true`).
func New(rules []string) (*Filler, error) {
	if len(rules) == 0 {
		return nil, errors.ErrNoRules
	}
	parsed := make([]Rule, 0, len(rules))
	for _, r := range rules {
		idx := strings.IndexByte(r, '=')
		if idx < 0 {
			return nil, errors.ErrBadRule
		}
		field := strings.TrimSpace(r[:idx])
		raw := strings.TrimSpace(r[idx+1:])
		if field == "" || raw == "" {
			return nil, errors.ErrBlankPart
		}
		var probe interface{}
		if err := json.Unmarshal([]byte(raw), &probe); err != nil {
			return nil, errors.ErrInvalidJSON
		}
		parsed = append(parsed, Rule{Field: field, Fallback: json.RawMessage(raw)})
	}
	return &Filler{rules: parsed}, nil
}

// Apply returns line unchanged when it is not valid JSON. For each rule, if
// the target field is absent, null, or an empty string the fallback value is
// written into the object.
func (f *Filler) Apply(line string) string {
	var obj map[string]json.RawMessage
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line
	}
	changed := false
	for _, r := range f.rules {
		if shouldFill(obj[r.Field]) {
			obj[r.Field] = r.Fallback
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

// shouldFill reports whether a raw JSON value is absent, null, or "".
func shouldFill(v json.RawMessage) bool {
	if v == nil {
		return true
	}
	s := strings.TrimSpace(string(v))
	return s == "null" || s == `""`
}
