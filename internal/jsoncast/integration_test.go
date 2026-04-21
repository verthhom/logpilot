package jsoncast_test

import (
	"encoding/json"
	"testing"

	"github.com/user/logpilot/internal/jsoncast"
)

// TestCaster_MultipleRulesApplied verifies that all configured rules are
// applied in a single pass over the same JSON object.
func TestCaster_MultipleRulesApplied(t *testing.T) {
	c, err := jsoncast.New([]string{"status:int", "latency:float", "retried:bool", "trace:string"})
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	input := `{"status":"201","latency":"0.75","retried":"false","trace":99}`
	out := c.Apply(input)

	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(out), &obj); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}

	if _, ok := obj["status"].(float64); !ok {
		t.Errorf("status: expected numeric, got %T", obj["status"])
	}
	if _, ok := obj["latency"].(float64); !ok {
		t.Errorf("latency: expected float64, got %T", obj["latency"])
	}
	if v, ok := obj["retried"].(bool); !ok || v {
		t.Errorf("retried: expected bool false, got %v (%T)", obj["retried"], obj["retried"])
	}
	if _, ok := obj["trace"].(string); !ok {
		t.Errorf("trace: expected string, got %T", obj["trace"])
	}
}
