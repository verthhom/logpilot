// Package jsonparseerr provides a transformer that injects a field
// indicating whether a log line is valid JSON, allowing downstream
// processors to distinguish well-formed records from raw text.
package jsonparseerr

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Transformer injects a boolean field into each JSON log line
// indicating parse validity. Non-JSON lines are passed through with
// the error field appended as a plain key=value suffix.
type Transformer struct {
	field string
}

// New constructs a Transformer that writes parse status into field.
// field must be a non-blank string.
func New(field string) (*Transformer, error) {
	if strings.TrimSpace(field) == "" {
		return nil, fmt.Errorf("jsonparseerr: field name must not be blank")
	}
	return &Transformer{field: field}, nil
}

// Apply examines line. If it is valid JSON, it injects field:true and
// returns the re-encoded object. If it is not valid JSON, it returns
// the original line unchanged with field:false appended as JSON would
// not be constructable — instead a sentinel suffix is added so the
// line remains useful in plain-text pipelines.
func (t *Transformer) Apply(line string) string {
	var obj map[string]json.RawMessage
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		// Not valid JSON — return as-is with a plain suffix.
		return line + fmt.Sprintf(" %s=false", t.field)
	}

	obj[t.field] = json.RawMessage("true")

	out, err := json.Marshal(obj)
	if err != nil {
		return line
	}
	return string(out)
}
