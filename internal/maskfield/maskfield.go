// Package maskfield partially masks JSON field values for privacy.
package maskfield

import (
	"encoding/json"
	"strings"
)

// Masker partially masks specified JSON string fields, keeping a configurable
// number of prefix characters visible and replacing the rest with asterisks.
type Masker struct {
	fields  map[string]bool
	visible int
}

// New creates a Masker that masks the given fields, keeping visible prefix chars.
func New(fields []string, visible int) (*Masker, error) {
	if len(fields) == 0 {
		return nil, ErrNoFields
	}
	if visible < 0 {
		return nil, ErrNegativeVisible
	}
	set := make(map[string]bool, len(fields))
	for _, f := range fields {
		if strings.TrimSpace(f) == "" {
			return nil, ErrBlankField
		}
		set[f] = true
	}
	return &Masker{fields: set, visible: visible}, nil
}

// Apply masks targeted fields in the JSON line and returns the result.
// Non-JSON lines are returned unchanged.
func (m *Masker) Apply(line string) string {
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line
	}
	for k := range obj {
		if !m.fields[k] {
			continue
		}
		s, ok := obj[k].(string)
		if !ok {
			continue
		}
		obj[k] = mask(s, m.visible)
	}
	out, err := json.Marshal(obj)
	if err != nil {
		return line
	}
	return string(out)
}

func mask(s string, visible int) string {
	runes := []rune(s)
	if visible >= len(runes) {
		return s
	}
	return string(runes[:visible]) + strings.Repeat("*", len(runes)-visible)
}
