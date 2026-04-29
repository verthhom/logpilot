// Package jsonregex provides a Replacer that applies regular-expression
// substitutions to nominated fields in structured JSON log lines.
//
// Rules are expressed as strings in the form:
//
//	"field=pattern:replacement"
//
// where pattern is a Go regular expression and replacement follows the
// syntax accepted by regexp.ReplaceAllString (e.g. "$1" for group
// references).  Non-JSON lines are passed through unmodified.  Fields
// that are absent or whose value is not a string are silently skipped.
//
// Example:
//
//	rep, err := jsonregex.New([]string{"message=error:\\d+:ERROR"})
//	if err != nil { ... }
//	out := rep.Apply(`{"message":"error:42 occurred"}`)
//	// out → {"message":"ERROR occurred"}
package jsonregex
