// Package jsonround rounds numeric fields in JSON log lines to a given
// number of decimal places. Non-numeric fields and non-JSON lines are
// passed through unchanged.
package jsonround

import (
	"encoding/json"
	"fmt"
	"math"
	"strings"

	"github.com/logpilot/internal/jsonround/internal"
)

// Rounder rounds configured numeric fields to a fixed number of decimal places.
type Rounder struct {
	rules []rule
}

type rule struct {
	field  string
	places int
}

// New creates a Rounder from a slice of rules in "field:places" format.
// Each places value must be >= 0.
func New(rules []string) (*Rounder, error) {
	if len(rules) == 0 {
		return nil, internal.ErrNoRules
	}
	parsed := make([]rule, 0, len(rules))
	for _, r := range rules {
		idx := strings.LastIndex(r, ":")
		if idx < 1 {
			return nil, fmt.Errorf("%w: %q", internal.ErrBadRule, r)
		}
		field := strings.TrimSpace(r[:idx])
		if field == "" {
			return nil, fmt.Errorf("%w: blank field in %q", internal.ErrBadRule, r)
		}
		var places int
		if _, err := fmt.Sscanf(strings.TrimSpace(r[idx+1:]), "%d", &places); err != nil {
			return nil, fmt.Errorf("%w: non-integer places in %q", internal.ErrBadRule, r)
		}
		if places < 0 {
			return nil, fmt.Errorf("%w: negative places in %q", internal.ErrBadRule, r)
		}
		parsed = append(parsed, rule{field: field, places: places})
	}
	return &Rounder{rules: parsed}, nil
}

// Apply rounds the configured fields in line and returns the result.
// If line is not valid JSON it is returned as-is.
func (r *Rounder) Apply(line string) string {
	var obj map[string]json.RawMessage
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line
	}
	for _, ru := range r.rules {
		raw, ok := obj[ru.field]
		if !ok {
			continue
		}
		var f float64
		if err := json.Unmarshal(raw, &f); err != nil {
			continue
		}
		scale := math.Pow(10, float64(ru.places))
		rounded := math.Round(f*scale) / scale
		enc, _ := json.Marshal(rounded)
		obj[ru.field] = json.RawMessage(enc)
	}
	out, err := json.Marshal(obj)
	if err != nil {
		return line
	}
	return string(out)
}
