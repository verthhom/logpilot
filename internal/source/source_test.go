package source_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/logpilot/internal/source"
)

func writeTempFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "test.log")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("writeTempFile: %v", err)
	}
	return p
}

func TestFileSource_Name(t *testing.T) {
	fs := source.NewFile("/var/log/app.log")
	if fs.Name() != "/var/log/app.log" {
		t.Errorf("expected /var/log/app.log, got %s", fs.Name())
	}
}

func TestFileSource_Tail_ReadsAllLines(t *testing.T) {
	path := writeTempFile(t, `{"level":"info","msg":"start"}
{"level":"error","msg":"boom"}
`)

	fs := source.NewFile(path)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	ch, err := fs.Tail(ctx)
	if err != nil {
		t.Fatalf("Tail error: %v", err)
	}

	var lines []source.Line
	for l := range ch {
		lines = append(lines, l)
	}

	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
	if lines[0].Source != path {
		t.Errorf("unexpected source: %s", lines[0].Source)
	}
}

func TestFileSource_Tail_ContextCancel(t *testing.T) {
	path := writeTempFile(t, "{\"a\":1}\n{\"b\":2}\n{\"c\":3}\n")

	fs := source.NewFile(path)
	ctx, cancel := context.WithCancel(context.Background())

	ch, err := fs.Tail(ctx)
	if err != nil {
		t.Fatalf("Tail error: %v", err)
	}

	// Read one line then cancel.
	<-ch
	cancel()

	// Drain remaining; channel must close eventually.
	timeout := time.After(2 * time.Second)
	for {
		select {
		case _, ok := <-ch:
			if !ok {
				return
			}
		case <-timeout:
			t.Fatal("channel did not close after context cancel")
		}
	}
}

func TestFileSource_Tail_FileNotFound(t *testing.T) {
	fs := source.NewFile("/nonexistent/path/app.log")
	ctx := context.Background()
	_, err := fs.Tail(ctx)
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestStdinSource_Name(t *testing.T) {
	s := source.NewStdin()
	if s.Name() != "stdin" {
		t.Errorf("expected stdin, got %s", s.Name())
	}
}
