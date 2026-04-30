package jsonspread

import (
	"encoding/json"
	"testing"
)

func TestNew_Valid(t *testing.T) {
	s, err := New([]string{"meta"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s == nil {
		t.Fatal("expected non-nil Spreader")
	}
}

func TestNew_NoFields(t *testing.T) {
	_, err := New([]string{})
	if err == nil {
		t.Fatal("expected error for empty fields slice")
	}
}

func TestNew_BlankField(t *testing.T) {
	_, err := New([]string{"ok", ""})
	if err == nil {
		t.Fatal("expected error for blank field name")
	}
}

func TestApply_NonJSONPassThrough(t *testing.T) {
	s, _ := New([]string{"meta"})
	input := "not json at all"
	if got := s.Apply(input); got != input {
		t.Fatalf("expected pass-through, got %q", got)
	}
}

func TestApply_SpreadsNestedObject(t *testing.T) {
	s, _ := New([]string{"meta"})
	input := `{"msg":"hello","meta":{"host":"srv1","env":"prod"}}`
	got := s.Apply(input)

	var doc map[string]json.RawMessage
	if err := json.Unmarshal([]byte(got), &doc); err != nil {
		t.Fatalf("result is not valid JSON: %v", err)
	}
	if _, ok := doc["meta"]; ok {
		t.Error("expected 'meta' field to be removed")
	}
	if _, ok := doc["host"]; !ok {
		t.Error("expected 'host' field to be promoted")
	}
	if _, ok := doc["env"]; !ok {
		t.Error("expected 'env' field to be promoted")
	}
}

func TestApply_DoesNotOverwriteExisting(t *testing.T) {
	s, _ := New([]string{"meta"})
	input := `{"host":"original","meta":{"host":"nested"}}`
	got := s.Apply(input)

	var doc map[string]string
	if err := json.Unmarshal([]byte(got), &doc); err != nil {
		t.Fatalf("result is not valid JSON: %v", err)
	}
	if doc["host"] != "original" {
		t.Errorf("expected 'host' to remain %q, got %q", "original", doc["host"])
	}
}

func TestApply_FieldNotObject_Unchanged(t *testing.T) {
	s, _ := New([]string{"meta"})
	input := `{"msg":"hi","meta":"scalar"}`
	got := s.Apply(input)

	var doc map[string]json.RawMessage
	if err := json.Unmarshal([]byte(got), &doc); err != nil {
		t.Fatalf("result is not valid JSON: %v", err)
	}
	if _, ok := doc["meta"]; !ok {
		t.Error("expected scalar 'meta' field to remain")
	}
}

func TestApply_MissingField_NoOp(t *testing.T) {
	s, _ := New([]string{"meta"})
	input := `{"msg":"hello"}`
	got := s.Apply(input)

	var doc map[string]json.RawMessage
	if err := json.Unmarshal([]byte(got), &doc); err != nil {
		t.Fatalf("result is not valid JSON: %v", err)
	}
	if _, ok := doc["msg"]; !ok {
		t.Error("expected 'msg' field to be present")
	}
}

func TestApply_MultipleFields(t *testing.T) {
	s, _ := New([]string{"meta", "tags"})
	input := `{"meta":{"host":"srv1"},"tags":{"region":"us-east"}}`
	got := s.Apply(input)

	var doc map[string]json.RawMessage
	if err := json.Unmarshal([]byte(got), &doc); err != nil {
		t.Fatalf("result is not valid JSON: %v", err)
	}
	for _, key := range []string{"host", "region"} {
		if _, ok := doc[key]; !ok {
			t.Errorf("expected %q to be promoted", key)
		}
	}
	for _, key := range []string{"meta", "tags"} {
		if _, ok := doc[key]; ok {
			t.Errorf("expected %q to be removed", key)
		}
	}
}
