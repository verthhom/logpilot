package jsongroup

import "errors"

// ErrBlankField is returned when the grouping field name is empty or whitespace.
var ErrBlankField = errors.New("jsongroup: field name must not be blank")
