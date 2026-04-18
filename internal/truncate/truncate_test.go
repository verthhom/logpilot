package truncate_test

import (
	"strings"
	"testing"

	"github.com/yourorg/logpilot/internal/truncate"
)

func TestNew_Valid(t *testing.T) {
	tr, err := truncate.New(80)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tr == nil {
		t.Fatal("expected non-nil Truncator")
	}
}

func TestNew_TooSmall(t *testing.T) {
	_, err := truncate.New(2)
	if err == nil {
		t.Fatal("expected error for maxLen <= 3")
	}
}

func TestNew_ExactlyEllipsis(t *testing.T) {
	_, err := truncate.New(3)
	if err == nil {
		t.Fatal("expected error for maxLen == len(ellipsis)")
	}
}

func TestApply_ShortLine(t *testing.T) {
	tr, _ := truncate.New(20)
	line := "hello world"
	if got := tr.Apply(line); got != line {
		t.Errorf("expected %q, got %q", line, got)
	}
}

func TestApply_ExactLength(t *testing.T) {
	tr, _ := truncate.New(10)
	line := "1234567890"
	if got := tr.Apply(line); got != line {
		t.Errorf("expected unchanged line, got %q", got)
	}
}

func TestApply_LongLine(t *testing.T) {
	maxLen := 20
	tr, _ := truncate.New(maxLen)
	line := strings.Repeat("a", 100)
	got := tr.Apply(line)
	if len(got) != maxLen {
		t.Errorf("expected length %d, got %d", maxLen, len(got))
	}
	if !strings.HasSuffix(got, "...") {
		t.Errorf("expected ellipsis suffix, got %q", got)
	}
}

func TestEnabled(t *testing.T) {
	tr, _ := truncate.New(40)
	if !tr.Enabled() {
		t.Error("expected Enabled() == true")
	}
}
