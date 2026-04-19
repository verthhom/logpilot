package jsonstrip

import (
	"encoding/json"
	"testing"
)

func TestNew_Valid(t *testing.T) {
	_, err := New([]string{"secret"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNew_NoFields(t *testing.T) {
	_, err := New([]string{})
	if err == nil {
		t.Fatal("expected error for empty fields")
	}
}

func TestNew_BlankField(t *testing.T) {
	_, err := New([]string{"ok", ""})
	if err == nil {
		t.Fatal("expected error for blank field")
	}
}

func TestApply_RemovesFields(t *testing.T) {
	s, _ := New([]string{"token", "password"})
	input := `{"user":"alice","token":"abc","password":"secret"}`
	out := s.Apply(input)

	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(out), &obj); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	if _, ok := obj["token"]; ok {
		t.Error("expected token to be removed")
	}
	if _, ok := obj["password"]; ok {
		t.Error("expected password to be removed")
	}
	if obj["user"] != "alice" {
		t.Errorf("expected user=alice, got %v", obj["user"])
	}
}

func TestApply_NonJSONPassThrough(t *testing.T) {
	s, _ := New([]string{"token"})
	input := "plain text line"
	if got := s.Apply(input); got != input {
		t.Errorf("expected pass-through, got %q", got)
	}
}

func TestApply_MissingFieldNoError(t *testing.T) {
	s, _ := New([]string{"nonexistent"})
	input := `{"level":"info","msg":"hello"}`
	out := s.Apply(input)
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(out), &obj); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(obj) != 2 {
		t.Errorf("expected 2 fields, got %d", len(obj))
	}
}

func TestApply_MultipleFields(t *testing.T) {
	s, _ := New([]string{"a", "b", "c"})
	input := `{"a":1,"b":2,"c":3,"d":4}`
	out := s.Apply(input)
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(out), &obj); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(obj) != 1 {
		t.Errorf("expected 1 field remaining, got %d", len(obj))
	}
}
