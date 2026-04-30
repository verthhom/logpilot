package jsonxform

import "errors"

// ErrEmptyRule is returned when an empty rule string is supplied.
var ErrEmptyRule = errors.New("jsonxform: rule must not be empty")

// ErrBadRule is returned when the rule string cannot be parsed.
var ErrBadRule = errors.New("jsonxform: invalid rule format, expected dest=src:template")

// ErrBadTemplate is returned when the template expression fails to compile.
var ErrBadTemplate = errors.New("jsonxform: invalid template expression")
