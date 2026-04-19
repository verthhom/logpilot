// Package jsonpath provides dot-notation field extraction from JSON log lines.
package jsonpath

import (
	"encoding/json"
	"fmt"
	"strings"
)

// ErrNotFound is returned when the path does not exist in the object.
var ErrNotFound = fmt.Errorf("jsonpath: field not found")

// Extractor resolves dot-notation paths against JSON objects.
type Extractor struct{}

// New returns a new Extractor.
func New() *Extractor { return &Extractor{} }

// Get returns the value at the dot-notation path within the JSON line.
// Returns ErrNotFound if any segment of the path is absent.
func (e *Extractor) Get(line, path string) (any, error) {
	var root map[string]any
	if err := json.Unmarshal([]byte(line), &root); err != nil {
		return nil, fmt.Errorf("jsonpath: invalid json: %w", err)
	}
	return walk(root, strings.Split(path, "."))
}

// GetString is a convenience wrapper that coerces the result to a string.
func (e *Extractor) GetString(line, path string) (string, error) {
	v, err := e.Get(line, path)
	if err != nil {
		return "", err
	}
	switch t := v.(type) {
	case string:
		return t, nil
	case float64:
		return fmt.Sprintf("%g", t), nil
	case bool:
		return fmt.Sprintf("%t", t), nil
	default:
		return fmt.Sprintf("%v", t), nil
	}
}

func walk(obj map[string]any, parts []string) (any, error) {
	v, ok := obj[parts[0]]
	if !ok {
		return nil, ErrNotFound
	}
	if len(parts) == 1 {
		return v, nil
	}
	nested, ok := v.(map[string]any)
	if !ok {
		return nil, ErrNotFound
	}
	return walk(nested, parts[1:])
}
