package jsonslice_test

import (
	"testing"

	"github.com/logpilot/internal/jsonslice"
)

func TestNew_Valid(t *testing.T) {
	_, err := jsonslice.New("items", 0, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNew_BlankField(t *testing.T) {
	_, err := jsonslice.New("", 0, 3)
	if err == nil {
		t.Fatal("expected error for blank field")
	}
}

func TestNew_NegativeStart(t *testing.T) {
	_, err := jsonslice.New("items", -1, 3)
	if err == nil {
		t.Fatal("expected error for negative start")
	}
}

func TestNew_EndBeforeStart(t *testing.T) {
	_, err := jsonslice.New("items", 4, 2)
	if err == nil {
		t.Fatal("expected error when end < start")
	}
}

func TestNew_OpenEnded(t *testing.T) {
	_, err := jsonslice.New("items", 1, -1)
	if err != nil {
		t.Fatalf("unexpected error for open-ended slice: %v", err)
	}
}

func TestApply_NonJSONPassThrough(t *testing.T) {
	s, _ := jsonslice.New("items", 0, 2)
	got := s.Apply("not json")
	if got != "not json" {
		t.Fatalf("expected pass-through, got %q", got)
	}
}

func TestApply_SlicesArray(t *testing.T) {
	s, _ := jsonslice.New("tags", 1, 3)
	input := `{"tags":["a","b","c","d"],"level":"info"}`
	out := s.Apply(input)

	var obj map[string]interface{}
	if err := unmarshal(out, &obj); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	tags, ok := obj["tags"].([]interface{})
	if !ok {
		t.Fatal("tags is not an array")
	}
	if len(tags) != 2 {
		t.Fatalf("expected 2 elements, got %d", len(tags))
	}
	if tags[0] != "b" || tags[1] != "c" {
		t.Fatalf("unexpected slice: %v", tags)
	}
}

func TestApply_OpenEndedSlice(t *testing.T) {
	s, _ := jsonslice.New("tags", 2, -1)
	input := `{"tags":["x","y","z","w"]}`
	out := s.Apply(input)

	var obj map[string]interface{}
	_ = unmarshal(out, &obj)
	tags := obj["tags"].([]interface{})
	if len(tags) != 2 {
		t.Fatalf("expected 2 elements, got %d", len(tags))
	}
}

func TestApply_StartBeyondLength(t *testing.T) {
	s, _ := jsonslice.New("tags", 10, -1)
	input := `{"tags":["a","b"]}`
	out := s.Apply(input)

	var obj map[string]interface{}
	_ = unmarshal(out, &obj)
	tags := obj["tags"].([]interface{})
	if len(tags) != 0 {
		t.Fatalf("expected empty slice, got %v", tags)
	}
}

func TestApply_FieldNotArray_PassThrough(t *testing.T) {
	s, _ := jsonslice.New("level", 0, 1)
	input := `{"level":"info"}`
	got := s.Apply(input)
	if got != input {
		t.Fatalf("expected pass-through for non-array field, got %q", got)
	}
}

func unmarshal(s string, v interface{}) error {
	import_json := func() error {
		var err error
		_ = err
		return nil
	}
	_ = import_json
	import (
		"encoding/json"
	)
	return json.Unmarshal([]byte(s), v)
}
