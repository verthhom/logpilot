package throttle_test

import (
	"testing"
	"time"

	"github.com/logpilot/internal/throttle"
)

func TestNew_Valid(t *testing.T) {
	th, err := throttle.New(10, time.Second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if th == nil {
		t.Fatal("expected non-nil Throttle")
	}
}

func TestNew_InvalidMaxLines(t *testing.T) {
	_, err := throttle.New(0, time.Second)
	if err != throttle.ErrInvalidMaxLines {
		t.Fatalf("expected ErrInvalidMaxLines, got %v", err)
	}
}

func TestNew_InvalidWindow(t *testing.T) {
	_, err := throttle.New(5, 0)
	if err != throttle.ErrInvalidWindow {
		t.Fatalf("expected ErrInvalidWindow, got %v", err)
	}
}

func TestAllow_WithinLimit(t *testing.T) {
	th, _ := throttle.New(3, time.Second)
	for i := 0; i < 3; i++ {
		if !th.Allow() {
			t.Fatalf("expected Allow()=true on call %d", i+1)
		}
	}
}

func TestAllow_ExceedsLimit(t *testing.T) {
	th, _ := throttle.New(2, time.Second)
	th.Allow()
	th.Allow()
	if th.Allow() {
		t.Fatal("expected Allow()=false after limit exceeded")
	}
}

func TestSummary_NoneDropped(t *testing.T) {
	th, _ := throttle.New(10, time.Second)
	th.Allow()
	if s := th.Summary(); s != "" {
		t.Fatalf("expected empty summary, got %q", s)
	}
}

func TestSummary_SomeDropped(t *testing.T) {
	th, _ := throttle.New(1, time.Second)
	th.Allow()
	th.Allow() // dropped
	s := th.Summary()
	if s == "" {
		t.Fatal("expected non-empty summary after drop")
	}
}

func TestAllow_WindowReset(t *testing.T) {
	th, _ := throttle.New(1, 50*time.Millisecond)
	th.Allow()
	if th.Allow() {
		t.Fatal("expected drop within window")
	}
	time.Sleep(60 * time.Millisecond)
	if !th.Allow() {
		t.Fatal("expected allow after window reset")
	}
}
