// Package throttle provides a line-level throughput limiter that drops
// lines exceeding a configured burst window, emitting a summary count.
package throttle

import (
	"fmt"
	"sync"
	"time"
)

// Throttle drops lines that exceed MaxLines within Window.
type Throttle struct {
	maxLines int
	window   time.Duration
	mu       sync.Mutex
	count    int
	dropped  int
	windowStart time.Time
}

// New creates a Throttle allowing at most maxLines per window duration.
func New(maxLines int, window time.Duration) (*Throttle, error) {
	if maxLines <= 0 {
		return nil, ErrInvalidMaxLines
	}
	if window <= 0 {
		return nil, ErrInvalidWindow
	}
	return &Throttle{
		maxLines:    maxLines,
		window:      window,
		windowStart: time.Now(),
	}, nil
}

// Allow returns true if the line should be passed through, false if dropped.
func (t *Throttle) Allow() bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	now := time.Now()
	if now.Sub(t.windowStart) >= t.window {
		t.count = 0
		t.dropped = 0
		t.windowStart = now
	}
	if t.count < t.maxLines {
		t.count++
		return true
	}
	t.dropped++
	return false
}

// Summary returns a human-readable drop summary if any lines were dropped,
// otherwise returns an empty string.
func (t *Throttle) Summary() string {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.dropped == 0 {
		return ""
	}
	return fmt.Sprintf("[throttle] dropped %d lines in last window", t.dropped)
}
