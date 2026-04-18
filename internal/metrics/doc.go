// Package metrics provides lightweight, thread-safe counters for tracking
// logpilot session statistics.
//
// A Collector accumulates counts for lines read, matched, dropped, and parse
// errors encountered during a pipeline run. Counters are updated via atomic
// operations so they are safe to call from concurrent goroutines.
//
// Usage:
//
//	c := metrics.New()
//	c.IncRead()
//	c.IncMatched()
//	snap := c.Snapshot()
//	fmt.Printf("read=%d matched=%d\n", snap.LinesRead, snap.LinesMatched)
package metrics
