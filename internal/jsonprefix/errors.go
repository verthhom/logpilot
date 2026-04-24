package jsonprefix

import "errors"

// ErrEmptyPrefix is returned by New when the supplied prefix is empty or
// contains only whitespace.
var ErrEmptyPrefix = errors.New("jsonprefix: prefix must not be empty")
