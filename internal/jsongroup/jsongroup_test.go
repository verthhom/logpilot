package jsongroup

import (
	"encoding/json"
	"testing"
)

func TestNew_Valid(t *testing.T) {
	g, err := New("service")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if g == nil {
		t.Fatal("expected non-nil Grouper")
	}
}

func TestNew_BlankField(t *testing.T) {
	_, err := New("")
	if err == nil {
		t.Fatal("expected error for blank field")
	}
	if err != ErrBlankField {
		t.Fatalf("want ErrBlankField, got %v", err)
	}
}

func TestNew_WhitespaceField(t *testing.T) {
	_, err := New("   ")
	if err != ErrBlankField {
		t.Fatalf("want ErrBlankField, got %v", err)
	}
}

func TestFeed_NonJSONIgnored(t *testing.T) {
	g, _ := New("svc")
	g.Feed("not json at all")
	if g.Len() != 0 {
		t.Fatalf("expected 0 buffered lines, got %d", g.Len())
	}
}

func TestFeed_MissingFieldIgnored(t *testing.T) {
	g, _ := New("svc")
	g.Feed(`{"other":"value"}`)
	if g.Len() != 0 {
		t.Fatalf("expected 0 buffered lines, got %d", g.Len())
	}
}

func TestFeed_GroupsByField(t *testing.T) {
	g, _ := New("service")
	g.Feed(`{"service":"auth","msg":"login"}`)
	g.Feed(`{"service":"auth","msg":"logout"}`)
	g.Feed(`{"service":"api","msg":"request"}`)
	if g.Len() != 3 {
		t.Fatalf("expected 3 buffered lines, got %d", g.Len())
	}
}

func TestFlush_ReturnsOneSummaryPerGroup(t *testing.T) {
	g, _ := New("service")
	g.Feed(`{"service":"auth","msg":"a"}`)
	g.Feed(`{"service":"auth","msg":"b"}`)
	g.Feed(`{"service":"api","msg":"c"}`)

	summaries := g.Flush()
	if len(summaries) != 2 {
		t.Fatalf("expected 2 summaries, got %d", len(summaries))
	}
}

func TestFlush_SummaryContainsCount(t *testing.T) {
	g, _ := New("service")
	g.Feed(`{"service":"auth","msg":"a"}`)
	g.Feed(`{"service":"auth","msg":"b"}`)

	summaries := g.Flush()
	if len(summaries) != 1 {
		t.Fatalf("expected 1 summary, got %d", len(summaries))
	}
	var obj map[string]json.RawMessage
	if err := json.Unmarshal([]byte(summaries[0]), &obj); err != nil {
		t.Fatalf("summary is not valid JSON: %v", err)
	}
	var count int
	if err := json.Unmarshal(obj["count"], &count); err != nil {
		t.Fatalf("cannot parse count: %v", err)
	}
	if count != 2 {
		t.Fatalf("expected count=2, got %d", count)
	}
}

func TestFlush_ResetsState(t *testing.T) {
	g, _ := New("service")
	g.Feed(`{"service":"auth","msg":"x"}`)
	g.Flush()
	if g.Len() != 0 {
		t.Fatalf("expected 0 after flush, got %d", g.Len())
	}
	if summaries := g.Flush(); len(summaries) != 0 {
		t.Fatalf("expected empty second flush, got %d summaries", len(summaries))
	}
}
