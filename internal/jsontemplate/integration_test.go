package jsontemplate_test

import (
	"encoding/json"
	"testing"

	"github.com/yourorg/logpilot/internal/jsontemplate"
)

// TestApplier_ChainedApplication verifies that two Appliers can be chained so
// that the second template can reference a field injected by the first.
func TestApplier_ChainedApplication(t *testing.T) {
	first, err := jsontemplate.New("full_msg", "[{{.level}}] {{.msg}}")
	if err != nil {
		t.Fatalf("first applier: %v", err)
	}
	second, err := jsontemplate.New("alert", "ALERT: {{.full_msg}}")
	if err != nil {
		t.Fatalf("second applier: %v", err)
	}

	line := `{"level":"warn","msg":"low memory"}`
	intermediate := first.Apply(line)
	final := second.Apply(intermediate)

	var obj map[string]any
	if err := json.Unmarshal([]byte(final), &obj); err != nil {
		t.Fatalf("final output not valid JSON: %v", err)
	}

	if obj["full_msg"] != "[warn] low memory" {
		t.Errorf("unexpected full_msg: %v", obj["full_msg"])
	}
	if obj["alert"] != "ALERT: [warn] low memory" {
		t.Errorf("unexpected alert: %v", obj["alert"])
	}
}
