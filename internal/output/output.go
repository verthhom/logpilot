// Package output handles formatting and writing of log entries to stdout.
package output

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

// Format controls how log lines are rendered.
type Format string

const (
	FormatJSON   Format = "json"
	FormatPretty Format = "pretty"
)

// Writer formats and writes log entries.
type Writer struct {
	format Format
	out    io.Writer
}

// New returns a Writer using the given format. If out is nil, os.Stdout is used.
func New(format Format, out io.Writer) (*Writer, error) {
	if format != FormatJSON && format != FormatPretty {
		return nil, fmt.Errorf("unknown output format %q: must be \"json\" or \"pretty\"", format)
	}
	if out == nil {
		out = os.Stdout
	}
	return &Writer{format: format, out: out}, nil
}

// Write formats a raw JSON log line and writes it to the output.
// Non-JSON lines are passed through as-is.
func (w *Writer) Write(source, line string) {
	line = strings.TrimSpace(line)
	if line == "" {
		return
	}
	if w.format == FormatJSON {
		fmt.Fprintln(w.out, line)
		return
	}
	// Pretty format
	var entry map[string]any
	if err := json.Unmarshal([]byte(line), &entry); err != nil {
		// Not JSON — print raw with source prefix.
		fmt.Fprintf(w.out, "[%s] %s\n", source, line)
		return
	}
	ts := extractTime(entry)
	level := extractString(entry, "level", "LOG")
	msg := extractString(entry, "msg", "message", "text")
	fmt.Fprintf(w.out, "%s [%s] [%s] %s\n", ts, strings.ToUpper(level), source, msg)
}

func extractTime(entry map[string]any) string {
	for _, key := range []string{"time", "ts", "timestamp"} {
		if v, ok := entry[key]; ok {
			if s, ok := v.(string); ok {
				return s
			}
		}
	}
	return time.Now().Format(time.RFC3339)
}

func extractString(entry map[string]any, keys ...string) string {
	for _, key := range keys {
		if v, ok := entry[key]; ok {
			if s, ok := v.(string); ok {
				return s
			}
		}
	}
	return ""
}
