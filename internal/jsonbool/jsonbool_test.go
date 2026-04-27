package jsonbool

import (
	"encoding/json"
	"testing"
)

func TestNew_Valid(t *testing.T) {
	_, err := New([]string{"active", "enabled"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNew_NoFields(t *testing.T) {
	_, err := New([]string{})
	if err == nil {
		t.Fatal("expected error for empty fields slice")
	}
}

func TestNew_BlankField(t *testing.T) {
	_, err := New([]string{"active", "  "})
	if err == nil {
		t.Fatal("expected error for blank field name")
	}
}

func TestApply_NonJSONPassThrough(t *testing.T) {
	n, _ := New([]string{"active"})
	input := "not json at all"
	if got := n.Apply(input); got != input {
		t.Fatalf("expected pass-through, got %q", got)
	}
}

func TestApply_TruthyStrings(t *testing.T) {
	cases := []struct {
		input string
		want bool
	}{
		{`{"active":"true"}`, true},
		{`{"active":"1"}`, true},
		{`{"active":"yes"}`, true},
		{`{"active":"on"}`, true},
		{`{"active":"TRUE"}`, true},
	}
	n, _ := New([]string{"active"})
	for _, tc := range cases {
		out := n.Apply(tc.input)
		var obj map[string]interface{}
		if err := json.Unmarshal([]byte(out), &obj); err != nil {
			t.Fatalf("invalid JSON output for %q: %v", tc.input, err)
		}
		got, ok := obj["active"].(bool)
		if !ok || got != tc.want {
			t.Errorf("input %q: want active=%v, got %v", tc.input, tc.want, obj["active"])
		}
	}
}

func TestApply_FalsyStrings(t *testing.T) {
	cases := []string{
		`{"active":"false"}`,
		`{"active":"0"}`,
		`{"active":"no"}`,
		`{"active":"off"}`,
	}
	n, _ := New([]string{"active"})
	for _, tc := range cases {
		out := n.Apply(tc)
		var obj map[string]interface{}
		if err := json.Unmarshal([]byte(out), &obj); err != nil {
			t.Fatalf("invalid JSON output for %q: %v", tc, err)
		}
		got, ok := obj["active"].(bool)
		if !ok || got != false {
			t.Errorf("input %q: want active=false, got %v", tc, obj["active"])
		}
	}
}

func TestApply_AlreadyBoolean(t *testing.T) {
	n, _ := New([]string{"active"})
	input := `{"active":true,"name":"test"}`
	out := n.Apply(input)
	// Should be unchanged (no re-marshal needed, but value preserved).
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(out), &obj); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if obj["active"] != true {
		t.Errorf("expected active=true, got %v", obj["active"])
	}
}

func TestApply_UnknownStringUnchanged(t *testing.T) {
	n, _ := New([]string{"active"})
	input := `{"active":"maybe"}`
	out := n.Apply(input)
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(out), &obj); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if obj["active"] != "maybe" {
		t.Errorf("expected active to remain 'maybe', got %v", obj["active"])
	}
}
