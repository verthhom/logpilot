package levelfilter

import "errors"

// ErrEmptyField is returned when the field name is blank.
var ErrEmptyField = errors.New("levelfilter: field name must not be empty")

// ErrUnknownLevel is returned when the supplied level string is not recognised.
var ErrUnknownLevel = errors.New("levelfilter: unknown log level")
