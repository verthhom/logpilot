package jsonregex

import (
	"encoding/json"
	"regexp"
	"strings"
)

// Replacer replaces field values matching a regular expression with a
// substitution string. Non-JSON lines are passed through unchanged.
type Replacer struct {
	rules []rule
}

type rule struct {
	field string
	re    *regexp.Regexp
	subst string
}

// New creates a Replacer from a slice of rules in the form "field=pattern:replacement".
// Returns an error if any rule is malformed or contains an invalid pattern.
func New(rules []string) (*Replacer, error) {
	if len(rules) == 0 {
		return nil, ErrNoRules
	}
	parsed := make([]rule, 0, len(rules))
	for _, r := range rules {
		eq := strings.IndexByte(r, '=')
		if eq < 1 {
			return nil, ErrBadRule(r)
		}
		field := r[:eq]
		rest := r[eq+1:]
		colon := strings.IndexByte(rest, ':')
		if colon < 0 {
			return nil, ErrBadRule(r)
		}
		pattern := rest[:colon]
		subst := rest[colon+1:]
		if strings.TrimSpace(field) == "" || pattern == "" {
			return nil, ErrBadRule(r)
		}
		re, err := regexp.Compile(pattern)
		if err != nil {
			return nil, ErrBadPattern{Rule: r, Err: err}
		}
		parsed = append(parsed, rule{field: field, re: re, subst: subst})
	}
	return &Replacer{rules: parsed}, nil
}

// Apply processes a single log line, replacing matched field values.
// Non-JSON lines are returned unchanged.
func (rep *Replacer) Apply(line string) string {
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line
	}
	for _, r := range rep.rules {
		v, ok := obj[r.field]
		if !ok {
			continue
		}
		s, ok := v.(string)
		if !ok {
			continue
		}
		obj[r.field] = r.re.ReplaceAllString(s, r.subst)
	}
	out, err := json.Marshal(obj)
	if err != nil {
		return line
	}
	return string(out)
}
