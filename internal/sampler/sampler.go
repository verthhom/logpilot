// Package sampler provides probabilistic log line sampling.
package sampler

import (
	"fmt"
	"math/rand"
)

// Sampler drops log lines probabilistically, keeping approximately
// 1-in-N lines based on the configured rate.
type Sampler struct {
	rate int
	rng  *rand.Rand
}

// New creates a Sampler that retains roughly 1/rate of all lines.
// rate must be >= 1; a rate of 1 keeps every line.
func New(rate int, seed int64) (*Sampler, error) {
	if rate < 1 {
		return nil, fmt.Errorf("sampler: rate must be >= 1, got %d", rate)
	}
	return &Sampler{
		rate: rate,
		rng:  rand.New(rand.NewSource(seed)),
	}, nil
}

// Keep returns true if the line should be passed through.
func (s *Sampler) Keep() bool {
	if s.rate == 1 {
		return true
	}
	return s.rng.Intn(s.rate) == 0
}

// Rate returns the configured sampling rate.
func (s *Sampler) Rate() int {
	return s.rate
}
