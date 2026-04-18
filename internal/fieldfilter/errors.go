package fieldfilter

import "errors"

// ErrNoFields is returned when no fields are provided to New.
var ErrNoFields = errors.New("fieldfilter: at least one field is required")

// ErrInvalidField is returned when a blank field name is provided.
var ErrInvalidField = errors.New("fieldfilter: invalid field")
