package metrics

import "sync/atomic"

// Collector tracks runtime counters for a logpilot session.
type Collector struct {
	linesRead     atomic.Int64
	linesMatched  atomic.Int64
	linesDropped  atomic.Int64
	parseErrors   atomic.Int64
}

// New returns a new Collector.
func New() *Collector {
	return &Collector{}
}

// IncRead increments the lines-read counter.
func (c *Collector) IncRead() { c.linesRead.Add(1) }

// IncMatched increments the lines-matched counter.
func (c *Collector) IncMatched() { c.linesMatched.Add(1) }

// IncDropped increments the lines-dropped counter.
func (c *Collector) IncDropped() { c.linesDropped.Add(1) }

// IncParseError increments the parse-error counter.
func (c *Collector) IncParseError() { c.parseErrors.Add(1) }

// Snapshot returns a point-in-time copy of all counters.
func (c *Collector) Snapshot() Snapshot {
	return Snapshot{
		LinesRead:    c.linesRead.Load(),
		LinesMatched: c.linesMatched.Load(),
		LinesDropped: c.linesDropped.Load(),
		ParseErrors:  c.parseErrors.Load(),
	}
}

// Snapshot holds a point-in-time view of collected metrics.
type Snapshot struct {
	LinesRead    int64
	LinesMatched int64
	LinesDropped int64
	ParseErrors  int64
}
