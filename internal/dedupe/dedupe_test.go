package dedupe

import (
	"fmt"
	"testing"
)

func TestNew_Valid(t *testing.T) {
	d, err := New(10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d == nil {
		t.Fatal("expected non-nil deduplicator")
	}
}

func TestNew_InvalidWindow(t *testing.T) {
	_, err := New(0)
	if err != ErrInvalidWindow {
		t.Fatalf("expected ErrInvalidWindow, got %v", err)
	}
}

func TestIsDuplicate_FirstOccurrence(t *testing.T) {
	d, _ := New(5)
	if d.IsDuplicate("hello") {
		t.Error("first occurrence should not be duplicate")
	}
}

func TestIsDuplicate_SecondOccurrence(t *testing.T) {
	d, _ := New(5)
	d.IsDuplicate("hello")
	if !d.IsDuplicate("hello") {
		t.Error("second occurrence should be duplicate")
	}
}

func TestIsDuplicate_WindowEviction(t *testing.T) {
	d, _ := New(3)
	d.IsDuplicate("a") // enters window
	d.IsDuplicate("b")
	d.IsDuplicate("c")
	d.IsDuplicate("d") // evicts "a"

	// "a" should no longer be in window
	if d.IsDuplicate("a") {
		t.Error("'a' should have been evicted from window")
	}
}

func TestIsDuplicate_UniqueLines(t *testing.T) {
	d, _ := New(100)
	for i := 0; i < 50; i++ {
		line := fmt.Sprintf(`{"n":%d}`, i)
		if d.IsDuplicate(line) {
			t.Errorf("line %d should not be duplicate", i)
		}
	}
}

func TestReset_ClearsState(t *testing.T) {
	d, _ := New(5)
	d.IsDuplicate("x")
	d.Reset()
	if d.IsDuplicate("x") {
		t.Error("after reset 'x' should not be duplicate")
	}
}
