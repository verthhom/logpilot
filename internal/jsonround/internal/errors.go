package internal

import "errors"

// ErrNoRules is returned when New is called with an empty rule slice.
var ErrNoRules = errors.New("jsonround: at least one rule is required")

// ErrBadRule is returned when a rule string cannot be parsed.
var ErrBadRule = errors.New("jsonround: invalid rule")
