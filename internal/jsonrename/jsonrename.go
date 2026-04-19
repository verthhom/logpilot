package jsonrename

import (
	"encoding/json"
	"strings"
)

// Renamer renames JSON fields according to a set of from:to rules.
type Renamer struct {
	rules map[string]string
}

// New creates a Renamer from a slice of "from:to" rule strings.
func New(rules []string) (*Renamer, error) {
	if len(rules) == 0 {
		return nil, ErrNoRules
	}
	m := make(map[string]string, len(rules))
	for _, r := range rules {
		parts := strings.SplitN(r, ":", 2)
		if len(parts) != 2 || strings.TrimSpace(parts[0]) == "" || strings.TrimSpace(parts[1]) == "" {
			return nil, ErrInvalidRule(r)
		}
		m[parts[0]] = parts[1]
	}
	return &Renamer{rules: m}, nil
}

// Apply renames fields in a JSON log line. Non-JSON lines pass through unchanged.
func (r *Renamer) Apply(line string) string {
	var obj map[string]json.RawMessage
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line
	}
	for from, to := range r.rules {
		if val, ok := obj[from]; ok {
			obj[to] = val
			delete(obj, from)
		}
	}
	out, err := json.Marshal(obj)
	if err != nil {
		return line
	}
	return string(out)
}
