package tagfield

import (
	"encoding/json"
	"testing"
)

func TestNew_Valid(t *testing.T) {
	_, err := New("env", "production")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNew_BlankField(t *testing.T) {
	_, err := New("", "production")
	if err != ErrBlankField {
		t.Fatalf("expected ErrBlankField, got %v", err)
	}
}

func TestNew_BlankValue(t *testing.T) {
	_, err := New("env", "")
	if err != ErrBlankValue {
		t.Fatalf("expected ErrBlankValue, got %v", err)
	}
}

func TestApply_InjectsTag(t *testing.T) {
	tg, _ := New("env", "staging")
	out := tg.Apply(`{"msg":"hello"}`)
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(out), &obj); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if obj["env"] != "staging" {
		t.Errorf("expected env=staging, got %v", obj["env"])
	}
	if obj["msg"] != "hello" {
		t.Errorf("expected msg=hello, got %v", obj["msg"])
	}
}

func TestApply_NonJSONPassThrough(t *testing.T) {
	tg, _ := New("env", "staging")
	line := "plain text line"
	if got := tg.Apply(line); got != line {
		t.Errorf("expected passthrough, got %q", got)
	}
}

func TestApply_OverwritesExistingField(t *testing.T) {
	tg, _ := New("env", "production")
	out := tg.Apply(`{"env":"dev","msg":"hi"}`)
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(out), &obj); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if obj["env"] != "production" {
		t.Errorf("expected env=production, got %v", obj["env"])
	}
}
