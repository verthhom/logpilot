package jsoncoalesce

import (
	"encoding/json"
	"testing"
)

func TestNew_Valid(t *testing.T) {
	_, err := New("message", []string{"msg", "text"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNew_EmptyDest(t *testing.T) {
	_, err := New("", []string{"msg"})
	if err != ErrEmptyDest {
		t.Fatalf("expected ErrEmptyDest, got %v", err)
	}
}

func TestNew_NoSources(t *testing.T) {
	_, err := New("message", nil)
	if err != ErrNoSources {
		t.Fatalf("expected ErrNoSources, got %v", err)
	}
}

func TestNew_BlankSource(t *testing.T) {
	_, err := New("message", []string{"msg", ""})
	if err != ErrBlankSource {
		t.Fatalf("expected ErrBlankSource, got %v", err)
	}
}

func TestApply_FirstSourceWins(t *testing.T) {
	c, _ := New("message", []string{"msg", "text"})
	line := `{"msg":"hello","text":"world"}`
	out := c.Apply(line)

	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(out), &obj); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if obj["message"] != "hello" {
		t.Errorf("expected 'hello', got %v", obj["message"])
	}
}

func TestApply_FallsBackToSecondSource(t *testing.T) {
	c, _ := New("message", []string{"msg", "text"})
	line := `{"text":"fallback"}`
	out := c.Apply(line)

	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(out), &obj); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if obj["message"] != "fallback" {
		t.Errorf("expected 'fallback', got %v", obj["message"])
	}
}

func TestApply_NoSourcePresent_Unchanged(t *testing.T) {
	c, _ := New("message", []string{"msg", "text"})
	line := `{"level":"info"}`
	out := c.Apply(line)
	if out != line {
		t.Errorf("expected unchanged line, got %q", out)
	}
}

func TestApply_NonJSONPassThrough(t *testing.T) {
	c, _ := New("message", []string{"msg"})
	line := "not json at all"
	out := c.Apply(line)
	if out != line {
		t.Errorf("expected pass-through, got %q", out)
	}
}

func TestApply_EmptyStringSourceSkipped(t *testing.T) {
	c, _ := New("message", []string{"msg", "text"})
	line := `{"msg":"","text":"used"}`
	out := c.Apply(line)

	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(out), &obj); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if obj["message"] != "used" {
		t.Errorf("expected 'used', got %v", obj["message"])
	}
}
