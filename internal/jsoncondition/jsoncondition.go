// Package jsoncondition provides conditional field injection for JSON log lines.
// If a specified field matches a given value, a tag field is set to a target value.
package jsoncondition

import (
	"encoding/json"
	"strings"
)

// Conditioner applies a conditional transformation to JSON log lines.
type Conditioner struct {
	srcField  string
	srcValue  string
	destField string
	destValue string
}

// New creates a Conditioner from a rule string of the form:
//
//	"src_field=src_value:dest_field=dest_value"
//
// Returns an error if the rule is malformed or any part is blank.
func New(rule string) (*Conditioner, error) {
	parts := strings.SplitN(rule, ":", 2)
	if len(parts) != 2 {
		return nil, ErrBadRule
	}
	lhs := strings.SplitN(parts[0], "=", 2)
	rhs := strings.SplitN(parts[1], "=", 2)
	if len(lhs) != 2 || len(rhs) != 2 {
		return nil, ErrBadRule
	}
	for _, s := range []string{lhs[0], lhs[1], rhs[0], rhs[1]} {
		if strings.TrimSpace(s) == "" {
			return nil, ErrBlankPart
		}
	}
	return &Conditioner{
		srcField:  lhs[0],
		srcValue:  lhs[1],
		destField: rhs[0],
		destValue: rhs[1],
	}, nil
}

// Apply evaluates the condition against line. If the line is not valid JSON or
// the condition is not satisfied, the original line is returned unchanged.
func (c *Conditioner) Apply(line string) string {
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line
	}
	v, ok := obj[c.srcField]
	if !ok {
		return line
	}
	var strVal string
	switch val := v.(type) {
	case string:
		strVal = val
	default:
		return line
	}
	if strVal != c.srcValue {
		return line
	}
	obj[c.destField] = c.destValue
	out, err := json.Marshal(obj)
	if err != nil {
		return line
	}
	return string(out)
}
