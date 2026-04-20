package jsoncoalesce

import "errors"

// Sentinel errors returned by New.
var (
	ErrEmptyDest   = errors.New("jsoncoalesce: dest field must not be empty")
	ErrNoSources   = errors.New("jsoncoalesce: at least one source field is required")
	ErrBlankSource = errors.New("jsoncoalesce: source field names must not be blank")
)
