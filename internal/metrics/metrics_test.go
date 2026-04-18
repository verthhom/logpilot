package metrics

import (
	"testing"
)

func TestNew(t *testing.T) {
	c := New()
	if c == nil {
		t.Fatal("expected non-nil Collector")
	}
}

func TestCounters_InitiallyZero(t *testing.T) {
	c := New()
	s := c.Snapshot()
	if s.LinesRead != 0 || s.LinesMatched != 0 || s.LinesDropped != 0 || s.ParseErrors != 0 {
		t.Errorf("expected all zeros, got %+v", s)
	}
}

func TestIncRead(t *testing.T) {
	c := New()
	c.IncRead()
	c.IncRead()
	if got := c.Snapshot().LinesRead; got != 2 {
		t.Errorf("LinesRead = %d, want 2", got)
	}
}

func TestIncMatched(t *testing.T) {
	c := New()
	c.IncMatched()
	if got := c.Snapshot().LinesMatched; got != 1 {
		t.Errorf("LinesMatched = %d, want 1", got)
	}
}

func TestIncDropped(t *testing.T) {
	c := New()
	c.IncDropped()
	c.IncDropped()
	c.IncDropped()
	if got := c.Snapshot().LinesDropped; got != 3 {
		t.Errorf("LinesDropped = %d, want 3", got)
	}
}

func TestIncParseError(t *testing.T) {
	c := New()
	c.IncParseError()
	if got := c.Snapshot().ParseErrors; got != 1 {
		t.Errorf("ParseErrors = %d, want 1", got)
	}
}

func TestSnapshot_Independence(t *testing.T) {
	c := New()
	s1 := c.Snapshot()
	c.IncRead()
	s2 := c.Snapshot()
	if s1.LinesRead != 0 {
		t.Errorf("s1.LinesRead should still be 0, got %d", s1.LinesRead)
	}
	if s2.LinesRead != 1 {
		t.Errorf("s2.LinesRead = %d, want 1", s2.LinesRead)
	}
}
