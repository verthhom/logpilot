package maskfield

import "errors"

// Sentinel errors for maskfield construction.
var (
	ErrNoFields        = errors.New("maskfield: at least one field required")
	ErrBlankField      = errors.New("maskfield: field name must not be blank")
	ErrNegativeVisible = errors.New("maskfield: visible prefix length must be non-negative")
)
