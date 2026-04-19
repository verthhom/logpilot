// Package jsonstrip provides a Stripper that removes specified keys from
// structured JSON log lines.
//
// Usage:
//
//	s, err := jsonstrip.New([]string{"password", "token"})
//	if err != nil {
//		log.Fatal(err)
//	}
//	clean := s.Apply(line)
//
// Non-JSON lines are passed through unchanged.
package jsonstrip
