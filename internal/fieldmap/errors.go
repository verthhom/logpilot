package fieldmap

import "fmt"

// ErrInvalidRule is returned when a field mapping rule cannot be parsed.
type ErrInvalidRule struct {
	Rule string
}

func (e *ErrInvalidRule) Error() string {
	return fmt.Sprintf("fieldmap: invalid rule %q: expected \"oldField=newField\"", e.Rule)
}
