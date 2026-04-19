package jsonstrip

import (
	"encoding/json"
	"fmt"
)

// Stripper removes specified keys from JSON log lines.
type Stripper struct {
	fields map[string]struct{}
}

// New creates a Stripper that removes the given fields from JSON objects.
// Returns an error if fields is empty or any field name is blank.
func New(fields []string) (*Stripper, error) {
	if len(fields) == 0 {
		return nil, fmt.Errorf("jsonstrip: at least one field required")
	}
	m := make(map[string]struct{}, len(fields))
	for _, f := range fields {
		if f == "" {
			return nil, fmt.Errorf("jsonstrip: field name must not be blank")
		}
		m[f] = struct{}{}
	}
	return &Stripper{fields: m}, nil
}

// Apply removes the configured fields from line if it is valid JSON.
// Non-JSON lines are returned unchanged.
func (s *Stripper) Apply(line string) string {
	var obj map[string]json.RawMessage
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line
	}
	for f := range s.fields {
		delete(obj, f)
	}
	out, err := json.Marshal(obj)
	if err != nil {
		return line
	}
	return string(out)
}
