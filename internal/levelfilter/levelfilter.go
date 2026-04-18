// Package levelfilter provides log-level-based filtering for structured JSON logs.
package levelfilter

import (
	"encoding/json"
	"strings"
)

// priority maps level strings to numeric priority.
var priority = map[string]int{
	"trace": 0,
	"debug": 1,
	"info":  2,
	"warn":  3,
	"warning": 3,
	"error": 4,
	"fatal": 5,
}

// Filter drops log lines whose level is below the configured minimum.
type Filter struct {
	minPriority int
	field       string
}

// New creates a Filter that passes lines at or above minLevel.
// field is the JSON key used to read the level (e.g. "level").
func New(minLevel, field string) (*Filter, error) {
	if field == "" {
		return nil, ErrEmptyField
	}
	p, ok := priority[strings.ToLower(minLevel)]
	if !ok {
		return nil, ErrUnknownLevel
	}
	return &Filter{minPriority: p, field: field}, nil
}

// Allow returns true if the line should be kept.
func (f *Filter) Allow(line string) bool {
	var obj map[string]any
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		// non-JSON lines pass through
		return true
	}
	val, ok := obj[f.field]
	if !ok {
		return true
	}
	s, ok := val.(string)
	if !ok {
		return true
	}
	p, ok := priority[strings.ToLower(s)]
	if !ok {
		return true
	}
	return p >= f.minPriority
}
