package jsonlookup

import "errors"

// ErrBlankField is returned when the source field name is empty or whitespace.
var ErrBlankField = errors.New("jsonlookup: field name must not be blank")

// ErrEmptyTable is returned when the lookup table contains no entries.
var ErrEmptyTable = errors.New("jsonlookup: lookup table must not be empty")
