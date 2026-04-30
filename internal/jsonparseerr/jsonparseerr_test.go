package jsonparseerr

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestNew_Valid(t *testing.T) {
	_, err := New("parse_ok")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNew_BlankField(t *testing.T) {
	for _, f := range []string{"", "   "} {
		_, err := New(f)
		if err == nil {
			t.Fatalf("expected error for blank field %q", f)
		}
	}
}

func TestApply_ValidJSON_InjectsTrue(t *testing.T) {
	tr, _ := New("valid")
	out := tr.Apply(`{"level":"info","msg":"hello"}`)

	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(out), &obj); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	v, ok := obj["valid"]
	if !ok {
		t.Fatal("field 'valid' missing from output")
	}
	if v != true {
		t.Fatalf("expected true, got %v", v)
	}
}

func TestApply_NonJSON_AppendsSuffix(t *testing.T) {
	tr, _ := New("parse_ok")
	raw := "this is not json at all"
	out := tr.Apply(raw)

	if !strings.HasPrefix(out, raw) {
		t.Fatalf("expected output to start with original line, got %q", out)
	}
	if !strings.Contains(out, "parse_ok=false") {
		t.Fatalf("expected suffix 'parse_ok=false' in output, got %q", out)
	}
}

func TestApply_EmptyObject_InjectsTrue(t *testing.T) {
	tr, _ := New("ok")
	out := tr.Apply(`{}`)

	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(out), &obj); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if obj["ok"] != true {
		t.Fatalf("expected true, got %v", obj["ok"])
	}
}

func TestApply_PartialJSON_AppendsSuffix(t *testing.T) {
	tr, _ := New("valid")
	out := tr.Apply(`{"key":"val"`) // missing closing brace

	if !strings.Contains(out, "valid=false") {
		t.Fatalf("expected suffix for partial JSON, got %q", out)
	}
}

func TestApply_PreservesExistingFields(t *testing.T) {
	tr, _ := New("parsed")
	out := tr.Apply(`{"service":"api","level":"warn"}`)

	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(out), &obj); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if obj["service"] != "api" {
		t.Errorf("service field lost, got %v", obj["service"])
	}
	if obj["level"] != "warn" {
		t.Errorf("level field lost, got %v", obj["level"])
	}
}
