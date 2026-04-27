// Package jsonclean removes null, empty-string, and optionally empty-object
// or empty-array fields from a JSON log line.
package jsonclean

import (
	"encoding/json"
)

// Options controls which zero-value kinds are stripped.
type Options struct {
	RemoveNull        bool
	RemoveEmptyString bool
	RemoveEmptyArray  bool
	RemoveEmptyObject bool
}

// Cleaner strips unwanted zero-value fields from JSON lines.
type Cleaner struct {
	opts Options
}

// New returns a Cleaner configured with opts.
func New(opts Options) *Cleaner {
	return &Cleaner{opts: opts}
}

// Apply removes zero-value fields from line according to the configured
// options. Non-JSON lines are returned unchanged.
func (c *Cleaner) Apply(line string) string {
	var obj map[string]json.RawMessage
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line
	}

	for k, v := range obj {
		if c.shouldRemove(v) {
			delete(obj, k)
		}
	}

	out, err := json.Marshal(obj)
	if err != nil {
		return line
	}
	return string(out)
}

func (c *Cleaner) shouldRemove(v json.RawMessage) bool {
	s := string(v)
	if c.opts.RemoveNull && s == "null" {
		return true
	}
	if c.opts.RemoveEmptyString && s == `""` {
		return true
	}
	if c.opts.RemoveEmptyArray && s == "[]" {
		return true
	}
	if c.opts.RemoveEmptyObject && s == "{}" {
		return true
	}
	return false
}
