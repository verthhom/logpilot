// Package source provides abstractions for reading structured JSON log lines
// from various input sources such as files, stdin, or named pipes.
package source

import (
	"bufio"
	"context"
	"io"
	"os"
)

// Line represents a single log line read from a source.
type Line struct {
	// Source is the name or path of the origin (e.g. filename or "stdin").
	Source string
	// Raw is the raw bytes of the log line.
	Raw []byte
}

// Source is the interface that wraps the basic Tail method.
type Source interface {
	// Tail emits log lines to the returned channel until ctx is cancelled or
	// the source is exhausted.
	Tail(ctx context.Context) (<-chan Line, error)
	// Name returns a human-readable identifier for the source.
	Name() string
}

// FileSource tails a single file.
type FileSource struct {
	path string
}

// NewFile creates a new FileSource for the given path.
func NewFile(path string) *FileSource {
	return &FileSource{path: path}
}

// Name returns the file path.
func (f *FileSource) Name() string { return f.path }

// Tail opens the file and streams lines until EOF or ctx cancellation.
func (f *FileSource) Tail(ctx context.Context) (<-chan Line, error) {
	file, err := os.Open(f.path)
	if err != nil {
		return nil, err
	}

	ch := make(chan Line)
	go func() {
		defer close(ch)
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			select {
			case <-ctx.Done():
				return
			case ch <- Line{Source: f.path, Raw: append([]byte(nil), scanner.Bytes()...)}:
			}
		}
	}()
	return ch, nil
}

// StdinSource reads log lines from os.Stdin.
type StdinSource struct{}

// NewStdin creates a new StdinSource.
func NewStdin() *StdinSource { return &StdinSource{} }

// Name returns "stdin".
func (s *StdinSource) Name() string { return "stdin" }

// Tail streams lines from stdin until EOF or ctx cancellation.
func (s *StdinSource) Tail(ctx context.Context) (<-chan Line, error) {
	return tailReader(ctx, "stdin", os.Stdin), nil
}

func tailReader(ctx context.Context, name string, r io.Reader) <-chan Line {
	ch := make(chan Line)
	go func() {
		defer close(ch)
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			select {
			case <-ctx.Done():
				return
			case ch <- Line{Source: name, Raw: append([]byte(nil), scanner.Bytes()...)}:
			}
		}
	}()
	return ch
}
