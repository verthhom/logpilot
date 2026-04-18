// Package replay provides functionality to replay structured JSON log files
// at a controlled rate, simulating a live log stream from a static file.
//
// # Overview
//
// A Replayer reads a file line by line and emits each line to a channel with
// an optional delay between lines. This is useful for testing pipelines or
// demonstrating logpilot against a known dataset without a live source.
//
// # Usage
//
//	r, err := replay.New("app.log", 100*time.Millisecond)
//	if err != nil { ... }
//
//	ch := make(chan string, 16)
//	go r.Run(ctx, ch)
//	for line := range ch {
//		fmt.Println(line)
//	}
package replay
