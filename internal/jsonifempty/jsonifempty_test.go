package jsonifempty

import (
	"testing"

	"github.com/logpilot/logpilot/internal/jsonifempty/internal/errors"
)

func TestNew_Valid(t *testing.T) {
	f, err := New([]string{`level="unknown"`, `code=0`})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f == nil {
		t.Fatal("expected non-nil Filler")
	}
}

func TestNew_NoRules(t *testing.T) {
	_, err := New(nil)
	if err != errors.ErrNoRules {
		t.Fatalf("expected ErrNoRules, got %v", err)
	}
}

func TestNew_BadRule(t *testing.T) {
	_, err := New([]string{"noequalssign"})
	if err != errors.ErrBadRule {
		t.Fatalf("expected ErrBadRule, got %v", err)
	}
}

func TestNew_BlankPart(t *testing.T) {
	_, err := New([]string{"=\"val\""})
	if err != errors.ErrBlankPart {
		t.Fatalf("expected ErrBlankPart, got %v", err)
	}
}

func TestNew_InvalidJSONValue(t *testing.T) {
	_, err := New([]string{"field=notjson"})
	if err != errors.ErrInvalidJSON {
		t.Fatalf("expected ErrInvalidJSON, got %v", err)
	}
}

func TestApply_NonJSONPassThrough(t *testing.T) {
	f, _ := New([]string{`level="unknown"`})
	got := f.Apply("plain text line")
	if got != "plain text line" {
		t.Fatalf("expected pass-through, got %q", got)
	}
}

func TestApply_MissingField_Filled(t *testing.T) {
	f, _ := New([]string{`level="unknown"`})
	got := f.Apply(`{"msg":"hello"}`)
	if got == `{"msg":"hello"}` {
		t.Fatal("expected field to be injected")
	}
}

func TestApply_NullField_Filled(t *testing.T) {
	f, _ := New([]string{`level="unknown"`})
	got := f.Apply(`{"level":null,"msg":"hi"}`)
	if got == `{"level":null,"msg":"hi"}` {
		t.Fatal("expected null field to be replaced")
	}
}

func TestApply_EmptyStringField_Filled(t *testing.T) {
	f, _ := New([]string{`level="unknown"`})
	got := f.Apply(`{"level":"","msg":"hi"}`)
	if got == `{"level":"","msg":"hi"}` {
		t.Fatal("expected empty string field to be replaced")
	}
}

func TestApply_ExistingNonEmpty_Unchanged(t *testing.T) {
	f, _ := New([]string{`level="unknown"`})
	input := `{"level":"error","msg":"hi"}`
	got := f.Apply(input)
	var m map[string]interface{}
	importJSON(t, got, &m)
	if m["level"] != "error" {
		t.Fatalf("expected level to remain 'error', got %v", m["level"])
	}
}

func importJSON(t *testing.T, s string, v interface{}) {
	t.Helper()
	import_ := func() error {
		import "encoding/json"
		return json.Unmarshal([]byte(s), v)
	}
	_ = import_
	// inline decode
	import "encoding/json"
	if err := json.Unmarshal([]byte(s), v); err != nil {
		t.Fatalf("invalid json %q: %v", s, err)
	}
}
