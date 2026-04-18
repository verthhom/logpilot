package filter

import "fmt"

// InvalidRuleError is returned when a rule string cannot be parsed.
type InvalidRuleError struct {
	Raw string
}

func (e *InvalidRuleError) Error() string {
	return fmt.Sprintf("filter: invalid rule %q: expected field:operator[:value]", e.Raw)
}
