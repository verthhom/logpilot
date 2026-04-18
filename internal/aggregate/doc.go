// Package aggregate provides a time-windowed log aggregator that groups
// structured JSON log lines by a configurable field and emits count summaries
// at the end of each window.
//
// Usage:
//
//	agg, err := aggregate.New("level", 10*time.Second)
//	if err != nil { ... }
//	defer agg.Stop()
//
//	for _, line := range lines {
//		agg.Feed(line)
//	}
//
//	for summary := range agg.Out() {
//		fmt.Println(summary)
//	}
package aggregate
