package jsonclamp

import (
	"testing"
)

func TestNew_Valid(t *testing.T) {
	c, err := New([]string{"latency:0:1000"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil Clamper")
	}
}

func TestNew_NoRules(t *testing.T) {
	_, err := New(nil)
	if err == nil {
		t.Fatal("expected error for empty rules")
	}
}

func TestNew_BadRule(t *testing.T) {
	_, err := New([]string{"badformat"})
	if err == nil {
		t.Fatal("expected error for bad rule format")
	}
}

func TestNew_BlankField(t *testing.T) {
	_, err := New([]string{":0:100"})
	if err == nil {
		t.Fatal("expected error for blank field")
	}
}

func TestNew_InvalidRange(t *testing.T) {
	_, err := New([]string{"score:100:0"})
	if err == nil {
		t.Fatal("expected error when min > max")
	}
}

func TestApply_ClampsBelow(t *testing.T) {
	c, _ := New([]string{"score:0:100"})
	out := c.Apply(`{"score":-5}`)
	if out != `{"score":0}` {
		t.Errorf("expected score clamped to 0, got %s", out)
	}
}

func TestApply_ClampsAbove(t *testing.T) {
	c, _ := New([]string{"score:0:100"})
	out := c.Apply(`{"score":200}`)
	if out != `{"score":100}` {
		t.Errorf("expected score clamped to 100, got %s", out)
	}
}

func TestApply_WithinRange_Unchanged(t *testing.T) {
	c, _ := New([]string{"score:0:100"})
	out := c.Apply(`{"score":50}`)
	if out != `{"score":50}` {
		t.Errorf("expected score unchanged, got %s", out)
	}
}

func TestApply_NonJSONPassThrough(t *testing.T) {
	c, _ := New([]string{"score:0:100"})
	input := "not json at all"
	out := c.Apply(input)
	if out != input {
		t.Errorf("expected passthrough, got %s", out)
	}
}

func TestApply_MissingField_PassThrough(t *testing.T) {
	c, _ := New([]string{"score:0:100"})
	input := `{"level":"info"}`
	out := c.Apply(input)
	if out != input {
		t.Errorf("expected passthrough for missing field, got %s", out)
	}
}

func TestApply_NonNumericField_PassThrough(t *testing.T) {
	c, _ := New([]string{"score:0:100"})
	input := `{"score":"high"}`
	out := c.Apply(input)
	if out != input {
		t.Errorf("expected passthrough for non-numeric field, got %s", out)
	}
}
