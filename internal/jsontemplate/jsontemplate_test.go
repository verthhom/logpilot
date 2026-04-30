package jsontemplate

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"
)

func TestNew_Valid(t *testing.T) {
	a, err := New("summary", "{{.level}}: {{.msg}}")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a == nil {
		t.Fatal("expected non-nil Applier")
	}
}

func TestNew_BlankField(t *testing.T) {
	_, err := New("  ", "{{.msg}}")
	if !errors.Is(err, ErrBlankField) {
		t.Fatalf("expected ErrBlankField, got %v", err)
	}
}

func TestNew_BlankTemplate(t *testing.T) {
	_, err := New("summary", "")
	if !errors.Is(err, ErrBlankTemplate) {
		t.Fatalf("expected ErrBlankTemplate, got %v", err)
	}
}

func TestNew_BadTemplate(t *testing.T) {
	_, err := New("summary", "{{.level")
	if !errors.Is(err, ErrBadTemplate) {
		t.Fatalf("expected ErrBadTemplate, got %v", err)
	}
}

func TestApply_InjectsField(t *testing.T) {
	a, _ := New("summary", "{{.level}}: {{.msg}}")
	out := a.Apply(`{"level":"error","msg":"disk full"}`)

	var obj map[string]any
	if err := json.Unmarshal([]byte(out), &obj); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	got, ok := obj["summary"].(string)
	if !ok {
		t.Fatalf("summary field missing or wrong type")
	}
	if got != "error: disk full" {
		t.Errorf("expected 'error: disk full', got %q", got)
	}
}

func TestApply_NonJSONPassThrough(t *testing.T) {
	a, _ := New("summary", "{{.msg}}")
	line := "not json at all"
	if out := a.Apply(line); out != line {
		t.Errorf("expected pass-through, got %q", out)
	}
}

func TestApply_MissingKeyZero(t *testing.T) {
	a, _ := New("label", "svc={{.service}}")
	out := a.Apply(`{"level":"info"}`)

	var obj map[string]any
	if err := json.Unmarshal([]byte(out), &obj); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	got, _ := obj["label"].(string)
	if !strings.HasPrefix(got, "svc=") {
		t.Errorf("unexpected label value: %q", got)
	}
}

func TestApply_OverwritesExistingField(t *testing.T) {
	a, _ := New("msg", "PREFIXED: {{.msg}}")
	out := a.Apply(`{"msg":"hello"}`)

	var obj map[string]any
	if err := json.Unmarshal([]byte(out), &obj); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if obj["msg"] != "PREFIXED: hello" {
		t.Errorf("expected overwritten msg, got %v", obj["msg"])
	}
}
