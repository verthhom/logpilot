package fieldmap_test

import (
	"encoding/json"
	"testing"

	"github.com/logpilot/internal/fieldmap"
)

func TestNew_ValidRules(t *testing.T) {
	m, err := fieldmap.New([]string{"msg=message", "ts=timestamp"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	rules := m.Rules()
	if rules["msg"] != "message" || rules["ts"] != "timestamp" {
		t.Errorf("unexpected rules: %v", rules)
	}
}

func TestNew_InvalidRule(t *testing.T) {
	for _, bad := range []string{"nodash", "=nokey", "noval="} {
		_, err := fieldmap.New([]string{bad})
		if err == nil {
			t.Errorf("expected error for rule %q", bad)
		}
	}
}

func TestApply_RemapsFields(t *testing.T) {
	m, _ := fieldmap.New([]string{"msg=message", "ts=timestamp"})
	input := `{"ts":"2024-01-01","msg":"hello","level":"info"}`
	out := m.Apply(input)
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(out), &obj); err != nil {
		t.Fatalf("output not valid JSON: %v", err)
	}
	if _, ok := obj["message"]; !ok {
		t.Error("expected field 'message'")
	}
	if _, ok := obj["timestamp"]; !ok {
		t.Error("expected field 'timestamp'")
	}
	if _, ok := obj["msg"]; ok {
		t.Error("old field 'msg' should be removed")
	}
}

func TestApply_NonJSONPassThrough(t *testing.T) {
	m, _ := fieldmap.New([]string{"msg=message"})
	line := "plain text log line"
	if got := m.Apply(line); got != line {
		t.Errorf("expected passthrough, got %q", got)
	}
}

func TestApply_NoRules(t *testing.T) {
	m, _ := fieldmap.New(nil)
	line := `{"msg":"hi"}`
	if got := m.Apply(line); got != line {
		t.Errorf("expected unchanged line, got %q", got)
	}
}

func TestApply_MissingField(t *testing.T) {
	m, _ := fieldmap.New([]string{"missing=new"})
	line := `{"msg":"hi"}`
	out := m.Apply(line)
	var obj map[string]interface{}
	json.Unmarshal([]byte(out), &obj)
	if _, ok := obj["new"]; ok {
		t.Error("field 'new' should not exist when source field is absent")
	}
}
