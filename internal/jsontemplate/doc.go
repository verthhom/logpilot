// Package jsontemplate provides a log-line processor that evaluates a Go
// text/template against each structured JSON log entry and stores the rendered
// string as a new field.
//
// The template receives the full parsed JSON object as its dot value
// (map[string]any), so any top-level field can be referenced with the standard
// {{.fieldName}} syntax.
//
// Non-JSON lines are passed through without modification. Lines where template
// execution fails are also passed through unchanged so that the pipeline
// remains resilient to malformed data.
//
// Example
//
//	applier, err := jsontemplate.New("summary", "{{.level}}: {{.msg}}")
//	if err != nil {
//		log.Fatal(err)
//	}
//	result := applier.Apply(`{"level":"error","msg":"disk full"}`)
//	// result contains a "summary" field: "error: disk full"
package jsontemplate
