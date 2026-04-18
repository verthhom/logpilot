// Package cli parses command-line arguments and wires together the
// source, filter, output, and pipeline components.
package cli

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/user/logpilot/internal/filter"
	"github.com/user/logpilot/internal/output"
	"github.com/user/logpilot/internal/pipeline"
	"github.com/user/logpilot/internal/source"
)

// Config holds parsed CLI options.
type Config struct {
	Files   []string
	Format  string
	Filters []string
}

// Run parses args and executes the pipeline.
func Run(ctx context.Context, args []string) error {
	cfg, output_writer, err := parse(args, os.Stdout)
	if err != nil {
		return err
	}
	return execute(ctx, cfg, output_writer)
}

func parse(args []string, stdout io.Writer) (*Config, io.Writer, error) {
	fs := flag.NewFlagSet("logpilot", flag.ContinueOnError)
	fs.SetOutput(stdout)

	format := fs.String("format", "pretty", "output format: json|pretty")
	filterFlag := fs.String("filter", "", "comma-separated filter rules, e.g. level=error,msg~timeout")

	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return nil, stdout, nil
		}
		return nil, stdout, fmt.Errorf("parse flags: %w", err)
	}

	var filters []string
	if *filterFlag != "" {
		filters = strings.Split(*filterFlag, ",")
	}

	return &Config{
		Files:   fs.Args(),
		Format:  *format,
		Filters: filters,
	}, stdout, nil
}

func execute(ctx context.Context, cfg *Config, w io.Writer) error {
	if cfg == nil {
		return nil
	}

	out, err := output.New(cfg.Format, w)
	if err != nil {
		return fmt.Errorf("output: %w", err)
	}

	var f *filter.Filter
	if len(cfg.Filters) > 0 {
		f, err = filter.New(cfg.Filters)
		if err != nil {
			return fmt.Errorf("filter: %w", err)
		}
	}

	var src source.Source
	if len(cfg.Files) == 0 {
		src = source.NewStdin()
	} else {
		src, err = source.NewFile(cfg.Files[0])
		if err != nil {
			return fmt.Errorf("source: %w", err)
		}
	}

	p := pipeline.New(src, f, out)
	return p.Run(ctx)
}
