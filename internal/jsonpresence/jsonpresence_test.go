package jsonpresence

import (
	"testing"
)

func TestNew_Valid(t *testing.T) {
	_, err := New([]string{"level"}, []string{"debug"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNew_OnlyRequired(t *testing.T) {
	_, err := New([]string{"msg"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNew_OnlyForbidden(t *testing.T) {
	_, err := New(nil, []string{"trace"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNew_NoFields(t *testing.T) {
	_, err := New(nil, nil)
	if err == nil {
		t.Fatal("expected error for empty field lists")
	}
}

func TestNew_BlankRequired(t *testing.T) {
	_, err := New([]string{""}, nil)
	if err == nil {
		t.Fatal("expected error for blank required field")
	}
}

func TestNew_BlankForbidden(t *testing.T) {
	_, err := New(nil, []string{""})
	if err == nil {
		t.Fatal("expected error for blank forbidden field")
	}
}

func TestAllow_RequiredPresent(t *testing.T) {
	c, _ := New([]string{"level", "msg"}, nil)
	if !c.Allow(`{"level":"info","msg":"hello"}`) {
		t.Fatal("expected Allow to return true")
	}
}

func TestAllow_RequiredMissing(t *testing.T) {
	c, _ := New([]string{"level", "msg"}, nil)
	if c.Allow(`{"level":"info"}`) {
		t.Fatal("expected Allow to return false when required field missing")
	}
}

func TestAllow_ForbiddenAbsent(t *testing.T) {
	c, _ := New(nil, []string{"debug"})
	if !c.Allow(`{"level":"info","msg":"ok"}`) {
		t.Fatal("expected Allow to return true when forbidden field absent")
	}
}

func TestAllow_ForbiddenPresent(t *testing.T) {
	c, _ := New(nil, []string{"debug"})
	if c.Allow(`{"level":"info","debug":true}`) {
		t.Fatal("expected Allow to return false when forbidden field present")
	}
}

func TestAllow_NonJSONPassThrough(t *testing.T) {
	c, _ := New([]string{"level"}, nil)
	if c.Allow("not json at all") {
		t.Fatal("expected Allow to return false for non-JSON input")
	}
}

func TestAllow_BothConstraints(t *testing.T) {
	c, _ := New([]string{"msg"}, []string{"error"})
	if !c.Allow(`{"msg":"hi","level":"info"}`) {
		t.Fatal("expected true: required present, forbidden absent")
	}
	if c.Allow(`{"msg":"hi","error":"oops"}`) {
		t.Fatal("expected false: forbidden field present")
	}
	if c.Allow(`{"level":"info"}`) {
		t.Fatal("expected false: required field missing")
	}
}
