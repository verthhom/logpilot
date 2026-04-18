// Package multiline provides a line joiner that coalesces multi-line log
// events into a single string before they are handed to the rest of the
// logpilot pipeline.
//
// A new event begins whenever an incoming line matches a caller-supplied
// regular expression (e.g. a timestamp prefix or a JSON opening brace).
// Lines that do not match are appended to the current event buffer.
//
// Usage:
//
//	j, err := multiline.New(`^\{`, " ")
//	for _, raw := range lines {
//		if event, ok := j.Feed(raw); ok {
//			process(event)
//		}
//	}
//	if event, ok := j.Flush(); ok {
//		process(event)
//	}
package multiline
