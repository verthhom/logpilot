package replay

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"
)

func writeTempLog(t *testing.T, lines []string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "replay-*.log")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	defer f.Close()
	_, _ = f.WriteString(strings.Join(lines, "\n") + "\n")
	return f.Name()
}

func TestNew_EmptyPath(t *testing.T) {
	_, err := New("", 0)
	if err != ErrEmptyPath {
		t.Fatalf("expected ErrEmptyPath, got %v", err)
	}
}

func TestNew_NegativeDelay(t *testing.T) {
	_, err := New("file.log", -1*time.Millisecond)
	if err != ErrNegativeDelay {
		t.Fatalf("expected ErrNegativeDelay, got %v", err)
	}
}

func TestRun_AllLinesReceived(t *testing.T) {
	lines := []string{`{"level":"info"}`, `{"level":"warn"}`, `{"level":"error"}`}
	path := writeTempLog(t, lines)

	r, err := New(path, 0)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	ch := make(chan string, 10)
	ctx := context.Background()
	if err := r.Run(ctx, ch); err != nil {
		t.Fatalf("Run: %v", err)
	}

	var got []string
	for line := range ch {
		got = append(got, line)
	}
	if len(got) != len(lines) {
		t.Fatalf("expected %d lines, got %d", len(lines), len(got))
	}
}

func TestRun_ContextCancel(t *testing.T) {
	lines := make([]string, 50)
	for i := range lines {
		lines[i] = `{"n":1}`
	}
	path := writeTempLog(t, lines)

	r, err := New(path, 20*time.Millisecond)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	ch := make(chan string, 5)

	done := make(chan error, 1)
	go func() { done <- r.Run(ctx, ch) }()

	time.Sleep(30 * time.Millisecond)
	cancel()

	err = <-done
	if err == nil {
		t.Fatal("expected context cancellation error, got nil")
	}
}

func TestRun_FileNotFound(t *testing.T) {
	r, _ := New("/nonexistent/path/file.log", 0)
	ch := make(chan string, 1)
	err := r.Run(context.Background(), ch)
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}
