package multiline

import "errors"

// ErrEmptyPattern is returned when an empty start pattern is provided.
var ErrEmptyPattern = errors.New("multiline: start pattern must not be empty")

// ErrBadPattern is returned when the start pattern cannot be compiled.
var ErrBadPattern = errors.New("multiline: start pattern is not a valid regular expression")
