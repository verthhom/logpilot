package filter

import (
	"strings"
)

// Rule defines a single filter condition on a JSON log field.
type Rule struct {
	Field    string
	Operator string // eq, contains, exists
	Value    string
}

// Filter holds a set of rules applied with AND logic.
type Filter struct {
	Rules []Rule
}

// New creates a Filter from a slice of raw rule strings (field:op:value).
func New(rawRules []string) (*Filter, error) {
	f := &Filter{}
	for _, r := range rawRules {
		parts := strings.SplitN(r, ":", 3)
		if len(parts) < 2 {
			return nil, &InvalidRuleError{Raw: r}
		}
		rule := Rule{
			Field:    parts[0],
			Operator: parts[1],
		}
		if len(parts) == 3 {
			rule.Value = parts[2]
		}
		f.Rules = append(f.Rules, rule)
	}
	return f, nil
}

// Match returns true if the log entry (as a map) satisfies all filter rules.
func (f *Filter) Match(entry map[string]interface{}) bool {
	for _, rule := range f.Rules {
		val, exists := entry[rule.Field]
		switch rule.Operator {
		case "exists":
			if !exists {
				return false
			}
		case "eq":
			if !exists || toString(val) != rule.Value {
				return false
			}
		case "contains":
			if !exists || !strings.Contains(toString(val), rule.Value) {
				return false
			}
		default:
			return false
		}
	}
	return true
}

func toString(v interface{}) string {
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}
