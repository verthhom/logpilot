package replay

import "errors"

// ErrEmptyPath is returned when an empty file path is provided.
var ErrEmptyPath = errors.New("replay: file path must not be empty")

// ErrNegativeDelay is returned when a negative delay duration is provided.
var ErrNegativeDelay = errors.New("replay: delay must not be negative")
