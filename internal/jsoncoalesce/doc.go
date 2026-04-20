// Package jsoncoalesce provides a transformer that selects the first
// non-empty string value from an ordered list of JSON source fields and
// writes it to a destination field.
//
// This is useful when log producers use inconsistent field names for the
// same semantic value (e.g. "msg", "message", "text") and you want to
// normalise them into a single canonical field before further processing.
//
// Usage:
//
//	c, err := jsoncoalesce.New("message", []string{"msg", "text", "body"})
//	if err != nil { ... }
//	output := c.Apply(line)
package jsoncoalesce
