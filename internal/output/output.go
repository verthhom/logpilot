package output

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"
	"time"

	"github.com/user/logpilot/internal/highlight"
)

// Format represents the output rendering format.
type Format string

const (
	FormatJSON   Format = "json"
	FormatPretty Format = "pretty"
)

var validFormats = map[Format]bool{
	FormatJSON:   true,
	FormatPretty: true,
}

// Writer formats and writes log lines to an io.Writer.
type Writer struct {
	out    io.Writer
	format Format
	hl     *highlight.Highlighter
}

// New creates a Writer for the given format. color controls ANSI highlighting
// (only applied in pretty mode). Returns an error for unknown formats.
func New(out io.Writer, format Format, color bool) (*Writer, error) {
	if !validFormats[format] {
		return nil, fmt.Errorf("unknown format %q: must be json or pretty", format)
	}
	return &Writer{out: out, format: format, hl: highlight.New(color)}, nil
}

// Write formats line and writes it to the underlying writer.
func (w *Writer) Write(line string) error {
	switch w.format {
	case FormatJSON:
		_, err := fmt.Fprintln(w.out, line)
		return err
	case FormatPretty:
		return w.writePretty(line)
	}
	return nil
}

func (w *Writer) writePretty(line string) error {
	var fields map[string]interface{}
	if err := json.Unmarshal([]byte(line), &fields); err != nil {
		_, err = fmt.Fprintln(w.out, line)
		return err
	}

	ts := extractTime(fields)
	level := extractString(fields, "level")
	msg := extractString(fields, "message", "msg")

	var sb strings.Builder
	if ts != "" {
		sb.WriteString(ts)
		sb.WriteString(" ")
	}
	if level != "" {
		sb.WriteString(w.hl.Level(level))
		sb.WriteString(" ")
	}
	if msg != "" {
		sb.WriteString(w.hl.Value(msg))
	}

	skip := map[string]bool{"time": true, "timestamp": true, "ts": true,
		"level": true, "message": true, "msg": true}
	keys := make([]string, 0, len(fields))
	for k := range fields {
		if !skip[k] {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)
	for _, k := range keys {
		sb.WriteString(" ")
		sb.WriteString(w.hl.Key(k))
		sb.WriteString("=")
		sb.WriteString(fmt.Sprintf("%v", fields[k]))
	}

	_, err := fmt.Fprintln(w.out, sb.String())
	return err
}

func extractTime(fields map[string]interface{}) string {
	for _, key := range []string{"time", "timestamp", "ts"} {
		if v, ok := fields[key]; ok {
			s := fmt.Sprintf("%v", v)
			if t, err := time.Parse(time.RFC3339, s); err == nil {
				return t.Format("15:04:05")
			}
			return s
		}
	}
	return ""
}

func extractString(fields map[string]interface{}, keys ...string) string {
	for _, k := range keys {
		if v, ok := fields[k]; ok {
			return fmt.Sprintf("%v", v)
		}
	}
	return ""
}
