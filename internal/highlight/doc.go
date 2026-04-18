// Package highlight provides ANSI terminal color support for logpilot output.
//
// A Highlighter can be created in enabled or disabled mode. When disabled,
// all methods return plain text, making it safe to use unconditionally and
// letting callers toggle color based on whether the output is a TTY.
//
// Example usage:
//
//	h := highlight.New(isatty.IsTerminal(os.Stdout.Fd()))
//	fmt.Println(h.Level("error"))   // red + bold "error" on a TTY
//	fmt.Println(h.Key("timestamp")) // cyan "timestamp" on a TTY
//
// The Strip function removes all ANSI escape sequences from a string, which
// is useful when writing colored output to a file or comparing in tests.
package highlight
