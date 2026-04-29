package jsonregex

import (
	"encoding/json"
	"testing"
)

func TestNew_Valid(t *testing.T) {
	_, err := New([]string{"msg=error:\\d+:ERROR"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNew_NoRules(t *testing.T) {
	_, err := New(nil)
	if err == nil {
		t.Fatal("expected error for empty rules")
	}
}

func TestNew_BadRule(t *testing.T) {
	for _, r := range []string{"nodot", "=pattern:sub", "field="} {
		_, err := New([]string{r})
		if err == nil {
			t.Fatalf("expected error for rule %q", r)
		}
	}
}

func TestNew_BadPattern(t *testing.T) {
	_, err := New([]string{"field=[invalid:sub"})
	if err == nil {
		t.Fatal("expected error for invalid regex")
	}
}

func TestApply_ReplacesMatch(t *testing.T) {
	rep, _ := New([]string{"msg=\\d+:NUM"})
	out := rep.Apply(`{"msg":"got 42 items"}`)
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(out), &obj); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if obj["msg"] != "got NUM items" {
		t.Errorf("expected 'got NUM items', got %q", obj["msg"])
	}
}

func TestApply_NonJSONPassThrough(t *testing.T) {
	rep, _ := New([]string{"msg=x:y"})
	line := "plain text line"
	if got := rep.Apply(line); got != line {
		t.Errorf("expected passthrough, got %q", got)
	}
}

func TestApply_MissingField_NoChange(t *testing.T) {
	rep, _ := New([]string{"other=\\d+:NUM"})
	line := `{"msg":"hello 99"}`
	out := rep.Apply(line)
	var obj map[string]interface{}
	json.Unmarshal([]byte(out), &obj)
	if obj["msg"] != "hello 99" {
		t.Errorf("unexpected change: %q", obj["msg"])
	}
}

func TestApply_NonStringField_Skipped(t *testing.T) {
	rep, _ := New([]string{"count=\\d+:NUM"})
	line := `{"count":42}`
	out := rep.Apply(line)
	var obj map[string]interface{}
	json.Unmarshal([]byte(out), &obj)
	// numeric field must remain numeric
	if _, ok := obj["count"].(float64); !ok {
		t.Errorf("expected numeric count, got %T", obj["count"])
	}
}

func TestApply_MultipleRules(t *testing.T) {
	rep, _ := New([]string{"a=foo:bar", "b=baz:qux"})
	out := rep.Apply(`{"a":"foo world","b":"baz world"}`)
	var obj map[string]interface{}
	json.Unmarshal([]byte(out), &obj)
	if obj["a"] != "bar world" {
		t.Errorf("a: expected 'bar world', got %q", obj["a"])
	}
	if obj["b"] != "qux world" {
		t.Errorf("b: expected 'qux world', got %q", obj["b"])
	}
}
