// Package jsonxform provides field-level text/template transformation for
// structured JSON log lines.
//
// A rule has the form:
//
//	dest=src:template
//
// where dest is the output field name, src is the source field to read, and
// template is a Go text/template expression that may reference {{.Value}} for
// the source value.  The rendered string is always stored as a JSON string.
//
// Example:
//
//	xform, _ := jsonxform.New("msg_upper=msg:{{printf "%s" .Value | upper}}")
//	out := xform.Apply(`{"msg":"hello","level":"info"}`)
//	// out: {"level":"info","msg":"hello","msg_upper":"hello"}
package jsonxform
