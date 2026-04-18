// Package snapshot provides a lightweight mechanism for capturing and
// persisting point-in-time pipeline metric summaries.
//
// A Store holds the most recent Snapshot in memory and atomically writes
// it to a JSON file on disk so that external tooling (dashboards, health
// checks) can read pipeline state without coupling to the running process.
//
// Typical usage:
//
//	s, err := snapshot.New("/var/run/logpilot/state.json")
//	if err != nil { ... }
//	s.Save(snapshot.Snapshot{
//		Source:  "app.log",
//		Read:    m.Read(),
//		Matched: m.Matched(),
//		Dropped: m.Dropped(),
//	})
package snapshot
