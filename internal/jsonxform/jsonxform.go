// Package jsonxform applies a Go text/template expression to transform
// a single JSON field value and write the result back to the same or a
// different field.
package jsonxform

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"text/template"
)

// Transformer rewrites one JSON field using a template expression.
type Transformer struct {
	src  string
	dest string
	tmpl *template.Template
}

// New returns a Transformer built from rule, which must have the form
// "dest=src:template" where template may reference {{.Value}} for the
// source field's current value.
func New(rule string) (*Transformer, error) {
	if rule == "" {
		return nil, ErrEmptyRule
	}

	eqIdx := strings.Index(rule, "=")
	if eqIdx < 1 {
		return nil, fmt.Errorf("%w: %q", ErrBadRule, rule)
	}
	dest := strings.TrimSpace(rule[:eqIdx])
	rest := rule[eqIdx+1:]

	colIdx := strings.Index(rest, ":")
	if colIdx < 1 {
		return nil, fmt.Errorf("%w: %q", ErrBadRule, rule)
	}
	src := strings.TrimSpace(rest[:colIdx])
	expr := strings.TrimSpace(rest[colIdx+1:])

	if dest == "" || src == "" || expr == "" {
		return nil, fmt.Errorf("%w: blank part in %q", ErrBadRule, rule)
	}

	tmpl, err := template.New("").Parse(expr)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrBadTemplate, err)
	}

	return &Transformer{src: src, dest: dest, tmpl: tmpl}, nil
}

// Apply transforms the source field and writes the result into dest.
// Non-JSON lines are passed through unchanged.
func (t *Transformer) Apply(line string) string {
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line
	}

	val, ok := obj[t.src]
	if !ok {
		return line
	}

	var buf bytes.Buffer
	if err := t.tmpl.Execute(&buf, map[string]interface{}{"Value": val}); err != nil {
		return line
	}

	obj[t.dest] = buf.String()

	out, err := json.Marshal(obj)
	if err != nil {
		return line
	}
	return string(out)
}
