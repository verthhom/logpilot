package jsoncast

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// Caster rewrites JSON field values to a target type.
type Caster struct {
	rules []rule
}

type rule struct {
	field    string
	targetType string
}

// New creates a Caster from a slice of "field:type" rule strings.
// Supported types: string, int, float, bool.
func New(rules []string) (*Caster, error) {
	if len(rules) == 0 {
		return nil, ErrNoRules
	}
	parsed := make([]rule, 0, len(rules))
	for _, r := range rules {
		field, typ, err := splitRule(r)
		if err != nil {
			return nil, err
		}
		parsed = append(parsed, rule{field: field, targetType: typ})
	}
	return &Caster{rules: parsed}, nil
}

// Apply casts fields in a JSON log line according to configured rules.
// Non-JSON lines are returned unchanged.
func (c *Caster) Apply(line string) string {
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line
	}
	for _, r := range c.rules {
		val, ok := obj[r.field]
		if !ok {
			continue
		}
		casted, err := castValue(val, r.targetType)
		if err != nil {
			continue
		}
		obj[r.field] = casted
	}
	out, err := json.Marshal(obj)
	if err != nil {
		return line
	}
	return string(out)
}

func castValue(val interface{}, typ string) (interface{}, error) {
	str := fmt.Sprintf("%v", val)
	switch typ {
	case "string":
		return str, nil
	case "int":
		f, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return nil, err
		}
		return int64(f), nil
	case "float":
		return strconv.ParseFloat(str, 64)
	case "bool":
		return strconv.ParseBool(str)
	default:
		return nil, fmt.Errorf("unsupported type: %s", typ)
	}
}

func splitRule(r string) (string, string, error) {
	for i, ch := range r {
		if ch == ':' {
			field := r[:i]
			typ := r[i+1:]
			if field == "" || typ == "" {
				return "", "", ErrBadRule(r)
			}
			switch typ {
			case "string", "int", "float", "bool":
				return field, typ, nil
			default:
				return "", "", ErrUnknownType(typ)
			}
		}
	}
	return "", "", ErrBadRule(r)
}
