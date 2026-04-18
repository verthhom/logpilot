package aggregate

import (
	"encoding/json"
	"testing"
	"time"
)

func TestNew_Valid(t *testing.T) {
	agg, err := New("level", 10*time.Second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	agg.Stop()
}

func TestNew_EmptyField(t *testing.T) {
	_, err := New("", time.Second)
	if err != ErrEmptyField {
		t.Fatalf("expected ErrEmptyField, got %v", err)
	}
}

func TestNew_InvalidWindow(t *testing.T) {
	_, err := New("level", 0)
	if err != ErrInvalidWindow {
		t.Fatalf("expected ErrInvalidWindow, got %v", err)
	}
}

func TestFeed_NonJSON_Ignored(t *testing.T) {
	agg, _ := New("level", time.Hour)
	agg.Feed("not json")
	agg.Stop()
	// no summary emitted for non-JSON
	if _, ok := <-agg.Out(); ok {
		t.Fatal("expected no output for non-JSON input")
	}
}

func TestFeed_MissingField_Ignored(t *testing.T) {
	agg, _ := New("level", time.Hour)
	agg.Feed(`{"msg":"hello"}`)
	agg.Stop()
	if _, ok := <-agg.Out(); ok {
		t.Fatal("expected no output when field absent")
	}
}

func TestStop_EmitsSummary(t *testing.T) {
	agg, _ := New("level", time.Hour)
	agg.Feed(`{"level":"info","msg":"a"}`)
	agg.Feed(`{"level":"info","msg":"b"}`)
	agg.Feed(`{"level":"error","msg":"c"}`)
	agg.Stop()

	line, ok := <-agg.Out()
	if !ok {
		t.Fatal("expected a summary line")
	}
	var m map[string]any
	if err := json.Unmarshal([]byte(line), &m); err != nil {
		t.Fatalf("summary not valid JSON: %v", err)
	}
	counts, ok := m["counts"].(map[string]any)
	if !ok {
		t.Fatal("counts field missing or wrong type")
	}
	if counts["info"].(float64) != 2 {
		t.Errorf("expected info=2, got %v", counts["info"])
	}
	if counts["error"].(float64) != 1 {
		t.Errorf("expected error=1, got %v", counts["error"])
	}
}
