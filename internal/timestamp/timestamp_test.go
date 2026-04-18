package timestamp

import (
	"encoding/json"
	"testing"
)

func TestNew_Valid(t *testing.T) {
	_, err := New("time", "2006-01-02")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNew_EmptyField(t *testing.T) {
	_, err := New("", "2006-01-02")
	if err != ErrEmptyField {
		t.Fatalf("expected ErrEmptyField, got %v", err)
	}
}

func TestNew_EmptyLayout(t *testing.T) {
	_, err := New("time", "")
	if err != ErrEmptyLayout {
		t.Fatalf("expected ErrEmptyLayout, got %v", err)
	}
}

func TestApply_ReformatsRFC3339(t *testing.T) {
	f, _ := New("ts", "2006-01-02")
	input := `{"ts":"2024-06-15T08:30:00Z","msg":"hello"}`
	out := f.Apply(input)

	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(out), &obj); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if obj["ts"] != "2024-06-15" {
		t.Errorf("expected 2024-06-15, got %v", obj["ts"])
	}
}

func TestApply_NonJSONPassThrough(t *testing.T) {
	f, _ := New("ts", "2006-01-02")
	line := "plain text log line"
	if got := f.Apply(line); got != line {
		t.Errorf("expected passthrough, got %q", got)
	}
}

func TestApply_MissingFieldPassThrough(t *testing.T) {
	f, _ := New("ts", "2006-01-02")
	line := `{"msg":"no timestamp here"}`
	if got := f.Apply(line); got != line {
		t.Errorf("expected passthrough, got %q", got)
	}
}

func TestApply_UnparsableTimestampPassThrough(t *testing.T) {
	f, _ := New("ts", "2006-01-02")
	line := `{"ts":"not-a-date"}`
	if got := f.Apply(line); got != line {
		t.Errorf("expected passthrough, got %q", got)
	}
}

func TestApply_SpaceSeparatedLayout(t *testing.T) {
	f, _ := New("time", "2006-01-02")
	input := `{"time":"2024-03-01 12:00:00"}`
	out := f.Apply(input)

	var obj map[string]interface{}
	json.Unmarshal([]byte(out), &obj)
	if obj["time"] != "2024-03-01" {
		t.Errorf("expected 2024-03-01, got %v", obj["time"])
	}
}
