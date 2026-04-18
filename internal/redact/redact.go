// Package redact provides field-level redaction for structured JSON log lines.
package redact

import (
	"encoding/json"
	"errors"
	"strings"
)

const redactedValue = "[REDACTED]"

// Redactor replaces the values of specified JSON fields with a placeholder.
type Redactor struct {
	fields map[string]struct{}
}

// New creates a Redactor that will redact the given field names.
// Field names are case-sensitive. Returns an error if fields is empty.
func New(fields []string) (*Redactor, error) {
	if len(fields) == 0 {
		return nil, errors.New("redact: at least one field name is required")
	}
	m := make(map[string]struct{}, len(fields))
	for _, f := range fields {
		f = strings.TrimSpace(f)
		if f == "" {
			return nil, errors.New("redact: field name must not be blank")
		}
		m[f] = struct{}{}
	}
	return &Redactor{fields: m}, nil
}

// Apply redacts configured fields from a JSON log line.
// Non-JSON lines are returned unchanged.
func (r *Redactor) Apply(line string) string {
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line
	}
	for field := range r.fields {
		if _, ok := obj[field]; ok {
			obj[field] = redactedValue
		}
	}
	out, err := json.Marshal(obj)
	if err != nil {
		return line
	}
	return string(out)
}
