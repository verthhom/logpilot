package jsontypecheck_test

import (
	"testing"

	"github.com/yourorg/logpilot/internal/jsontypecheck"
)

func TestNew_Valid(t *testing.T) {
	c, err := jsontypecheck.New([]string{"level:string", "count:number"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil Checker")
	}
}

func TestNew_NoRules(t *testing.T) {
	_, err := jsontypecheck.New(nil)
	if err == nil {
		t.Fatal("expected error for empty rules")
	}
}

func TestNew_BadRule(t *testing.T) {
	_, err := jsontypecheck.New([]string{"noseparator"})
	if err == nil {
		t.Fatal("expected error for rule without colon")
	}
}

func TestNew_BlankPart(t *testing.T) {
	for _, r := range []string{":string", "field:"} {
		_, err := jsontypecheck.New([]string{r})
		if err == nil {
			t.Fatalf("expected error for rule %q", r)
		}
	}
}

func TestNew_UnknownType(t *testing.T) {
	_, err := jsontypecheck.New([]string{"field:integer"})
	if err == nil {
		t.Fatal("expected error for unknown type")
	}
}

func TestCheck_CorrectType_PassThrough(t *testing.T) {
	c, _ := jsontypecheck.New([]string{"level:string", "retries:number"})
	line := `{"level":"info","retries":3}`
	if got := c.Check(line); got != line {
		t.Fatalf("expected pass-through, got %q", got)
	}
}

func TestCheck_WrongType_Dropped(t *testing.T) {
	c, _ := jsontypecheck.New([]string{"level:string"})
	line := `{"level":42}`
	if got := c.Check(line); got != "" {
		t.Fatalf("expected empty string for type mismatch, got %q", got)
	}
}

func TestCheck_AbsentField_PassThrough(t *testing.T) {
	c, _ := jsontypecheck.New([]string{"level:string"})
	line := `{"msg":"hello"}`
	if got := c.Check(line); got != line {
		t.Fatalf("expected pass-through for absent field, got %q", got)
	}
}

func TestCheck_NonJSON_PassThrough(t *testing.T) {
	c, _ := jsontypecheck.New([]string{"level:string"})
	line := "plain text line"
	if got := c.Check(line); got != line {
		t.Fatalf("expected pass-through for non-JSON, got %q", got)
	}
}

func TestCheck_BoolField(t *testing.T) {
	c, _ := jsontypecheck.New([]string{"ok:bool"})
	good := `{"ok":true}`
	if got := c.Check(good); got != good {
		t.Fatalf("expected pass-through, got %q", got)
	}
	bad := `{"ok":"yes"}`
	if got := c.Check(bad); got != "" {
		t.Fatalf("expected drop, got %q", got)
	}
}

func TestCheck_ArrayField(t *testing.T) {
	c, _ := jsontypecheck.New([]string{"tags:array"})
	good := `{"tags":["a","b"]}`
	if got := c.Check(good); got != good {
		t.Fatalf("expected pass-through, got %q", got)
	}
}
