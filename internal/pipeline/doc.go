// Package pipeline provides the core processing pipeline for logpilot.
//
// It connects one or more [source.Source] instances to a single [output.Output],
// optionally applying a [filter.Filter] to each log line before writing.
//
// # Usage
//
//	sources := []source.Source{fileSource, stdinSource}
//	f, _ := filter.New([]string{"level=error"})
//	out, _ := output.New("pretty", os.Stdout)
//	p := pipeline.New(sources, f, out)
//	if err := p.Run(ctx); err != nil {
//		log.Fatal(err)
//	}
//
// All sources are tailed concurrently. Non-JSON lines are silently dropped
// when a filter is active. Cancelling the context stops all sources and
// causes Run to return.
package pipeline
