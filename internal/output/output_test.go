package output

import (
	"bytes"
	"strings"
	"testing"
)

func TestNew_ValidFormats(t *testing.T) {
	for _, f := range []Format{FormatJSON, FormatPretty} {
		_, err := New(f, nil)
		if err != nil {
			t.Errorf("expected no error for format %q, got %v", f, err)
		}
	}
}

func TestNew_InvalidFormat(t *testing.T) {
	_, err := New("xml", nil)
	if err == nil {
		t.Fatal("expected error for unknown format")
	}
}

func TestWrite_JSONFormat_PassThrough(t *testing.T) {
	var buf bytes.Buffer
	w, _ := New(FormatJSON, &buf)
	line := `{"level":"info","msg":"hello"}`
	w.Write("app", line)
	if got := strings.TrimSpace(buf.String()); got != line {
		t.Errorf("expected %q, got %q", line, got)
	}
}

func TestWrite_PrettyFormat_JSONLine(t *testing.T) {
	var buf bytes.Buffer
	w, _ := New(FormatPretty, &buf)
	w.Write("svc", `{"time":"2024-01-01T00:00:00Z","level":"warn","msg":"disk full"}`)
	out := buf.String()
	if !strings.Contains(out, "WARN") {
		t.Errorf("expected WARN in output, got: %s", out)
	}
	if !strings.Contains(out, "disk full") {
		t.Errorf("expected message in output, got: %s", out)
	}
	if !strings.Contains(out, "[svc]") {
		t.Errorf("expected source in output, got: %s", out)
	}
}

func TestWrite_PrettyFormat_NonJSON(t *testing.T) {
	var buf bytes.Buffer
	w, _ := New(FormatPretty, &buf)
	w.Write("src", "plain text log line")
	out := buf.String()
	if !strings.Contains(out, "[src]") {
		t.Errorf("expected source prefix, got: %s", out)
	}
	if !strings.Contains(out, "plain text log line") {
		t.Errorf("expected raw line in output, got: %s", out)
	}
}

func TestWrite_EmptyLine(t *testing.T) {
	var buf bytes.Buffer
	w, _ := New(FormatPretty, &buf)
	w.Write("src", "   ")
	if buf.Len() != 0 {
		t.Errorf("expected no output for blank line")
	}
}
