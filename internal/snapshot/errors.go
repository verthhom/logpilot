package snapshot

import "errors"

// ErrEmptyPath is returned when an empty file path is provided to New.
var ErrEmptyPath = errors.New("snapshot: path must not be empty")
