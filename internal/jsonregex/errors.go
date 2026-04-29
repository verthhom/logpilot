package jsonregex

import "fmt"

// ErrNoRules is returned when no rules are provided.
const ErrNoRules = jsonregexError("jsonregex: at least one rule is required")

type jsonregexError string

func (e jsonregexError) Error() string { return string(e) }

// ErrBadRule is returned when a rule string cannot be parsed.
type ErrBadRule string

func (e ErrBadRule) Error() string {
	return fmt.Sprintf("jsonregex: malformed rule %q: expected field=pattern:replacement", string(e))
}

// ErrBadPattern is returned when a rule's regular expression fails to compile.
type ErrBadPattern struct {
	Rule string
	Err  error
}

func (e ErrBadPattern) Error() string {
	return fmt.Sprintf("jsonregex: invalid pattern in rule %q: %v", e.Rule, e.Err)
}
