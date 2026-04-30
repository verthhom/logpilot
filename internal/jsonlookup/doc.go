// Package jsonlookup provides a transformer that replaces a JSON field's value
// by consulting a static key→value table.
//
// # Usage
//
//	table := map[string]string{
//		"us-east-1": "US East (N. Virginia)",
//		"eu-west-1": "EU West (Ireland)",
//	}
//	l, err := jsonlookup.New("region", table, "region_name")
//	if err != nil { ... }
//	out := l.Apply(line)
//
// If the field value is absent from the table the line is passed through
// unchanged. Non-JSON lines are also passed through unchanged.
package jsonlookup
