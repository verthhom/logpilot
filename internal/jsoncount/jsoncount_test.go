package jsoncount_test

import (
	"encoding/json"
	"testing"

	"github.com/yourorg/logpilot/internal/jsoncount"
)

func TestNew_Valid(t *testing.T) {
	c, err := jsoncount.New("_count")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil Counter")
	}
}

func TestNew_BlankField(t *testing.T) {
	_, err := jsoncount.New("")
	if err == nil {
		t.Fatal("expected error for blank field")
	}
}

func TestApply_NonJSONPassThrough(t *testing.T) {
	c, _ := jsoncount.New("_count")
	input := "not json at all"
	if got := c.Apply(input); got != input {
		t.Fatalf("expected pass-through, got %q", got)
	}
}

func TestApply_InjectsCount(t *testing.T) {
	c, _ := jsoncount.New("_count")
	input := `{"level":"info","msg":"hello"}`
	out := c.Apply(input)

	var obj map[string]json.RawMessage
	if err := json.Unmarshal([]byte(out), &obj); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}

	raw, ok := obj["_count"]
	if !ok {
		t.Fatal("expected _count field in output")
	}

	// original had 2 fields; after injection total is 3, but count reflects
	// the number of keys at injection time (including the new field itself).
	if string(raw) != "3" {
		t.Fatalf("expected count 3, got %s", raw)
	}
}

func TestApply_EmptyObject(t *testing.T) {
	c, _ := jsoncount.New("n")
	out := c.Apply(`{}`)

	var obj map[string]json.RawMessage
	if err := json.Unmarshal([]byte(out), &obj); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if string(obj["n"]) != "1" {
		t.Fatalf("expected count 1 (just the injected field), got %s", obj["n"])
	}
}

func TestApply_OverwritesExistingField(t *testing.T) {
	c, _ := jsoncount.New("_count")
	// pre-existing _count should be overwritten
	input := `{"_count":99,"x":1}`
	out := c.Apply(input)

	var obj map[string]json.RawMessage
	if err := json.Unmarshal([]byte(out), &obj); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	// map had 2 keys before overwrite; after overwrite still 2 keys
	if string(obj["_count"]) != "2" {
		t.Fatalf("expected count 2, got %s", obj["_count"])
	}
}
