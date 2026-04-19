package jsonrename_test

import (
	"encoding/json"
	"testing"

	"github.com/yourorg/logpilot/internal/jsonrename"
)

func TestRenamer_ChainedApplication(t *testing.T) {
	r, err := jsonrename.New([]string{"msg:message", "lvl:level"})
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	lines := []string{
		`{"msg":"started","lvl":"info"}`,
		`{"msg":"failed","lvl":"error","code":500}`,
		"plain text line",
	}

	for _, line := range lines[:2] {
		out := r.Apply(line)
		var obj map[string]json.RawMessage
		if err := json.Unmarshal([]byte(out), &obj); err != nil {
			t.Fatalf("invalid json output: %v", err)
		}
		if _, ok := obj["message"]; !ok {
			t.Errorf("missing 'message' in: %s", out)
		}
		if _, ok := obj["level"]; !ok {
			t.Errorf("missing 'level' in: %s", out)
		}
	}

	if got := r.Apply(lines[2]); got != lines[2] {
		t.Errorf("plain line should pass through, got %q", got)
	}
}
