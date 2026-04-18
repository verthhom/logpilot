// Package fieldmap provides utilities for remapping JSON log field names
// before output, allowing users to normalise fields across different log sources.
package fieldmap

import (
	"encoding/json"
	"strings"
)

// Mapper holds a set of field rename rules.
type Mapper struct {
	rules map[string]string
}

// New creates a Mapper from a slice of "old=new" rule strings.
// Returns an error if any rule is malformed.
func New(rules []string) (*Mapper, error) {
	m := &Mapper{rules: make(map[string]string, len(rules))}
	for _, r := range rules {
		parts := strings.SplitN(r, "=", 2)
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			return nil, &ErrInvalidRule{Rule: r}
		}
		m.rules[parts[0]] = parts[1]
	}
	return m, nil
}

// Apply remaps fields in a raw JSON line according to the configured rules.
// Non-JSON lines are returned unchanged.
func (m *Mapper) Apply(line string) string {
	if len(m.rules) == 0 {
		return line
	}
	var obj map[string]json.RawMessage
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line
	}
	for old, new := range m.rules {
		if val, ok := obj[old]; ok {
			delete(obj, old)
			obj[new] = val
		}
	}
	out, err := json.Marshal(obj)
	if err != nil {
		return line
	}
	return string(out)
}

// Rules returns a copy of the configured rename rules.
func (m *Mapper) Rules() map[string]string {
	copy := make(map[string]string, len(m.rules))
	for k, v := range m.rules {
		copy[k] = v
	}
	return copy
}
