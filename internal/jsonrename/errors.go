package jsonrename

import (
	"errors"
	"fmt"
)

// ErrNoRules is returned when no rename rules are provided.
var ErrNoRules = errors.New("jsonrename: at least one rule is required")

// ErrInvalidRule is returned when a rule does not follow the "from:to" format.
type ErrInvalidRule string

func (e ErrInvalidRule) Error() string {
	return fmt.Sprintf("jsonrename: invalid rule %q: expected \"from:to\" format", string(e))
}
