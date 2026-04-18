package multiline

import (
	"regexp"
	"strings"
)

// Joiner accumulates lines that belong to a single logical log event.
// A new event begins whenever a line matches the configured start pattern.
type Joiner struct {
	start   *regexp.Regexp
	buf     []string
	sep     string
}

// New returns a Joiner that treats lines matching pattern as the start of a
// new event. sep is placed between joined lines (usually " " or "\n").
func New(pattern, sep string) (*Joiner, error) {
	if pattern == "" {
		return nil, ErrEmptyPattern
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, ErrBadPattern
	}
	if sep == "" {
		sep = " "
	}
	return &Joiner{start: re, sep: sep}, nil
}

// Feed adds a line to the joiner. If the line starts a new event and a
// previous event was buffered, the previous event is returned together with
// true. Otherwise ("" , false) is returned.
func (j *Joiner) Feed(line string) (string, bool) {
	if j.start.MatchString(line) {
		if len(j.buf) == 0 {
			j.buf = append(j.buf, line)
			return "", false
		}
		out := strings.Join(j.buf, j.sep)
		j.buf = []string{line}
		return out, true
	}
	j.buf = append(j.buf, line)
	return "", false
}

// Flush returns any remaining buffered event and resets state.
func (j *Joiner) Flush() (string, bool) {
	if len(j.buf) == 0 {
		return "", false
	}
	out := strings.Join(j.buf, j.sep)
	j.buf = nil
	return out, true
}
