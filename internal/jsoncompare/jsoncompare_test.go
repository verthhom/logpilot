package jsoncompare

import (
	"encoding/json"
	"testing"
)

func TestNew_Valid(t *testing.T) {
	c, err := New("a:b:result")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.left != "a" || c.right != "b" || c.dest != "result" {
		t.Fatalf("unexpected fields: %+v", c)
	}
}

func TestNew_BadRule(t *testing.T) {
	_, err := New("a:b")
	if err != ErrBadRule {
		t.Fatalf("expected ErrBadRule, got %v", err)
	}
}

func TestNew_BlankPart(t *testing.T) {
	for _, rule := range []string{":b:dest", "a::dest", "a:b:"} {
		_, err := New(rule)
		if err != ErrBlankPart {
			t.Fatalf("rule %q: expected ErrBlankPart, got %v", rule, err)
		}
	}
}

func result(t *testing.T, line string) string {
	t.Helper()
	var obj map[string]any
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	v, _ := obj["cmp"].(string)
	return v
}

func TestApply_LessThan(t *testing.T) {
	c, _ := New("x:y:cmp")
	out := c.Apply(`{"x":1,"y":5}`)
	if got := result(t, out); got != "lt" {
		t.Fatalf("expected lt, got %q", got)
	}
}

func TestApply_GreaterThan(t *testing.T) {
	c, _ := New("x:y:cmp")
	out := c.Apply(`{"x":10,"y":2}`)
	if got := result(t, out); got != "gt" {
		t.Fatalf("expected gt, got %q", got)
	}
}

func TestApply_Equal(t *testing.T) {
	c, _ := New("x:y:cmp")
	out := c.Apply(`{"x":7,"y":7}`)
	if got := result(t, out); got != "eq" {
		t.Fatalf("expected eq, got %q", got)
	}
}

func TestApply_MissingField_PassThrough(t *testing.T) {
	c, _ := New("x:y:cmp")
	input := `{"x":3}`
	out := c.Apply(input)
	if out != input {
		t.Fatalf("expected pass-through, got %q", out)
	}
}

func TestApply_NonJSONPassThrough(t *testing.T) {
	c, _ := New("x:y:cmp")
	input := "not json at all"
	if out := c.Apply(input); out != input {
		t.Fatalf("expected pass-through, got %q", out)
	}
}
