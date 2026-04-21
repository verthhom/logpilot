package jsoncondition

import "errors"

// ErrBadRule is returned when the rule string does not follow the
// expected format "src_field=src_value:dest_field=dest_value".
var ErrBadRule = errors.New("jsoncondition: rule must be in the form 'src_field=src_value:dest_field=dest_value'")

// ErrBlankPart is returned when any component of the rule is blank.
var ErrBlankPart = errors.New("jsoncondition: rule parts must not be blank")
