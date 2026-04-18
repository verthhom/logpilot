package aggregate

import "errors"

// ErrEmptyField is returned when the grouping field name is blank.
var ErrEmptyField = errors.New("aggregate: field must not be empty")

// ErrInvalidWindow is returned when the window duration is not positive.
var ErrInvalidWindow = errors.New("aggregate: window must be greater than zero")
