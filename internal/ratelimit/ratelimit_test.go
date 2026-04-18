package ratelimit_test

import (
	"context"
	"testing"
	"time"

	"github.com/yourorg/logpilot/internal/ratelimit"
)

func TestNew_ValidRate(t *testing.T) {
	l, err := ratelimit.New(10)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	l.Stop()
}

func TestNew_InvalidRate(t *testing.T) {
	for _, rate := range []int{0, -1, -100} {
		_, err := ratelimit.New(rate)
		if err == nil {
			t.Errorf("expected error for rate %d, got nil", rate)
		}
	}
}

func TestWait_ContextCancel(t *testing.T) {
	// Rate of 1 line/sec means next tick is ~1s away; cancel immediately.
	l, err := ratelimit.New(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer l.Stop()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err = l.Wait(ctx)
	if err == nil {
		t.Fatal("expected context error, got nil")
	}
}

func TestWait_TokenGranted(t *testing.T) {
	// High rate so the tick fires quickly.
	l, err := ratelimit.New(1000)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer l.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	if err := l.Wait(ctx); err != nil {
		t.Fatalf("expected token to be granted, got: %v", err)
	}
}

func TestWait_RespectsThroughput(t *testing.T) {
	const rate = 50
	l, err := ratelimit.New(rate)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer l.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	count := 0
	start := time.Now()
	for time.Since(start) < time.Second {
		if err := l.Wait(ctx); err != nil {
			break
		}
		count++
	}
	// Allow ±30% tolerance.
	if count < rate*7/10 || count > rate*13/10 {
		t.Errorf("expected ~%d lines/sec, got %d", rate, count)
	}
}
