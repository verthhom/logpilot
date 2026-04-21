package jsoncompare

import "errors"

// ErrBadRule is returned when the rule string does not contain exactly two colons.
var ErrBadRule = errors.New("jsoncompare: rule must be in the form 'left:right:dest'")

// ErrBlankPart is returned when any part of the rule is blank or whitespace.
var ErrBlankPart = errors.New("jsoncompare: rule parts must not be blank")
