package aggregate

import (
	"encoding/json"
	"sync"
	"time"
)

// Aggregator groups log lines by a key field and counts occurrences within a
// rolling time window, emitting a summary line when the window closes.
type Aggregator struct {
	field  string
	window time.Duration
	mu     sync.Mutex
	bucket map[string]int
	timer  *time.Timer
	out    chan string
}

// New creates an Aggregator that groups by field over window duration.
func New(field string, window time.Duration) (*Aggregator, error) {
	if field == "" {
		return nil, ErrEmptyField
	}
	if window <= 0 {
		return nil, ErrInvalidWindow
	}
	a := &Aggregator{
		field:  field,
		window: window,
		bucket: make(map[string]int),
		out:    make(chan string, 64),
	}
	a.timer = time.AfterFunc(window, a.flush)
	return a, nil
}

// Feed accepts a raw log line and accumulates counts.
func (a *Aggregator) Feed(line string) {
	var m map[string]any
	if err := json.Unmarshal([]byte(line), &m); err != nil {
		return
	}
	val, ok := m[a.field]
	if !ok {
		return
	}
	key, ok := val.(string)
	if !ok {
		return
	}
	a.mu.Lock()
	a.bucket[key]++
	a.mu.Unlock()
}

// Out returns the channel on which summary lines are emitted.
func (a *Aggregator) Out() <-chan string { return a.out }

// Stop cancels the internal timer and flushes remaining counts.
func (a *Aggregator) Stop() {
	a.timer.Stop()
	a.flush()
	close(a.out)
}

func (a *Aggregator) flush() {
	a.mu.Lock()
	snap := a.bucket
	a.bucket = make(map[string]int)
	a.mu.Unlock()
	if len(snap) == 0 {
		return
	}
	summary := map[string]any{
		"type":    "aggregate",
		"field":   a.field,
		"window":  a.window.String(),
		"counts":  snap,
		"time":    time.Now().UTC().Format(time.RFC3339),
	}
	b, _ := json.Marshal(summary)
	a.out <- string(b)
	a.timer.Reset(a.window)
}
