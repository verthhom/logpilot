package cli

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func TestParse_Defaults(t *testing.T) {
	cfg, _, err := parse([]string{}, &bytes.Buffer{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Format != "pretty" {
		t.Errorf("expected format=pretty, got %s", cfg.Format)
	}
	if len(cfg.Filters) != 0 {
		t.Errorf("expected no filters, got %v", cfg.Filters)
	}
}

func TestParse_WithFlags(t *testing.T) {
	args := []string{"-format", "json", "-filter", "level=error,msg~fail", "app.log"}
	cfg, _, err := parse(args, &bytes.Buffer{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Format != "json" {
		t.Errorf("expected json, got %s", cfg.Format)
	}
	if len(cfg.Filters) != 2 {
		t.Errorf("expected 2 filters, got %d", len(cfg.Filters))
	}
	if len(cfg.Files) != 1 || cfg.Files[0] != "app.log" {
		t.Errorf("unexpected files: %v", cfg.Files)
	}
}

func TestParse_InvalidFlag(t *testing.T) {
	_, _, err := parse([]string{"-unknown"}, &bytes.Buffer{})
	if err == nil {
		t.Fatal("expected error for unknown flag")
	}
}

func TestExecute_InvalidFormat(t *testing.T) {
	cfg := &Config{Format: "xml", Filters: nil, Files: []string{}}
	err := execute(context.Background(), cfg, &bytes.Buffer{})
	if err == nil || !strings.Contains(err.Error(), "output") {
		t.Errorf("expected output error, got %v", err)
	}
}

func TestExecute_NilConfig(t *testing.T) {
	if err := execute(context.Background(), nil, &bytes.Buffer{}); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
