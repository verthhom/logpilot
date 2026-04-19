// Package tagfield appends a static or dynamic tag to JSON log lines.
package tagfield

import (
	"encoding/json"
	"fmt"
)

// Tagger injects a tag field into JSON log lines.
type Tagger struct {
	field string
	value string
}

// New creates a Tagger that sets field to value on every JSON line.
// Returns an error if field or value is blank.
func New(field, value string) (*Tagger, error) {
	if field == "" {
		return nil, ErrBlankField
	}
	if value == "" {
		return nil, ErrBlankValue
	}
	return &Tagger{field: field, value: value}, nil
}

// Apply injects the tag field into the JSON line.
// Non-JSON lines are returned unchanged.
func (t *Tagger) Apply(line string) string {
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line
	}
	obj[t.field] = t.value
	b, err := json.Marshal(obj)
	if err != nil {
		return line
	}
	return fmt.Sprintf("%s", b)
}
