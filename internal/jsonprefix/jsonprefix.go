// Package jsonprefix adds a string prefix to every key in a JSON log line.
// Keys that already carry the prefix are left unchanged.
package jsonprefix

import (
	"encoding/json"
	"strings"
)

// Prefixer adds a configurable string prefix to every top-level key of a
// JSON object. Non-JSON lines are passed through without modification.
type Prefixer struct {
	prefix string
}

// New returns a Prefixer that will prepend prefix to every top-level key.
// An empty prefix is rejected with ErrEmptyPrefix.
func New(prefix string) (*Prefixer, error) {
	if strings.TrimSpace(prefix) == "" {
		return nil, ErrEmptyPrefix
	}
	return &Prefixer{prefix: prefix}, nil
}

// Apply rewrites line so that every top-level JSON key is prefixed.
// If the line is not valid JSON, it is returned unchanged.
func (p *Prefixer) Apply(line string) string {
	var obj map[string]json.RawMessage
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line
	}

	out := make(map[string]json.RawMessage, len(obj))
	for k, v := range obj {
		newKey := k
		if !strings.HasPrefix(k, p.prefix) {
			newKey = p.prefix + k
		}
		out[newKey] = v
	}

	b, err := json.Marshal(out)
	if err != nil {
		return line
	}
	return string(b)
}
