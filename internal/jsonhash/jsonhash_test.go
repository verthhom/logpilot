package jsonhash

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"testing"
)

func TestNew_Valid(t *testing.T) {
	h, err := New("_hash", []string{"msg", "level"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if h == nil {
		t.Fatal("expected non-nil Hasher")
	}
}

func TestNew_BlankDest(t *testing.T) {
	_, err := New("", []string{"msg"})
	if err == nil {
		t.Fatal("expected error for blank dest")
	}
}

func TestNew_NoFields(t *testing.T) {
	_, err := New("_hash", []string{})
	if err == nil {
		t.Fatal("expected error for empty fields slice")
	}
}

func TestNew_BlankField(t *testing.T) {
	_, err := New("_hash", []string{"msg", ""})
	if err == nil {
		t.Fatal("expected error for blank field name")
	}
}

func TestApply_NonJSONPassThrough(t *testing.T) {
	h, _ := New("_hash", []string{"msg"})
	input := "not json at all"
	if got := h.Apply(input); got != input {
		t.Fatalf("expected pass-through, got %q", got)
	}
}

func TestApply_InjectsHash(t *testing.T) {
	h, _ := New("_hash", []string{"msg"})
	input := `{"msg":"hello","level":"info"}`
	output := h.Apply(input)

	var obj map[string]json.RawMessage
	if err := json.Unmarshal([]byte(output), &obj); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if _, ok := obj["_hash"]; !ok {
		t.Fatal("expected _hash field in output")
	}
}

func TestApply_HashIsCorrect(t *testing.T) {
	h, _ := New("_hash", []string{"msg"})
	input := `{"msg":"hello"}`
	output := h.Apply(input)

	var obj map[string]string
	if err := json.Unmarshal([]byte(output), &obj); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}

	// The hasher concatenates raw JSON values: "hello" (with quotes).
	sum := sha256.Sum256([]byte(`"hello"`))
	want := fmt.Sprintf("%x", sum)
	if obj["_hash"] != want {
		t.Fatalf("hash mismatch: got %q, want %q", obj["_hash"], want)
	}
}

func TestApply_MissingFieldIgnored(t *testing.T) {
	h, _ := New("_hash", []string{"missing"})
	input := `{"msg":"hello"}`
	output := h.Apply(input)

	// Hash of empty string
	sum := sha256.Sum256([]byte(""))
	want := fmt.Sprintf("%x", sum)

	var obj map[string]string
	if err := json.Unmarshal([]byte(output), &obj); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if obj["_hash"] != want {
		t.Fatalf("expected hash of empty string, got %q", obj["_hash"])
	}
}
