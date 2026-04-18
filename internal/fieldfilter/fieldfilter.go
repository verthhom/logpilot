// Package fieldfilter provides a processor that retains only specified fields
// from structured JSON log lines, dropping all others.
package fieldfilter

import (
	"encoding/json"
	"fmt"
)

// Filter retains only the configured fields from each JSON log line.
type Filter struct {
	fields map[string]struct{}
}

// New creates a Filter that keeps only the given fields.
// Returns an error if fields is empty or any field name is blank.
func New(fields []string) (*Filter, error) {
	if len(fields) == 0 {
		return nil, ErrNoFields
	}
	m := make(map[string]struct{}, len(fields))
	for _, f := range fields {
		if f == "" {
			return nil, fmt.Errorf("%w: field name must not be blank", ErrInvalidField)
		}
		m[f] = struct{}{}
	}
	return &Filter{fields: m}, nil
}

// Apply returns a new JSON line containing only the retained fields.
// Non-JSON lines are returned unchanged.
func (f *Filter) Apply(line string) string {
	var obj map[string]json.RawMessage
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line
	}
	out := make(map[string]json.RawMessage, len(f.fields))
	for k, v := range obj {
		if _, ok := f.fields[k]; ok {
			out[k] = v
		}
	}
	b, err := json.Marshal(out)
	if err != nil {
		return line
	}
	return string(b)
}
