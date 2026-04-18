package multiline_test

import (
	"testing"

	"github.com/logpilot/internal/multiline"
)

func TestJoiner_FullSequence(t *testing.T) {
	j, err := multiline.New(`^\{`, "")
	if err != nil {
		t.Fatalf("setup: %v", err)
	}

	input := []string{
		`{"level":"info",`,
		`"msg":"hello"}`,
		`{"level":"error",`,
		`"msg":"oops",`,
		`"code":500}`,
	}

	var events []string
	for _, line := range input {
		if ev, ok := j.Feed(line); ok {
			events = append(events, ev)
		}
	}
	if ev, ok := j.Flush(); ok {
		events = append(events, ev)
	}

	if len(events) != 2 {
		t.Fatalf("expected 2 events, got %d: %v", len(events), events)
	}

	want0 := `{"level":"info","msg":"hello"}`
	if events[0] != want0 {
		t.Errorf("event[0]: got %q, want %q", events[0], want0)
	}

	want1 := `{"level":"error","msg":"oops","code":500}`
	if events[1] != want1 {
		t.Errorf("event[1]: got %q, want %q", events[1], want1)
	}
}
