package jsonpick_test

import (
	"encoding/json"
	"testing"

	"github.com/yourorg/logpilot/internal/jsonpick"
)

func TestNew_Valid(t *testing.T) {
	p, err := jsonpick.New(3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p == nil {
		t.Fatal("expected non-nil Picker")
	}
}

func TestNew_ZeroLimit(t *testing.T) {
	_, err := jsonpick.New(0)
	if err == nil {
		t.Fatal("expected error for limit=0")
	}
}

func TestNew_NegativeLimit(t *testing.T) {
	_, err := jsonpick.New(-5)
	if err == nil {
		t.Fatal("expected error for negative limit")
	}
}

func TestApply_NonJSONPassThrough(t *testing.T) {
	p, _ := jsonpick.New(2)
	input := "not json at all"
	got := p.Apply(input)
	if got != input {
		t.Errorf("expected passthrough, got %q", got)
	}
}

func TestApply_LimitExceedsFields(t *testing.T) {
	p, _ := jsonpick.New(10)
	input := `{"a":1,"b":2,"c":3}`
	got := p.Apply(input)

	var m map[string]interface{}
	if err := json.Unmarshal([]byte(got), &m); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	if len(m) != 3 {
		t.Errorf("expected 3 fields, got %d", len(m))
	}
}

func TestApply_LimitReducesFields(t *testing.T) {
	p, _ := jsonpick.New(2)
	input := `{"x":1,"y":2,"z":3,"w":4}`
	got := p.Apply(input)

	var m map[string]interface{}
	if err := json.Unmarshal([]byte(got), &m); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	if len(m) != 2 {
		t.Errorf("expected 2 fields, got %d", len(m))
	}
}

func TestApply_LimitOne(t *testing.T) {
	p, _ := jsonpick.New(1)
	input := `{"level":"info","msg":"hello","ts":"2024-01-01"}`
	got := p.Apply(input)

	var m map[string]interface{}
	if err := json.Unmarshal([]byte(got), &m); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	if len(m) != 1 {
		t.Errorf("expected 1 field, got %d", len(m))
	}
}

func TestApply_EmptyObject(t *testing.T) {
	p, _ := jsonpick.New(3)
	input := `{}`
	got := p.Apply(input)

	var m map[string]interface{}
	if err := json.Unmarshal([]byte(got), &m); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	if len(m) != 0 {
		t.Errorf("expected 0 fields, got %d", len(m))
	}
}
