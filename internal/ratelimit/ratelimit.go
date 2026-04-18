// Package ratelimit provides a simple token-bucket rate limiter for
// controlling how many log lines are processed per second.
package ratelimit

import (
	"context"
	"errors"
	"time"
)

// ErrInvalidRate is returned when a non-positive rate is provided.
var ErrInvalidRate = errors.New("ratelimit: rate must be greater than zero")

// Limiter controls the flow of log lines using a token-bucket approach.
type Limiter struct {
	ticker *time.Ticker
	done   chan struct{}
}

// New creates a Limiter that allows up to linesPerSec lines per second.
// Returns ErrInvalidRate if linesPerSec is less than 1.
func New(linesPerSec int) (*Limiter, error) {
	if linesPerSec <= 0 {
		return nil, ErrInvalidRate
	}
	interval := time.Second / time.Duration(linesPerSec)
	return &Limiter{
		ticker: time.NewTicker(interval),
		done:   make(chan struct{}),
	}, nil
}

// Wait blocks until the limiter grants a token or the context is cancelled.
// Returns ctx.Err() if the context expires before a token is available.
func (l *Limiter) Wait(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-l.ticker.C:
		return nil
	}
}

// Stop releases resources held by the Limiter.
func (l *Limiter) Stop() {
	l.ticker.Stop()
}
