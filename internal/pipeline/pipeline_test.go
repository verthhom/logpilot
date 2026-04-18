package pipeline_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/user/logpilot/internal/filter"
	"github.com/user/logpilot/internal/output"
	"github.com/user/logpilot/internal/pipeline"
	"github.com/user/logpilot/internal/source"
)

func TestPipeline_Run_NoFilter(t *testing.T) {
	src := source.NewStdin(strings.NewReader(`{"level":"info","msg":"hello"}
`))
	var buf strings.Builder
	out, _ := output.New("json", &buf)
	p := pipeline.New([]source.Source{src}, nil, out)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := p.Run(ctx); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "hello") {
		t.Errorf("expected output to contain 'hello', got: %s", buf.String())
	}
}

func TestPipeline_Run_WithFilter_Matches(t *testing.T) {
	src := source.NewStdin(strings.NewReader(`{"level":"error","msg":"boom"}
{"level":"info","msg":"ok"}
`))
	var buf strings.Builder
	out, _ := output.New("json", &buf)
	f, _ := filter.New([]string{"level=error"})
	p := pipeline.New([]source.Source{src}, f, out)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_ = p.Run(ctx)

	if !strings.Contains(buf.String(), "boom") {
		t.Errorf("expected 'boom' in output")
	}
	if strings.Contains(buf.String(), "ok") {
		t.Errorf("did not expect 'ok' in filtered output")
	}
}

func TestPipeline_Run_ContextCancel(t *testing.T) {
	src := source.NewStdin(strings.NewReader(""))
	var buf strings.Builder
	out, _ := output.New("json", &buf)
	p := pipeline.New([]source.Source{src}, nil, out)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if err := p.Run(ctx); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
