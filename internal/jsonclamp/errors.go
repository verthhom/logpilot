package jsonclamp

import "errors"

// ErrNoRules is returned when no clamp rules are provided.
var ErrNoRules = errors.New("jsonclamp: at least one rule is required")

// ErrBadRule is returned when a rule string cannot be parsed.
var ErrBadRule = errors.New("jsonclamp: invalid rule format, expected field:min:max")

// ErrBlankField is returned when the field name in a rule is empty.
var ErrBlankField = errors.New("jsonclamp: field name must not be blank")

// ErrInvalidRange is returned when min is greater than max.
var ErrInvalidRange = errors.New("jsonclamp: min must not be greater than max")
