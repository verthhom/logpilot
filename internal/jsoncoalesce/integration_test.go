package jsoncoalesce_test

import (
	"encoding/json"
	"testing"

	"github.com/logpilot/internal/jsoncoalesce"
)

// TestCoalescer_ChainedApplication verifies that two coalescers can be
// applied in sequence to normalise different field groups independently.
func TestCoalescer_ChainedApplication(t *testing.T) {
	msgCoalescer, err := jsoncoalesce.New("message", []string{"msg", "text", "body"})
	if err != nil {
		t.Fatalf("setup error: %v", err)
	}
	errCoalescer, err := jsoncoalesce.New("error", []string{"err", "error_msg"})
	if err != nil {
		t.Fatalf("setup error: %v", err)
	}

	lines := []struct {
		input       string
		wantMessage string
		wantError   string
	}{
		{`{"msg":"started","err":"none"}`, "started", "none"},
		{`{"body":"ready","error_msg":"timeout"}`, "ready", "timeout"},
		{`{"level":"info"}`, "", ""},
	}

	for _, tc := range lines {
		out := msgCoalescer.Apply(tc.input)
		out = errCoalescer.Apply(out)

		var obj map[string]interface{}
		if err := json.Unmarshal([]byte(out), &obj); err != nil {
			t.Fatalf("invalid JSON for input %q: %v", tc.input, err)
		}

		gotMsg, _ := obj["message"].(string)
		gotErr, _ := obj["error"].(string)

		if gotMsg != tc.wantMessage {
			t.Errorf("input %q: message = %q, want %q", tc.input, gotMsg, tc.wantMessage)
		}
		if gotErr != tc.wantError {
			t.Errorf("input %q: error = %q, want %q", tc.input, gotErr, tc.wantError)
		}
	}
}
