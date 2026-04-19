// Package jsonunwrap promotes a nested JSON object's fields to the top level.
package jsonunwrap

import (
	"encoding/json"
	"fmt"
)

// Unwrapper promotes fields from a nested key to the top-level object.
type Unwrapper struct {
	field string
}

// New returns an Unwrapper that promotes fields from the given nested key.
func New(field string) (*Unwrapper, error) {
	if field == "" {
		return nil, fmt.Errorf("jsonunwrap: field must not be empty")
	}
	return &Unwrapper{field: field}, nil
}

// Apply parses line as JSON, merges fields from the nested key into the root
// object (without overwriting existing keys), and returns the result.
// Non-JSON lines are returned unchanged.
func (u *Unwrapper) Apply(line string) string {
	var root map[string]json.RawMessage
	if err := json.Unmarshal([]byte(line), &root); err != nil {
		return line
	}

	raw, ok := root[u.field]
	if !ok {
		return line
	}

	var nested map[string]json.RawMessage
	if err := json.Unmarshal(raw, &nested); err != nil {
		return line
	}

	for k, v := range nested {
		if _, exists := root[k]; !exists {
			root[k] = v
		}
	}
	delete(root, u.field)

	out, err := json.Marshal(root)
	if err != nil {
		return line
	}
	return string(out)
}
