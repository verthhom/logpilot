package dedupe

import "errors"

// ErrInvalidWindow is returned when windowSize is less than 1.
var ErrInvalidWindow = errors.New("dedupe: window size must be at least 1")
