// Package flatten provides a processor that flattens nested JSON objects
// into a single-level map using dot-notation keys.
package flatten

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Flattener flattens nested JSON log lines into dot-notation fields.
type Flattener struct {
	separator string
}

// New returns a Flattener using the given separator (e.g. ".").
func New(separator string) (*Flattener, error) {
	if separator == "" {
		return nil, fmt.Errorf("flatten: separator must not be empty")
	}
	return &Flattener{separator: separator}, nil
}

// Apply flattens a JSON line. Non-JSON lines are returned unchanged.
func (f *Flattener) Apply(line string) string {
	var obj map[string]any
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line
	}
	flat := make(map[string]any)
	f.flatten("", obj, flat)
	b, err := json.Marshal(flat)
	if err != nil {
		return line
	}
	return string(b)
}

func (f *Flattener) flatten(prefix string, obj map[string]any, out map[string]any) {
	for k, v := range obj {
		key := k
		if prefix != "" {
			key = strings.Join([]string{prefix, k}, f.separator)
		}
		switch child := v.(type) {
		case map[string]any:
			f.flatten(key, child, out)
		default:
			out[key] = v
		}
	}
}
