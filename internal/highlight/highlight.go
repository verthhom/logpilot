package highlight

import (
	"fmt"
	"strings"
)

// ANSI color codes.
const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Cyan   = "\033[36m"
	Bold   = "\033[1m"
)

// LevelColors maps log level strings to ANSI color codes.
var LevelColors = map[string]string{
	"error": Red,
	"fatal": Red,
	"warn":  Yellow,
	"info":  Green,
	"debug": Cyan,
}

// Highlighter applies ANSI colors to log output.
type Highlighter struct {
	enabled bool
}

// New returns a new Highlighter. If enabled is false, all methods return
// plain text without any escape sequences.
func New(enabled bool) *Highlighter {
	return &Highlighter{enabled: enabled}
}

// Level wraps the given level string with its associated color, if known.
func (h *Highlighter) Level(level string) string {
	if !h.enabled {
		return level
	}
	color, ok := LevelColors[strings.ToLower(level)]
	if !ok {
		return level
	}
	return fmt.Sprintf("%s%s%s%s", Bold, color, level, Reset)
}

// Key wraps a JSON key with cyan color.
func (h *Highlighter) Key(key string) string {
	if !h.enabled {
		return key
	}
	return fmt.Sprintf("%s%s%s", Cyan, key, Reset)
}

// Value wraps a value string with bold formatting.
func (h *Highlighter) Value(value string) string {
	if !h.enabled {
		return value
	}
	return fmt.Sprintf("%s%s%s", Bold, value, Reset)
}

// Strip removes all ANSI escape sequences from s.
func Strip(s string) string {
	var b strings.Builder
	inEscape := false
	for _, r := range s {
		switch {
		case r == '\033':
			inEscape = true
		case inEscape && r == 'm':
			inEscape = false
		case !inEscape:
			b.WriteRune(r)
		}
	}
	return b.String()
}
