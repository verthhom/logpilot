// Package jsonspread spreads the key-value pairs of a nested JSON object
// stored under a given field into the top-level document, then removes the
// original field.
package jsonspread

import (
	"encoding/json"
	"fmt"
)

// Spreader promotes the contents of a nested object field to the top level.
type Spreader struct {
	fields []string
}

// New returns a Spreader that will spread each of the named fields.
// Every entry in fields must be a non-blank string.
func New(fields []string) (*Spreader, error) {
	if len(fields) == 0 {
		return nil, fmt.Errorf("jsonspread: at least one field is required")
	}
	for _, f := range fields {
		if f == "" {
			return nil, fmt.Errorf("jsonspread: field name must not be blank")
		}
	}
	out := make([]string, len(fields))
	copy(out, fields)
	return &Spreader{fields: out}, nil
}

// Apply spreads each configured field into the top-level document.
// Lines that are not valid JSON objects are returned unchanged.
func (s *Spreader) Apply(line string) string {
	var doc map[string]json.RawMessage
	if err := json.Unmarshal([]byte(line), &doc); err != nil {
		return line
	}

	for _, field := range s.fields {
		raw, ok := doc[field]
		if !ok {
			continue
		}
		var nested map[string]json.RawMessage
		if err := json.Unmarshal(raw, &nested); err != nil {
			// Not an object — leave field untouched.
			continue
		}
		delete(doc, field)
		for k, v := range nested {
			if _, exists := doc[k]; !exists {
				doc[k] = v
			}
		}
	}

	b, err := json.Marshal(doc)
	if err != nil {
		return line
	}
	return string(b)
}
