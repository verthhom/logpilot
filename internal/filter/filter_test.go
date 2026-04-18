package filter

import (
	"testing"
)

func TestNew_ValidRules(t *testing.T) {
	f, err := New([]string{"level:eq:error", "service:contains:auth", "request_id:exists"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(f.Rules) != 3 {
		t.Fatalf("expected 3 rules, got %d", len(f.Rules))
	}
}

func TestNew_InvalidRule(t *testing.T) {
	_, err := New([]string{"badrule"})
	if err == nil {
		t.Fatal("expected error for invalid rule")
	}
}

func TestMatch_Eq(t *testing.T) {
	f, _ := New([]string{"level:eq:error"})
	entry := map[string]interface{}{"level": "error", "msg": "boom"}
	if !f.Match(entry) {
		t.Error("expected match")
	}
	entry["level"] = "info"
	if f.Match(entry) {
		t.Error("expected no match")
	}
}

func TestMatch_Contains(t *testing.T) {
	f, _ := New([]string{"msg:contains:timeout"})
	entry := map[string]interface{}{"msg": "connection timeout occurred"}
	if !f.Match(entry) {
		t.Error("expected match")
	}
}

func TestMatch_Exists(t *testing.T) {
	f, _ := New([]string{"trace_id:exists"})
	if f.Match(map[string]interface{}{"level": "info"}) {
		t.Error("expected no match without trace_id")
	}
	if !f.Match(map[string]interface{}{"trace_id": "abc123"}) {
		t.Error("expected match with trace_id")
	}
}

func TestMatch_MultipleRules(t *testing.T) {
	f, _ := New([]string{"level:eq:error", "service:contains:auth"})
	entry := map[string]interface{}{"level": "error", "service": "auth-service"}
	if !f.Match(entry) {
		t.Error("expected match")
	}
	entry["service"] = "billing"
	if f.Match(entry) {
		t.Error("expected no match")
	}
}

func TestMatch_EmptyFilter(t *testing.T) {
	f, _ := New(nil)
	if !f.Match(map[string]interface{}{"level": "info"}) {
		t.Error("empty filter should match everything")
	}
}
