// Package jsonslice extracts a sub-slice of a JSON array field.
package jsonslice

import (
	"encoding/json"
	"fmt"

	"github.com/logpilot/internal/jsonslice/internal/errors"
)

// Slicer extracts elements [Start:End] from a named JSON array field.
type Slicer struct {
	field string
	start int
	end   int
}

// New creates a Slicer that replaces field with arr[start:end].
// end == -1 means "until the last element" (open-ended).
func New(field string, start, end int) (*Slicer, error) {
	if field == "" {
		return nil, errors.ErrBlankField
	}
	if start < 0 {
		return nil, errors.ErrNegativeIndex
	}
	if end != -1 && end < start {
		return nil, errors.ErrEndBeforeStart
	}
	return &Slicer{field: field, start: start, end: end}, nil
}

// Apply returns line with the named array field replaced by the requested
// sub-slice. Non-JSON lines and lines where the field is not an array are
// passed through unchanged.
func (s *Slicer) Apply(line string) string {
	var obj map[string]json.RawMessage
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line
	}

	raw, ok := obj[s.field]
	if !ok {
		return line
	}

	var arr []json.RawMessage
	if err := json.Unmarshal(raw, &arr); err != nil {
		return line
	}

	lo := s.start
	hi := s.end
	if lo > len(arr) {
		lo = len(arr)
	}
	if hi == -1 || hi > len(arr) {
		hi = len(arr)
	}

	sliced, err := json.Marshal(arr[lo:hi])
	if err != nil {
		return line
	}
	obj[s.field] = sliced

	out, err := json.Marshal(obj)
	if err != nil {
		return line
	}
	return fmt.Sprintf("%s", out)
}
