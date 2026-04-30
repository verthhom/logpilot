// Package jsonhash computes a hash of specified JSON fields and injects
// the result as a new field in the log line.
package jsonhash

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strings"
)

// Hasher computes a SHA-256 digest over one or more JSON field values
// and writes the hex-encoded result into a destination field.
type Hasher struct {
	fields []string
	dest   string
}

// New returns a Hasher that hashes the given source fields and stores
// the result in dest. At least one source field must be provided and
// dest must not be blank.
func New(dest string, fields []string) (*Hasher, error) {
	if strings.TrimSpace(dest) == "" {
		return nil, errBlankDest
	}
	if len(fields) == 0 {
		return nil, errNoFields
	}
	for _, f := range fields {
		if strings.TrimSpace(f) == "" {
			return nil, errBlankField
		}
	}
	return &Hasher{fields: fields, dest: dest}, nil
}

// Apply reads the configured fields from the JSON line, concatenates
// their raw values, hashes the result, and injects the hex digest into
// dest. Non-JSON lines are returned unchanged.
func (h *Hasher) Apply(line string) string {
	var obj map[string]json.RawMessage
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line
	}

	var sb strings.Builder
	for _, f := range h.fields {
		if v, ok := obj[f]; ok {
			sb.Write(v)
		}
	}

	sum := sha256.Sum256([]byte(sb.String()))
	obj[h.dest] = json.RawMessage(fmt.Sprintf("%q", fmt.Sprintf("%x", sum)))

	out, err := json.Marshal(obj)
	if err != nil {
		return line
	}
	return string(out)
}
