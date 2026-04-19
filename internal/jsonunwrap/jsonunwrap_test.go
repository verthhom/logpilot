package jsonunwrap

import (
	"encoding/json"
	"testing"
)

func TestNew_Valid(t *testing.T) {
	_, err := New("meta")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNew_EmptyField(t *testing.T) {
	_, err := New("")
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestApply_PromotesNestedFields(t *testing.T) {
	u, _ := New("meta")
	input := `{"level":"info","meta":{"host":"srv1","env":"prod"}}`
	out := u.Apply(input)

	var m map[string]interface{}
	if err := json.Unmarshal([]byte(out), &m); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if _, ok := m["meta"]; ok {
		t.Error("meta key should have been removed")
	}
	if m["host"] != "srv1" {
		t.Errorf("expected host=srv1, got %v", m["host"])
	}
	if m["env"] != "prod" {
		t.Errorf("expected env=prod, got %v", m["env"])
	}
	if m["level"] != "info" {
		t.Errorf("expected level=info, got %v", m["level"])
	}
}

func TestApply_DoesNotOverwriteExisting(t *testing.T) {
	u, _ := New("meta")
	input := `{"host":"original","meta":{"host":"nested"}}`
	out := u.Apply(input)

	var m map[string]interface{}
	if err := json.Unmarshal([]byte(out), &m); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if m["host"] != "original" {
		t.Errorf("expected host=original, got %v", m["host"])
	}
}

func TestApply_MissingField_PassThrough(t *testing.T) {
	u, _ := New("meta")
	input := `{"level":"warn"}`
	if got := u.Apply(input); got != input {
		t.Errorf("expected passthrough, got %q", got)
	}
}

func TestApply_NonJSONPassThrough(t *testing.T) {
	u, _ := New("meta")
	input := "not json at all"
	if got := u.Apply(input); got != input {
		t.Errorf("expected passthrough, got %q", got)
	}
}

func TestApply_NestedNotObject_PassThrough(t *testing.T) {
	u, _ := New("meta")
	input := `{"meta":"string-not-object"}`
	if got := u.Apply(input); got != input {
		t.Errorf("expected passthrough when nested value is not object, got %q", got)
	}
}
