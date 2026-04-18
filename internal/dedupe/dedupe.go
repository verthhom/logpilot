package dedupe

import (
	"crypto/md5"
	"fmt"
	"sync"
)

// Deduplicator drops repeated log lines within a sliding window of recent hashes.
type Deduplicator struct {
	mu      sync.Mutex
	seen    map[string]struct{}
	order   []string
	window  int
}

// New creates a Deduplicator that remembers the last windowSize unique lines.
// windowSize must be >= 1.
func New(windowSize int) (*Deduplicator, error) {
	if windowSize < 1 {
		return nil, ErrInvalidWindow
	}
	return &Deduplicator{
		seen:   make(map[string]struct{}, windowSize),
		order:  make([]string, 0, windowSize),
		window: windowSize,
	}, nil
}

// IsDuplicate returns true if line was seen within the current window.
// If not a duplicate, the line is recorded and older entries evicted.
func (d *Deduplicator) IsDuplicate(line string) bool {
	h := hash(line)
	d.mu.Lock()
	defer d.mu.Unlock()

	if _, ok := d.seen[h]; ok {
		return true
	}

	d.seen[h] = struct{}{}
	d.order = append(d.order, h)

	if len(d.order) > d.window {
		evict := d.order[0]
		d.order = d.order[1:]
		delete(d.seen, evict)
	}
	return false
}

// Reset clears all recorded hashes.
func (d *Deduplicator) Reset() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.seen = make(map[string]struct{}, d.window)
	d.order = d.order[:0]
}

func hash(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}
