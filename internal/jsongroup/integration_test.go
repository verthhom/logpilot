package jsongroup_test

import (
	"encoding/json"
	"testing"

	"github.com/yourorg/logpilot/internal/jsongroup"
)

func TestGrouper_FullRoundTrip(t *testing.T) {
	lines := []string{
		`{"env":"prod","level":"error","msg":"disk full"}`,
		`{"env":"prod","level":"info","msg":"started"}`,
		`{"env":"staging","level":"warn","msg":"slow query"}`,
		`{"env":"prod","level":"error","msg":"oom"}`,
		`not json`,
		`{"other":"field"}`,
	}

	g, err := jsongroup.New("env")
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	for _, l := range lines {
		g.Feed(l)
	}

	// 4 valid JSON lines with "env" field fed
	if g.Len() != 4 {
		t.Fatalf("expected 4 buffered, got %d", g.Len())
	}

	summaries := g.Flush()
	if len(summaries) != 2 {
		t.Fatalf("expected 2 groups (prod, staging), got %d", len(summaries))
	}

	counts := map[string]int{}
	for _, s := range summaries {
		var obj map[string]json.RawMessage
		if err := json.Unmarshal([]byte(s), &obj); err != nil {
			t.Fatalf("invalid summary JSON: %v", err)
		}
		var env string
		if err := json.Unmarshal(obj["env"], &env); err != nil {
			t.Fatalf("cannot decode env: %v", err)
		}
		var count int
		if err := json.Unmarshal(obj["count"], &count); err != nil {
			t.Fatalf("cannot decode count: %v", err)
		}
		counts[env] = count
	}

	if counts["prod"] != 3 {
		t.Errorf("expected prod count=3, got %d", counts["prod"])
	}
	if counts["staging"] != 1 {
		t.Errorf("expected staging count=1, got %d", counts["staging"])
	}
}
