package replay

import (
	"bufio"
	"context"
	"os"
	"time"
)

// Replayer replays a log file at a controlled rate.
type Replayer struct {
	path  string
	delay time.Duration
}

// New creates a new Replayer for the given file path.
// delay controls how long to wait between emitting each line.
func New(path string, delay time.Duration) (*Replayer, error) {
	if path == "" {
		return nil, ErrEmptyPath
	}
	if delay < 0 {
		return nil, ErrNegativeDelay
	}
	return &Replayer{path: path, delay: delay}, nil
}

// Run opens the file and sends each line to out, respecting ctx cancellation.
// It closes out when done or when ctx is cancelled.
func (r *Replayer) Run(ctx context.Context, out chan<- string) error {
	defer close(out)

	f, err := os.Open(r.path)
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		line := scanner.Text()
		if line == "" {
			continue
		}

		select {
		case out <- line:
		case <-ctx.Done():
			return ctx.Err()
		}

		if r.delay > 0 {
			select {
			case <-time.After(r.delay):
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}
	return scanner.Err()
}
