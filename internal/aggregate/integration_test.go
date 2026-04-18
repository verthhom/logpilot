package aggregate_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/logpilot/internal/aggregate"
)

func TestAggregator_WindowFlush(t *testing.T) {
	window := 80 * time.Millisecond
	agg, err := aggregate.New("service", window)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer agg.Stop()

	lines := []string{
		`{"service":"auth","msg":"login"}`,
		`{"service":"auth","msg":"logout"}`,
		`{"service":"api","msg":"request"}`,
	}
	for _, l := range lines {
		agg.Feed(l)
	}

	select {
	case summary := <-agg.Out():
		var m map[string]any
		if err := json.Unmarshal([]byte(summary), &m); err != nil {
			t.Fatalf("invalid JSON: %v", err)
		}
		if m["field"] != "service" {
			t.Errorf("expected field=service, got %v", m["field"])
		}
		counts := m["counts"].(map[string]any)
		if counts["auth"].(float64) != 2 {
			t.Errorf("expected auth=2, got %v", counts["auth"])
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timed out waiting for aggregate summary")
	}
}
