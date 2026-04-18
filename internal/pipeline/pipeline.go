// Package pipeline wires sources, filters, and outputs together.
package pipeline

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/user/logpilot/internal/filter"
	"github.com/user/logpilot/internal/output"
	"github.com/user/logpilot/internal/source"
)

// Pipeline reads from multiple sources, applies a filter, and writes to an output.
type Pipeline struct {
	sources []source.Source
	filter  *filter.Filter
	output  *output.Output
}

// New creates a Pipeline.
func New(sources []source.Source, f *filter.Filter, out *output.Output) *Pipeline {
	return &Pipeline{sources: sources, filter: f, output: out}
}

// Run starts tailing all sources concurrently until ctx is cancelled.
func (p *Pipeline) Run(ctx context.Context) error {
	lines := make(chan string, 256)
	var wg sync.WaitGroup

	for _, src := range p.sources {
		wg.Add(1)
		go func(s source.Source) {
			defer wg.Done()
			_ = s.Tail(ctx, lines)
		}(src)
	}

	go func() {
		wg.Wait()
		close(lines)
	}()

	for line := range lines {
		if p.filter != nil {
			var record map[string]interface{}
			if err := json.Unmarshal([]byte(line), &record); err != nil {
				continue
			}
			if !p.filter.Match(record) {
				continue
			}
		}
		if err := p.output.Write(line); err != nil {
			return err
		}
	}
	return nil
}
