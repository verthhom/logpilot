// Package errors defines sentinel errors for the jsonslice package.
package errors

import "errors"

// ErrBlankField is returned when the field name is empty.
var ErrBlankField = errors.New("jsonslice: field name must not be blank")

// ErrNegativeIndex is returned when start is negative.
var ErrNegativeIndex = errors.New("jsonslice: start index must be >= 0")

// ErrEndBeforeStart is returned when end is less than start (and not -1).
var ErrEndBeforeStart = errors.New("jsonslice: end index must be >= start or -1 for open-ended")
