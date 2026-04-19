package jsonrename

import (
	"encoding/json"
	"testing"
)

func TestNew_Valid(t *testing.T) {
	_, err := New([]string{"msg:message"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNew_NoRules(t *testing.T) {
	_, err := New([]string{})
	if err == nil {
		t.Fatal("expected error for empty rules")
	}
}

func TestNew_BadRule(t *testing.T) {
	_, err := New([]string{"nodivider"})
	if err == nil {
		t.Fatal("expected error for bad rule")
	}
}

func TestNew_BlankPart(t *testing.T) {
	_, err := New([]string{":to"})
	if err == nil {
		t.Fatal("expected error for blank from-part")
	}
}

func TestApply_RenamesField(t *testing.T) {
	r, _ := New([]string{"msg:message"})
	out := r.Apply(`{"msg":"hello","level":"info"}`)
	var obj map[string]json.RawMessage
	if err := json.Unmarshal([]byte(out), &obj); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if _, ok := obj["message"]; !ok {
		t.Error("expected field 'message'")
	}
	if _, ok := obj["msg"]; ok {
		t.Error("old field 'msg' should be removed")
	}
}

func TestApply_NonJSONPassThrough(t *testing.T) {
	r, _ := New([]string{"a:b"})
	line := "not json at all"
	if got := r.Apply(line); got != line {
		t.Errorf("expected passthrough, got %q", got)
	}
}

func TestApply_MissingFieldNoOp(t *testing.T) {
	r, _ := New([]string{"missing:other"})
	line := `{"level":"info"}`
	out := r.Apply(line)
	var obj map[string]json.RawMessage
	if err := json.Unmarshal([]byte(out), &obj); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if _, ok := obj["other"]; ok {
		t.Error("field 'other' should not exist")
	}
}

func TestApply_MultipleRules(t *testing.T) {
	r, _ := New([]string{"msg:message", "ts:timestamp"})
	out := r.Apply(`{"msg":"hi","ts":"2024-01-01T00:00:00Z"}`)
	var obj map[string]json.RawMessage
	if err := json.Unmarshal([]byte(out), &obj); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	for _, want := range []string{"message", "timestamp"} {
		if _, ok := obj[want]; !ok {
			t.Errorf("expected field %q", want)
		}
	}
}
