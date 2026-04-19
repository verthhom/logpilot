package jsonmerge

import (
	"encoding/json"
	"testing"
)

func TestNew_Valid(t *testing.T) {
	_, err := New(map[string]string{"env": "prod"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNew_NoFields(t *testing.T) {
	_, err := New(map[string]string{})
	if err == nil {
		t.Fatal("expected error for empty fields")
	}
}

func TestApply_InjectsField(t *testing.T) {
	m, _ := New(map[string]string{"env": "prod"})
	out := m.Apply(`{"msg":"hello"}`)

	var obj map[string]string
	if err := json.Unmarshal([]byte(out), &obj); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if obj["env"] != "prod" {
		t.Errorf("expected env=prod, got %q", obj["env"])
	}
	if obj["msg"] != "hello" {
		t.Errorf("original field lost")
	}
}

func TestApply_DoesNotOverwrite(t *testing.T) {
	m, _ := New(map[string]string{"env": "prod"})
	out := m.Apply(`{"env":"staging"}`)

	var obj map[string]string
	if err := json.Unmarshal([]byte(out), &obj); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if obj["env"] != "staging" {
		t.Errorf("existing field should not be overwritten, got %q", obj["env"])
	}
}

func TestApply_NonJSONPassThrough(t *testing.T) {
	m, _ := New(map[string]string{"env": "prod"})
	plain := "not json at all"
	if got := m.Apply(plain); got != plain {
		t.Errorf("expected pass-through, got %q", got)
	}
}

func TestApply_MultipleFields(t *testing.T) {
	m, _ := New(map[string]string{"env": "prod", "region": "us-east-1"})
	out := m.Apply(`{"msg":"ok"}`)

	var obj map[string]string
	if err := json.Unmarshal([]byte(out), &obj); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if obj["env"] != "prod" {
		t.Errorf("missing env field")
	}
	if obj["region"] != "us-east-1" {
		t.Errorf("missing region field")
	}
}
