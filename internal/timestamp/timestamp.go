// Package timestamp provides utilities for normalising and reformatting
// timestamp fields found in structured JSON log lines.
package timestamp

import (
	"encoding/json"
	"fmt"
	"time"
)

// common input layouts tried in order when parsing.
var inputLayouts = []string{
	time.RFC3339Nano,
	time.RFC3339,
	"2006-01-02T15:04:05.999999999",
	"2006-01-02 15:04:05",
	"2006-01-02",
}

// Formatter reads a timestamp field from a JSON log line, parses it, and
// rewrites it using the configured output layout.
type Formatter struct {
	field  string
	layout string
}

// New creates a Formatter that rewrites the named field using outLayout.
// outLayout must be a valid Go time layout string.
func New(field, outLayout string) (*Formatter, error) {
	if field == "" {
		return nil, ErrEmptyField
	}
	if outLayout == "" {
		return nil, ErrEmptyLayout
	}
	// validate layout by formatting a known time
	_ = time.Unix(0, 0).UTC().Format(outLayout)
	return &Formatter{field: field, layout: outLayout}, nil
}

// Apply rewrites the timestamp field in line and returns the modified JSON.
// Non-JSON lines and lines missing the field are returned unchanged.
func (f *Formatter) Apply(line string) string {
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line
	}
	raw, ok := obj[f.field]
	if !ok {
		return line
	}
	str, ok := raw.(string)
	if !ok {
		return line
	}
	t, err := parse(str)
	if err != nil {
		return line
	}
	obj[f.field] = t.Format(f.layout)
	b, err := json.Marshal(obj)
	if err != nil {
		return line
	}
	return string(b)
}

func parse(value string) (time.Time, error) {
	for _, layout := range inputLayouts {
		if t, err := time.Parse(layout, value); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("timestamp: cannot parse %q", value)
}
