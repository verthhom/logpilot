// Package truncate provides a Truncator that shortens log lines exceeding
// a configured maximum byte length.
//
// Usage:
//
//	tr, err := truncate.New(200)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	short := tr.Apply(rawLine)
//
// Lines within the limit are returned unchanged. Lines that exceed the
// limit are trimmed and suffixed with "..." so the total length equals
// maxLen exactly.
package truncate
