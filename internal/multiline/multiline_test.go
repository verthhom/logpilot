package multiline

import (
	"testing"
)

func TestNew_EmptyPattern(t *testing.T) {
	_, err := New("", " ")
	if err != ErrEmptyPattern {
		t.Fatalf("expected ErrEmptyPattern, got %v", err)
	}
}

func TestNew_BadPattern(t *testing.T) {
	_, err := New("[invalid", " ")
	if err != ErrBadPattern {
		t.Fatalf("expected ErrBadPattern, got %v", err)
	}
}

func TestNew_Valid(t *testing.T) {
	j, err := New(`^\d`, " ")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if j == nil {
		t.Fatal("expected non-nil Joiner")
	}
}

func TestFeed_BuffersUntilNewStart(t *testing.T) {
	j, _ := New(`^START`, " ")

	_, ok := j.Feed("START first")
	if ok {
		t.Fatal("should not emit on first start line")
	}
	j.Feed("continuation a")
	j.Feed("continuation b")

	event, ok := j.Feed("START second")
	if !ok {
		t.Fatal("expected emission when new start arrived")
	}
	expected := "START first continuation a continuation b"
	if event != expected {
		t.Fatalf("got %q, want %q", event, expected)
	}
}

func TestFlush_ReturnsPending(t *testing.T) {
	j, _ := New(`^START`, " ")
	j.Feed("START only")
	j.Feed("tail")

	event, ok := j.Flush()
	if !ok {
		t.Fatal("expected Flush to return buffered event")
	}
	if event != "START only tail" {
		t.Fatalf("unexpected event: %q", event)
	}
}

func TestFlush_EmptyBuffer(t *testing.T) {
	j, _ := New(`^START`, " ")
	_, ok := j.Flush()
	if ok {
		t.Fatal("expected false for empty buffer")
	}
}

func TestFeed_DefaultSep(t *testing.T) {
	j, err := New(`^A`, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	j.Feed("A line")
	j.Feed("extra")
	event, ok := j.Feed("A next")
	if !ok {
		t.Fatal("expected emission")
	}
	if event != "A line extra" {
		t.Fatalf("unexpected: %q", event)
	}
}
