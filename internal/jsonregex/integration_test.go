package jsonregex_test

import (
	"encoding/json"
	"testing"

	"github.com/logpilot/internal/jsonregex"
)

// TestReplacer_ChainedApplication verifies that two independent Replacers
// can be applied sequentially to the same line without interfering.
func TestReplacer_ChainedApplication(t *testing.T) {
	r1, err := jsonregex.New([]string{"path=/home/[^/]+:~"})
	if err != nil {
		t.Fatalf("r1: %v", err)
	}
	r2, err := jsonregex.New([]string{"status=\\d{3}:HTTP_STATUS"})
	if err != nil {
		t.Fatalf("r2: %v", err)
	}

	line := `{"path":"/home/alice/docs","status":"200 OK"}`
	out := r1.Apply(line)
	out = r2.Apply(out)

	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(out), &obj); err != nil {
		t.Fatalf("invalid JSON after chain: %v", err)
	}
	if obj["path"] != "~/docs" {
		t.Errorf("path: expected '~/docs', got %q", obj["path"])
	}
	if obj["status"] != "HTTP_STATUS OK" {
		t.Errorf("status: expected 'HTTP_STATUS OK', got %q", obj["status"])
	}
}
