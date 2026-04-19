package enrich

import "errors"

// Sentinel errors returned by New.
var (
	ErrNoRules = errors.New("enrich: at least one rule is required")
	ErrBadRule = errors.New("enrich: rule must be in key=value format")
	ErrBlankKey = errors.New("enrich: key must not be blank")
)
