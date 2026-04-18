// Package filter provides rule-based filtering for structured JSON log entries.
//
// Rules are expressed as colon-separated strings in the form:
//
//	field:operator[:value]
//
// Supported operators:
//
//	eq       - field value equals the given string
//	contains - field value contains the given substring
//	exists   - field is present in the log entry (no value required)
//
// Multiple rules are combined with AND logic: an entry must satisfy all
// rules to be considered a match.
//
// Example usage:
//
//	f, err := filter.New([]string{"level:eq:error", "service:contains:auth"})
//	if err != nil { ... }
//	if f.Match(entry) {
//	    // process matching log line
//	}
package filter
