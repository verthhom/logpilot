package redact

import (
	"encoding/json"
	"testing"
)

func TestNew_Valid(t *testing.T) {
	r, err := New([]string{"password", "token"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r == nil {
		t.Fatal("expected non-nil Redactor")
	}
}

func TestNew_EmptySlice(t *testing.T) {
	_, err := New([]string{})
	if err == nil {
		t.Fatal("expected error for empty fields")
	}
}

func TestNew_BlankField(t *testing.T) {
	_, err := New([]string{"  "})
	if err == nil {
		t.Fatal("expected error for blank field name")
	}
}

func TestApply_RedactsField(t *testing.T) {
	r, _ := New([]string{"password"})
	input := `{"user":"alice","password":"s3cr3t"}`
	out := r.Apply(input)

	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(out), &obj); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if obj["password"] != redactedValue {
		t.Errorf("expected password to be %q, got %v", redactedValue, obj["password"])
	}
	if obj["user"] != "alice" {
		t.Errorf("expected user to be unchanged")
	}
}

func TestApply_NonJSONPassThrough(t *testing.T) {
	r, _ := New([]string{"password"})
	input := "plain text log line"
	if got := r.Apply(input); got != input {
		t.Errorf("expected passthrough, got %q", got)
	}
}

func TestApply_MissingFieldNoChange(t *testing.T) {
	r, _ := New([]string{"secret"})
	input := `{"user":"bob","level":"info"}`
	out := r.Apply(input)

	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(out), &obj); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if _, ok := obj["secret"]; ok {
		t.Error("secret field should not exist")
	}
}

func TestApply_MultipleFields(t *testing.T) {
	r, _ := New([]string{"token", "ssn"})
	input := `{"name":"carol","token":"abc123","ssn":"123-45-6789"}`
	out := r.Apply(input)

	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(out), &obj); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	for _, f := range []string{"token", "ssn"} {
		if obj[f] != redactedValue {
			t.Errorf("expected %s to be redacted, got %v", f, obj[f])
		}
	}
}
