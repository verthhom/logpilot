package jsontemplate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"text/template"
)

// Applier renders a Go text/template against each JSON log line,
// injecting the result as a new string field.
type Applier struct {
	field string
	tmpl  *template.Template
}

// New creates an Applier that evaluates tmplSrc as a Go text/template
// and stores the rendered output in field.
// The template receives the parsed JSON object as its dot value (map[string]any).
func New(field, tmplSrc string) (*Applier, error) {
	if strings.TrimSpace(field) == "" {
		return nil, ErrBlankField
	}
	if strings.TrimSpace(tmplSrc) == "" {
		return nil, ErrBlankTemplate
	}
	tmpl, err := template.New("jsontemplate").Option("missingkey=zero").Parse(tmplSrc)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrBadTemplate, err.Error())
	}
	return &Applier{field: field, tmpl: tmpl}, nil
}

// Apply renders the template against the JSON object in line and injects the
// result under the configured field name. Non-JSON lines pass through unchanged.
func (a *Applier) Apply(line string) string {
	var obj map[string]any
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line
	}

	var buf bytes.Buffer
	if err := a.tmpl.Execute(&buf, obj); err != nil {
		return line
	}

	obj[a.field] = buf.String()

	out, err := json.Marshal(obj)
	if err != nil {
		return line
	}
	return string(out)
}
