package jsoncoalesce

import (
	"encoding/json"

	"github.com/logpilot/internal/jsoncoalesce/errors"
)

// Coalescer picks the first non-empty value from a list of source fields
// and writes it to a destination field in a JSON log line.
type Coalescer struct {
	dest    string
	sources []string
}

// New returns a Coalescer that reads from sources in order and writes the
// first non-empty string value to dest.
func New(dest string, sources []string) (*Coalescer, error) {
	if dest == "" {
		return nil, ErrEmptyDest
	}
	if len(sources) == 0 {
		return nil, ErrNoSources
	}
	for _, s := range sources {
		if s == "" {
			return nil, ErrBlankSource
		}
	}
	return &Coalescer{dest: dest, sources: sources}, nil
}

// Apply returns line with dest set to the first non-empty value found among
// sources. If no source yields a value the line is returned unchanged.
func (c *Coalescer) Apply(line string) string {
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line
	}

	for _, src := range c.sources {
		val, ok := obj[src]
		if !ok {
			continue
		}
		str, ok := val.(string)
		if !ok || str == "" {
			continue
		}
		obj[c.dest] = str
		out, err := json.Marshal(obj)
		if err != nil {
			return line
		}
		return string(out)
	}
	return line
}
