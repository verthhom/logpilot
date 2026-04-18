package fieldfilter_test

import (
	"testing"

	"github.com/yourorg/logpilot/internal/fieldfilter"
)

func TestNew_Valid(t *testing.T) {
	f, err := fieldfilter.New([]string{"level", "msg"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f == nil {
		t.Fatal("expected non-nil Filter")
	}
}

func TestNew_NoFields(t *testing.T) {
	_, err := fieldfilter.New([]string{})
	if err == nil {
		t.Fatal("expected error for empty fields")
	}
}

func TestNew_BlankField(t *testing.T) {
	_, err := fieldfilter.New([]string{"level", ""})
	if err == nil {
		t.Fatal("expected error for blank field name")
	}
}

func TestApply_RetainsOnlyFields(t *testing.T) {
	f, _ := fieldfilter.New([]string{"level", "msg"})
	line := `{"level":"info","msg":"hello","ts":"2024-01-01"}`
	out := f.Apply(line)
	for _, want := range []string{"level", "msg"} {
		if !contains(out, want) {
			t.Errorf("expected field %q in output %q", want, out)
		}
	}
	if contains(out, "ts") {
		t.Errorf("unexpected field 'ts' in output %q", out)
	}
}

func TestApply_NonJSONPassThrough(t *testing.T) {
	f, _ := fieldfilter.New([]string{"level"})
	line := "plain text log line"
	if got := f.Apply(line); got != line {
		t.Errorf("expected passthrough, got %q", got)
	}
}

func TestApply_MissingFieldsOmitted(t *testing.T) {
	f, _ := fieldfilter.New([]string{"level", "missing"})
	line := `{"level":"warn","msg":"hi"}`
	out := f.Apply(line)
	if !contains(out, "level") {
		t.Errorf("expected 'level' in output %q", out)
	}
	if contains(out, "msg") {
		t.Errorf("unexpected 'msg' in output %q", out)
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(sub) == 0 ||
		(func() bool {
			for i := 0; i <= len(s)-len(sub); i++ {
				if s[i:i+len(sub)] == sub {
					return true
				}
			}
			return false
		})())
}
