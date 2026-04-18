package throttle

import "errors"

// ErrInvalidMaxLines is returned when maxLines is not positive.
var ErrInvalidMaxLines = errors.New("throttle: maxLines must be greater than zero")

// ErrInvalidWindow is returned when the window duration is not positive.
var ErrInvalidWindow = errors.New("throttle: window must be greater than zero")
