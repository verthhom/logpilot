package tagfield

import "errors"

// ErrBlankField is returned when the tag field name is empty.
var ErrBlankField = errors.New("tagfield: field name must not be blank")

// ErrBlankValue is returned when the tag value is empty.
var ErrBlankValue = errors.New("tagfield: tag value must not be blank")
