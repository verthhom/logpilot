package jsonround_test

import (
	"encoding/json"
	"testing"

	"github.com/logpilot/internal/jsonround"
)

func TestNew_Valid(t *testing.T) {
	_, err := jsonround.New([]string{"latency:2", "ratio:4"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNew_NoRules(t *testing.T) {
	_, err := jsonround.New(nil)
	if err == nil {
		t.Fatal("expected error for empty rules")
	}
}

func TestNew_BadRule_NoColon(t *testing.T) {
	_, err := jsonround.New([]string{"latency2"})
	if err == nil {
		t.Fatal("expected error for rule without colon")
	}
}

func TestNew_BadRule_NegativePlaces(t *testing.T) {
	_, err := jsonround.New([]string{"latency:-1"})
	if err == nil {
		t.Fatal("expected error for negative places")
	}
}

func TestNew_BadRule_BlankField(t *testing.T) {
	_, err := jsonround.New([]string{":2"})
	if err == nil {
		t.Fatal("expected error for blank field")
	}
}

func TestApply_RoundsField(t *testing.T) {
	r, _ := jsonround.New([]string{"latency:2"})
	out := r.Apply(`{"latency":3.14159,"msg":"ok"}`)
	var obj map[string]json.RawMessage
	if err := json.Unmarshal([]byte(out), &obj); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	var v float64
	if err := json.Unmarshal(obj["latency"], &v); err != nil {
		t.Fatalf("cannot read latency: %v", err)
	}
	if v != 3.14 {
		t.Errorf("expected 3.14, got %v", v)
	}
}

func TestApply_ZeroPlaces(t *testing.T) {
	r, _ := jsonround.New([]string{"score:0"})
	out := r.Apply(`{"score":7.9}`)
	var obj map[string]json.RawMessage
	json.Unmarshal([]byte(out), &obj) //nolint:errcheck
	var v float64
	json.Unmarshal(obj["score"], &v) //nolint:errcheck
	if v != 8 {
		t.Errorf("expected 8, got %v", v)
	}
}

func TestApply_NonNumericField_Skipped(t *testing.T) {
	r, _ := jsonround.New([]string{"msg:2"})
	input := `{"msg":"hello"}`
	out := r.Apply(input)
	if out != input {
		t.Errorf("expected pass-through, got %q", out)
	}
}

func TestApply_NonJSONPassThrough(t *testing.T) {
	r, _ := jsonround.New([]string{"latency:2"})
	input := "not json at all"
	if got := r.Apply(input); got != input {
		t.Errorf("expected %q, got %q", input, got)
	}
}

func TestApply_MissingField_Ignored(t *testing.T) {
	r, _ := jsonround.New([]string{"missing:3"})
	input := `{"other":1}`
	out := r.Apply(input)
	var obj map[string]json.RawMessage
	if err := json.Unmarshal([]byte(out), &obj); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if _, ok := obj["missing"]; ok {
		t.Error("missing field should not be injected")
	}
}
