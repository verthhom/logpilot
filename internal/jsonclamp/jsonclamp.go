// Package jsonclamp clamps numeric fields in JSON log lines to a specified range.
package jsonclamp

import (
	"encoding/json"
	"fmt"
)

// Clamper clamps numeric JSON fields to [Min, Max].
type Clamper struct {
	rules []rule
}

type rule struct {
	field string
	min   float64
	max   float64
}

// New creates a Clamper from a slice of rule strings in the form "field:min:max".
func New(rules []string) (*Clamper, error) {
	if len(rules) == 0 {
		return nil, ErrNoRules
	}
	parsed := make([]rule, 0, len(rules))
	for _, r := range rules {
		var field string
		var min, max float64
		_, err := fmt.Sscanf(r, "%s", &field)
		if err != nil {
			return nil, fmt.Errorf("%w: %q", ErrBadRule, r)
		}
		n, err2 := fmt.Sscanf(r, "%s %f %f", &field, &min, &max)
		// Use colon-delimited parsing
		var f string
		n, err2 = fmt.Sscanf(r, "%[^:]:%f:%f", &f, &min, &max)
		if err2 != nil || n != 3 {
			return nil, fmt.Errorf("%w: %q", ErrBadRule, r)
		}
		if f == "" {
			return nil, fmt.Errorf("%w: %q", ErrBlankField, r)
		}
		if min > max {
			return nil, fmt.Errorf("%w: min %.4g > max %.4g in %q", ErrInvalidRange, min, max, r)
		}
		parsed = append(parsed, rule{field: f, min: min, max: max})
	}
	return &Clamper{rules: parsed}, nil
}

// Apply clamps numeric fields in the JSON line according to the configured rules.
// Non-JSON lines are returned unchanged.
func (c *Clamper) Apply(line string) string {
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line
	}
	for _, r := range c.rules {
		v, ok := obj[r.field]
		if !ok {
			continue
		}
		num, ok := v.(float64)
		if !ok {
			continue
		}
		if num < r.min {
			obj[r.field] = r.min
		} else if num > r.max {
			obj[r.field] = r.max
		}
	}
	out, err := json.Marshal(obj)
	if err != nil {
		return line
	}
	return string(out)
}
