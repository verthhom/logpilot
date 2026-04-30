// Package jsonpresence checks whether required or forbidden JSON fields
// are present in a log line, passing or dropping it accordingly.
package jsonpresence

import (
	"encoding/json"
	"fmt"
)

// Checker holds the configured required and forbidden field names.
type Checker struct {
	required []string
	forbidden []string
}

// New creates a Checker. requiredFields must all be present; forbiddenFields
// must all be absent for a line to pass. Either slice may be empty.
func New(requiredFields, forbiddenFields []string) (*Checker, error) {
	for _, f := range requiredFields {
		if f == "" {
			return nil, fmt.Errorf("jsonpresence: required field name must not be blank")
		}
	}
	for _, f := range forbiddenFields {
		if f == "" {
			return nil, fmt.Errorf("jsonpresence: forbidden field name must not be blank")
		}
	}
	if len(requiredFields) == 0 && len(forbiddenFields) == 0 {
		return nil, fmt.Errorf("jsonpresence: at least one required or forbidden field must be specified")
	}
	return &Checker{
		required:  append([]string(nil), requiredFields...),
		forbidden: append([]string(nil), forbiddenFields...),
	}, nil
}

// Allow returns true when line satisfies the presence constraints.
// Non-JSON lines always return false.
func (c *Checker) Allow(line string) bool {
	var obj map[string]json.RawMessage
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return false
	}
	for _, f := range c.required {
		if _, ok := obj[f]; !ok {
			return false
		}
	}
	for _, f := range c.forbidden {
		if _, ok := obj[f]; ok {
			return false
		}
	}
	return true
}
