// Package jsoncount injects a field containing the number of keys
// present in a JSON log line.
package jsoncount

import (
	"encoding/json"
	"fmt"
)

// Counter injects a count of top-level JSON keys into each log line.
type Counter struct {
	field string
}

// New returns a Counter that writes the key count into field.
// field must not be blank.
func New(field string) (*Counter, error) {
	if field == "" {
		return nil, fmt.Errorf("jsoncount: field must not be blank")
	}
	return &Counter{field: field}, nil
}

// Apply parses line as a JSON object, counts its top-level keys, and
// injects the count under the configured field name. Non-JSON lines
// are returned unchanged.
func (c *Counter) Apply(line string) string {
	var obj map[string]json.RawMessage
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line
	}

	obj[c.field] = json.RawMessage(fmt.Sprintf("%d", len(obj)))

	out, err := json.Marshal(obj)
	if err != nil {
		return line
	}
	return string(out)
}
