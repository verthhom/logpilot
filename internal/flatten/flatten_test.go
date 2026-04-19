package flatten_test

import (
	"encoding/json"
	"testing"

	"github.com/logpilot/internal/flatten"
)

func TestNew_Valid(t *testing.T) {
	_, err := flatten.New(".")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestNew_EmptySeparator(t *testing.T) {
	_, err := flatten.New("")
	if err == nil {
		t.Fatal("expected error for empty separator")
	}
}

func TestApply_NonJSON(t *testing.T) {
	f, _ := flatten.New(".")
	input := "plain text line"
	if got := f.Apply(input); got != input {
		t.Errorf("expected passthrough, got %q", got)
	}
}

func TestApply_AlreadyFlat(t *testing.T) {
	f, _ := flatten.New(".")
	input := `{"level":"info","msg":"hello"}`
	out := f.Apply(input)
	var m map[string]any
	if err := json.Unmarshal([]byte(out), &m); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if m["level"] != "info" || m["msg"] != "hello" {
		t.Errorf("unexpected map: %v", m)
	}
}

func TestApply_NestedObject(t *testing.T) {
	f, _ := flatten.New(".")
	input := `{"a":{"b":{"c":42}}}`
	out := f.Apply(input)
	var m map[string]any
	if err := json.Unmarshal([]byte(out), &m); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	val, ok := m["a.b.c"]
	if !ok {
		t.Fatalf("expected key a.b.c, got %v", m)
	}
	if val.(float64) != 42 {
		t.Errorf("expected 42, got %v", val)
	}
}

func TestApply_CustomSeparator(t *testing.T) {
	f, _ := flatten.New("_")
	input := `{"x":{"y":"z"}}`
	out := f.Apply(input)
	var m map[string]any
	if err := json.Unmarshal([]byte(out), &m); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if _, ok := m["x_y"]; !ok {
		t.Errorf("expected key x_y, got %v", m)
	}
}

func TestApply_MixedValues(t *testing.T) {
	f, _ := flatten.New(".")
	input := `{"top":"val","nested":{"n":1}}`
	out := f.Apply(input)
	var m map[string]any
	if err := json.Unmarshal([]byte(out), &m); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if m["top"] != "val" {
		t.Errorf("expected top=val, got %v", m["top"])
	}
	if _, ok := m["nested.n"]; !ok {
		t.Errorf("expected nested.n, got %v", m)
	}
}
