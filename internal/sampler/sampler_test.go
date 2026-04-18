package sampler

import (
	"testing"
)

func TestNew_Valid(t *testing.T) {
	s, err := New(10, 42)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Rate() != 10 {
		t.Errorf("expected rate 10, got %d", s.Rate())
	}
}

func TestNew_RateOne(t *testing.T) {
	s, err := New(1, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i := 0; i < 20; i++ {
		if !s.Keep() {
			t.Error("rate=1 should keep every line")
		}
	}
}

func TestNew_InvalidRate(t *testing.T) {
	_, err := New(0, 0)
	if err == nil {
		t.Fatal("expected error for rate=0")
	}
	_, err = New(-5, 0)
	if err == nil {
		t.Fatal("expected error for negative rate")
	}
}

func TestKeep_ApproximateRate(t *testing.T) {
	const rate = 10
	const iterations = 10000
	s, _ := New(rate, 99)

	kept := 0
	for i := 0; i < iterations; i++ {
		if s.Keep() {
			kept++
		}
	}

	// Expect ~1000 kept; allow ±300 tolerance.
	expected := iterations / rate
	if kept < expected-300 || kept > expected+300 {
		t.Errorf("expected ~%d kept, got %d", expected, kept)
	}
}

func TestKeep_DifferentSeeds(t *testing.T) {
	s1, _ := New(2, 1)
	s2, _ := New(2, 2)

	matches := 0
	const n = 100
	for i := 0; i < n; i++ {
		if s1.Keep() == s2.Keep() {
			matches++
		}
	}
	// Different seeds should produce different sequences.
	if matches == n {
		t.Error("different seeds produced identical sequences")
	}
}
