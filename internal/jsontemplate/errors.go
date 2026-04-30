package jsontemplate

import "errors"

// ErrBlankField is returned when the destination field name is empty or whitespace.
var ErrBlankField = errors.New("jsontemplate: field name must not be blank")

// ErrBlankTemplate is returned when the template source string is empty or whitespace.
var ErrBlankTemplate = errors.New("jsontemplate: template must not be blank")

// ErrBadTemplate is returned when the template source cannot be parsed.
var ErrBadTemplate = errors.New("jsontemplate: invalid template")
