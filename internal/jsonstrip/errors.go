package jsonstrip

import "errors"

// ErrNoFields is returned when New is called with an empty field list.
var ErrNoFields = errors.New("jsonstrip: no fields provided")
