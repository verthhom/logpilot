// Package jsonsort reorders fields in a JSON log line according to a
// caller-supplied field order. Fields not listed appear after the ordered
// fields in their original relative order.
package jsonsort

import (
	"encoding/json"
	"errors"
)

// ErrNoFields is returned when New is called with an empty field list.
var ErrNoFields = errors.New("jsonsort: at least one field required")

// ErrBlankField is returned when an empty string appears in the field list.
var ErrBlankField = errors.New("jsonsort: field name must not be blank")

// Sorter reorders JSON fields.
type Sorter struct {
	order []string
	index map[string]int
}

// New creates a Sorter that places fields in the given order first.
func New(fields []string) (*Sorter, error) {
	if len(fields) == 0 {
		return nil, ErrNoFields
	}
	idx := make(map[string]int, len(fields))
	for i, f := range fields {
		if f == "" {
			return nil, ErrBlankField
		}
		idx[f] = i
	}
	return &Sorter{order: fields, index: idx}, nil
}

// Apply reorders the fields of a JSON object line. Non-JSON lines are
// returned unchanged.
func (s *Sorter) Apply(line string) string {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal([]byte(line), &raw); err != nil {
		return line
	}

	out := make([]byte, 0, len(line))
	out = append(out, '{')

	written := make(map[string]bool, len(raw))
	first := true

	appendField := func(k string, v json.RawMessage) {
		if !first {
			out = append(out, ',')
		}
		first = false
		key, _ := json.Marshal(k)
		out = append(out, key...)
		out = append(out, ':')
		out = append(out, v...)
		written[k] = true
	}

	for _, f := range s.order {
		if v, ok := raw[f]; ok {
			appendField(f, v)
		}
	}

	for k, v := range raw {
		if !written[k] {
			appendField(k, v)
		}
	}

	out = append(out, '}')
	return string(out)
}
