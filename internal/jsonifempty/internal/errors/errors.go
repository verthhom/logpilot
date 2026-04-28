// Package errors defines sentinel errors for the jsonifempty package.
package errors

import "errors"

// ErrNoRules is returned when no rules are provided.
var ErrNoRules = errors.New("jsonifempty: at least one rule is required")

// ErrBadRule is returned when a rule string does not contain '='.
var ErrBadRule = errors.New("jsonifempty: rule must be in field=value format")

// ErrBlankPart is returned when the field name or value is blank.
var ErrBlankPart = errors.New("jsonifempty: field and value must not be blank")

// ErrInvalidJSON is returned when the fallback value is not valid JSON.
var ErrInvalidJSON = errors.New("jsonifempty: fallback value must be valid JSON")
