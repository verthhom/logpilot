package jsoncondition

import (
	"encoding/json"
	"testing"
)

func TestNew_Valid(t *testing.T) {
	c, err := New("level=error:alert=true")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil Conditioner")
	}
}

func TestNew_BadRule_NoColon(t *testing.T) {
	_, err := New("level=error")
	if err != ErrBadRule {
		t.Fatalf("expected ErrBadRule, got %v", err)
	}
}

func TestNew_BadRule_MissingEquals(t *testing.T) {
	_, err := New("levelERROR:alert=true")
	if err != ErrBadRule {
		t.Fatalf("expected ErrBadRule, got %v", err)
	}
}

func TestNew_BlankPart(t *testing.T) {
	_, err := New("=error:alert=true")
	if err != ErrBlankPart {
		t.Fatalf("expected ErrBlankPart, got %v", err)
	}
}

func TestApply_ConditionMatches(t *testing.T) {
	c, _ := New("level=error:alert=true")
	input := `{"level":"error","msg":"boom"}`
	out := c.Apply(input)
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(out), &obj); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if obj["alert"] != "true" {
		t.Errorf("expected alert=true, got %v", obj["alert"])
	}
}

func TestApply_ConditionNotMatched(t *testing.T) {
	c, _ := New("level=error:alert=true")
	input := `{"level":"info","msg":"ok"}`
	out := c.Apply(input)
	if out != input {
		t.Errorf("expected unchanged line, got %q", out)
	}
}

func TestApply_MissingSourceField(t *testing.T) {
	c, _ := New("level=error:alert=true")
	input := `{"msg":"no level here"}`
	out := c.Apply(input)
	if out != input {
		t.Errorf("expected unchanged line, got %q", out)
	}
}

func TestApply_NonJSONPassThrough(t *testing.T) {
	c, _ := New("level=error:alert=true")
	input := "plain text line"
	out := c.Apply(input)
	if out != input {
		t.Errorf("expected unchanged line, got %q", out)
	}
}

func TestApply_NonStringFieldPassThrough(t *testing.T) {
	c, _ := New("level=error:alert=true")
	input := `{"level":42}`
	out := c.Apply(input)
	if out != input {
		t.Errorf("expected unchanged line for non-string field, got %q", out)
	}
}
