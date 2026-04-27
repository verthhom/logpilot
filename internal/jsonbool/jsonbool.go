// Package jsonbool normalises boolean-like values in JSON log lines.
// String values such as "true", "1", "yes" and "on" are converted to
// the JSON boolean true; "false", "0", "no" and "off" become false.
package jsonbool

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Normaliser rewrites nominated fields so that boolean-like string or
// numeric values become proper JSON booleans.
type Normaliser struct {
	fields map[string]struct{}
}

var truthy = map[string]bool{
	"true": true, "1": true, "yes": true, "on": true,
	"false": false, "0": false, "no": false, "off": false,
}

// New returns a Normaliser that operates on the supplied field names.
// At least one field must be provided and no field name may be blank.
func New(fields []string) (*Normaliser, error) {
	if len(fields) == 0 {
		return nil, fmt.Errorf("jsonbool: at least one field is required")
	}
	fm := make(map[string]struct{}, len(fields))
	for _, f := range fields {
		if strings.TrimSpace(f) == "" {
			return nil, fmt.Errorf("jsonbool: field name must not be blank")
		}
		fm[f] = struct{}{}
	}
	return &Normaliser{fields: fm}, nil
}

// Apply parses line as a JSON object, converts any targeted fields that
// hold boolean-like values to proper booleans, and returns the result.
// Lines that are not valid JSON objects are returned unchanged.
func (n *Normaliser) Apply(line string) string {
	var obj map[string]json.RawMessage
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line
	}

	changed := false
	for field := range n.fields {
		raw, ok := obj[field]
		if !ok {
			continue
		}
		// Already a proper boolean — nothing to do.
		if string(raw) == "true" || string(raw) == "false" {
			continue
		}
		// Try unquoting as a string.
		var s string
		if err := json.Unmarshal(raw, &s); err != nil {
			continue
		}
		bv, known := truthy[strings.ToLower(strings.TrimSpace(s))]
		if !known {
			continue
		}
		if bv {
			obj[field] = json.RawMessage("true")
		} else {
			obj[field] = json.RawMessage("false")
		}
		changed = true
	}

	if !changed {
		return line
	}
	out, err := json.Marshal(obj)
	if err != nil {
		return line
	}
	return string(out)
}
