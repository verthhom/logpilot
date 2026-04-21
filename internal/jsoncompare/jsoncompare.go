// Package jsoncompare provides a processor that compares two numeric fields
// and injects a result field indicating their relationship (lt, eq, gt).
package jsoncompare

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Comparer compares two numeric JSON fields and writes a result field.
type Comparer struct {
	left   string
	right  string
	dest   string
}

// New creates a Comparer from a rule string of the form "left:right:dest".
// The left and right fields must be numeric in the processed JSON objects.
// The dest field will be set to "lt", "eq", or "gt".
func New(rule string) (*Comparer, error) {
	parts := strings.SplitN(rule, ":", 3)
	if len(parts) != 3 {
		return nil, ErrBadRule
	}
	for _, p := range parts {
		if strings.TrimSpace(p) == "" {
			return nil, ErrBlankPart
		}
	}
	return &Comparer{
		left:  parts[0],
		right: parts[1],
		dest:  parts[2],
	}, nil
}

// Apply reads left and right numeric fields from the JSON line, compares them,
// and injects dest with "lt", "eq", or "gt". Non-JSON lines pass through unchanged.
func (c *Comparer) Apply(line string) string {
	var obj map[string]any
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line
	}

	lv, lok := toFloat(obj[c.left])
	rv, rok := toFloat(obj[c.right])
	if !lok || !rok {
		return line
	}

	switch {
	case lv < rv:
		obj[c.dest] = "lt"
	case lv > rv:
		obj[c.dest] = "gt"
	default:
		obj[c.dest] = "eq"
	}

	out, err := json.Marshal(obj)
	if err != nil {
		return line
	}
	return fmt.Sprintf("%s", out)
}

func toFloat(v any) (float64, bool) {
	switch n := v.(type) {
	case float64:
		return n, true
	case json.Number:
		f, err := n.Float64()
		return f, err == nil
	}
	return 0, false
}
