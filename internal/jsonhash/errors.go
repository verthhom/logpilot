package jsonhash

import "errors"

var (
	errBlankDest  = errors.New("jsonhash: dest field must not be blank")
	errNoFields   = errors.New("jsonhash: at least one source field is required")
	errBlankField = errors.New("jsonhash: source field names must not be blank")
)
