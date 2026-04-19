// Package jsonrename provides a processor that renames fields in structured
// JSON log lines according to a configurable set of "from:to" rules.
//
// Usage:
//
//	r, err := jsonrename.New([]string{"msg:message", "ts:timestamp"})
//	if err != nil { ... }
//	output := r.Apply(line)
//
// Non-JSON lines are passed through without modification.
package jsonrename
