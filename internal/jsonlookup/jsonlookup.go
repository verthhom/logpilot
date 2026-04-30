package jsonlookup

import (
	"encoding/json"
	"strings"
)

// Lookup replaces the value of a JSON field by looking it up in a static map.
// If the field value is not found in the map the line is passed through unchanged.
type Lookup struct {
	field  string
	table  map[string]string
	destField string
}

// New creates a Lookup that reads field, consults table, and writes the result
// to destField. If destField is empty the source field is overwritten.
func New(field string, table map[string]string, destField string) (*Lookup, error) {
	if strings.TrimSpace(field) == "" {
		return nil, ErrBlankField
	}
	if len(table) == 0 {
		return nil, ErrEmptyTable
	}
	dest := destField
	if strings.TrimSpace(dest) == "" {
		dest = field
	}
	return &Lookup{field: field, table: table, destField: dest}, nil
}

// Apply performs the lookup on a single log line. Non-JSON lines are returned
// unchanged. Lines whose field value is not in the table are returned unchanged.
func (l *Lookup) Apply(line string) string {
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line
	}

	raw, ok := obj[l.field]
	if !ok {
		return line
	}

	var key string
	switch v := raw.(type) {
	case string:
		key = v
	default:
		return line
	}

	mapped, found := l.table[key]
	if !found {
		return line
	}

	obj[l.destField] = mapped

	out, err := json.Marshal(obj)
	if err != nil {
		return line
	}
	return string(out)
}
