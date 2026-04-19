// Package enrich provides field enrichment for JSON log lines,
// injecting static key-value pairs into every record.
package enrich

import (
	"encoding/json"
	"fmt"
)

// Enricher appends static fields to JSON log lines.
type Enricher struct {
	fields map[string]string
	keys   []string
}

// New creates an Enricher from a slice of "key=value" rules.
// Returns an error if any rule is malformed or contains a blank key.
func New(rules []string) (*Enricher, error) {
	if len(rules) == 0 {
		return nil, ErrNoRules
	}
	fields := make(map[string]string, len(rules))
	keys := make([]string, 0, len(rules))
	for _, r := range rules {
		k, v, err := splitRule(r)
		if err != nil {
			return nil, err
		}
		if _, dup := fields[k]; !dup {
			keys = append(keys, k)
		}
		fields[k] = v
	}
	return &Enricher{fields: fields, keys: keys}, nil
}

// Apply injects the static fields into line if it is valid JSON.
// Non-JSON lines are returned unchanged.
func (e *Enricher) Apply(line string) string {
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line
	}
	for _, k := range e.keys {
		if _, exists := obj[k]; !exists {
			obj[k] = e.fields[k]
		}
	}
	out, err := json.Marshal(obj)
	if err != nil {
		return line
	}
	return string(out)
}

func splitRule(rule string) (string, string, error) {
	for i, c := range rule {
		if c == '=' {
			k := rule[:i]
			if k == "" {
				return "", "", fmt.Errorf("%w: %q", ErrBlankKey, rule)
			}
			return k, rule[i+1:], nil
		}
	}
	return "", "", fmt.Errorf("%w: %q", ErrBadRule, rule)
}
