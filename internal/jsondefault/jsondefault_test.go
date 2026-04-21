package jsondefault

import (
	"encoding/json"
	"testing"
)

func TestNew_Valid(t *testing.T) {
	d, err := New([]string{`level="info"`, `retries=3`})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d == nil {
		t.Fatal("expected non-nil Defaulter")
	}
}

func TestNew_NoRules(t *testing.T) {
	_, err := New(nil)
	if err != ErrNoRules {
		t.Fatalf("expected ErrNoRules, got %v", err)
	}
}

func TestNew_BadRule(t *testing.T) {
	_, err := New([]string{"nodivider"})
	if err == nil {
		t.Fatal("expected error for bad rule")
	}
}

func TestNew_BlankField(t *testing.T) {
	_, err := New([]string{`=value`})
	if err == nil {
		t.Fatal("expected error for blank field")
	}
}

func TestNew_InvalidJSONValue(t *testing.T) {
	_, err := New([]string{`level=not-json`})
	if err == nil {
		t.Fatal("expected error for invalid JSON value")
	}
}

func TestApply_SetsMissingField(t *testing.T) {
	d, _ := New([]string{`level="info"`})
	result := d.Apply(`{"msg":"hello"}`)
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(result), &obj); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	if obj["level"] != "info" {
		t.Errorf("expected level=info, got %v", obj["level"])
	}
}

func TestApply_DoesNotOverwriteExisting(t *testing.T) {
	d, _ := New([]string{`level="info"`})
	result := d.Apply(`{"level":"error","msg":"oops"}`)
	var obj map[string]interface{}
	json.Unmarshal([]byte(result), &obj)
	if obj["level"] != "error" {
		t.Errorf("expected existing level=error to be preserved, got %v", obj["level"])
	}
}

func TestApply_OverwritesNull(t *testing.T) {
	d, _ := New([]string{`level="info"`})
	result := d.Apply(`{"level":null,"msg":"test"}`)
	var obj map[string]interface{}
	json.Unmarshal([]byte(result), &obj)
	if obj["level"] != "info" {
		t.Errorf("expected null level to be replaced with info, got %v", obj["level"])
	}
}

func TestApply_NonJSONPassThrough(t *testing.T) {
	d, _ := New([]string{`level="info"`})
	plain := "not a json line"
	if got := d.Apply(plain); got != plain {
		t.Errorf("expected pass-through, got %q", got)
	}
}

func TestApply_NumericDefault(t *testing.T) {
	d, _ := New([]string{`retries=3`})
	result := d.Apply(`{"msg":"retry"}`)
	var obj map[string]interface{}
	json.Unmarshal([]byte(result), &obj)
	if obj["retries"] != float64(3) {
		t.Errorf("expected retries=3, got %v", obj["retries"])
	}
}
