// Package jsonmerge merges a set of static key-value fields into every
// JSON log line that passes through the pipeline.  Non-JSON lines are
// passed through unchanged.
package jsonmerge

import (
	"encoding/json"
	"errors"
)

// ErrNoFields is returned when New is called with an empty fields map.
var ErrNoFields = errors.New("jsonmerge: at least one field is required")

// Merger merges static fields into JSON log lines.
type Merger struct {
	fields map[string]string
}

// New creates a Merger that will inject the supplied fields into every JSON
// log line.  The fields map must contain at least one entry.
func New(fields map[string]string) (*Merger, error) {
	if len(fields) == 0 {
		return nil, ErrNoFields
	}
	copy := make(map[string]string, len(fields))
	for k, v := range fields {
		copy[k] = v
	}
	return &Merger{fields: copy}, nil
}

// Apply merges the static fields into line.  Existing keys are NOT
// overwritten.  Non-JSON lines are returned as-is.
func (m *Merger) Apply(line string) string {
	var obj map[string]json.RawMessage
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line
	}
	for k, v := range m.fields {
		if _, exists := obj[k]; !exists {
			encoded, _ := json.Marshal(v)
			obj[k] = json.RawMessage(encoded)
		}
	}
	out, err := json.Marshal(obj)
	if err != nil {
		return line
	}
	return string(out)
}
