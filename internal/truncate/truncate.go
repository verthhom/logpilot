// Package truncate provides line truncation for long log entries.
package truncate

import "errors"

const ellipsis = "..."

// Truncator trims lines that exceed a maximum byte length.
type Truncator struct {
	maxLen int
}

// New returns a Truncator that trims lines longer than maxLen bytes.
// maxLen must be greater than len(ellipsis).
func New(maxLen int) (*Truncator, error) {
	if maxLen <= len(ellipsis) {
		return nil, errors.New("truncate: maxLen must be greater than 3")
	}
	return &Truncator{maxLen: maxLen}, nil
}

// Apply returns line unchanged if it fits within maxLen, otherwise it
// trims the line and appends an ellipsis marker.
func (t *Truncator) Apply(line string) string {
	if len(line) <= t.maxLen {
		return line
	}
	return line[:t.maxLen-len(ellipsis)] + ellipsis
}

// Enabled reports whether truncation is active (always true for a
// constructed Truncator; provided for interface symmetry).
func (t *Truncator) Enabled() bool { return true }
