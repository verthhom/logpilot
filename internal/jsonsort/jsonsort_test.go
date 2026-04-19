package jsonsort

import (
	"encoding/json"
	"testing"
)

func TestNew_Valid(t *testing.T) {
	_, err := New([]string{"time", "level", "msg"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNew_NoFields(t *testing.T) {
	_, err := New([]string{})
	if err != ErrNoFields {
		t.Fatalf("expected ErrNoFields, got %v", err)
	}
}

func TestNew_BlankField(t *testing.T) {
	_, err := New([]string{"time", ""})
	if err != ErrBlankField {
		t.Fatalf("expected ErrBlankField, got %v", err)
	}
}

func TestApply_NonJSONPassThrough(t *testing.T) {
	s, _ := New([]string{"time"})
	input := "not json at all"
	if got := s.Apply(input); got != input {
		t.Fatalf("expected pass-through, got %q", got)
	}
}

func TestApply_OrderedFieldsFirst(t *testing.T) {
	s, _ := New([]string{"time", "level", "msg"})
	input := `{"msg":"hello","level":"info","time":"2024-01-01","extra":"x"}`
	result := s.Apply(input)

	var m map[string]string
	if err := json.Unmarshal([]byte(result), &m); err != nil {
		t.Fatalf("result is not valid JSON: %v", err)
	}

	// Verify all fields preserved
	for _, k := range []string{"time", "level", "msg", "extra"} {
		if _, ok := m[k]; !ok {
			t.Errorf("missing field %q in result", k)
		}
	}

	// Verify ordering via raw bytes
	expectOrder := []string{"time", "level", "msg", "extra"}
	last := 0
	for _, f := range expectOrder {
		idx := indexOf(result, `"`+f+`"`)
		if idx < last {
			t.Errorf("field %q out of order in %s", f, result)
		}
		last = idx
	}
}

func TestApply_MissingOrderedField(t *testing.T) {
	s, _ := New([]string{"time", "level"})
	input := `{"msg":"hello","level":"info"}`
	result := s.Apply(input)

	var m map[string]string
	if err := json.Unmarshal([]byte(result), &m); err != nil {
		t.Fatalf("result is not valid JSON: %v", err)
	}
	if m["level"] != "info" || m["msg"] != "hello" {
		t.Fatalf("unexpected values: %v", m)
	}
}

func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
