package levelfilter_test

import (
	"testing"

	"github.com/yourorg/logpilot/internal/levelfilter"
)

func TestNew_Valid(t *testing.T) {
	f, err := levelfilter.New("warn", "level")
	if err != nil || f == nil {
		t.Fatalf("expected valid filter, got err=%v", err)
	}
}

func TestNew_EmptyField(t *testing.T) {
	_, err := levelfilter.New("info", "")
	if err != levelfilter.ErrEmptyField {
		t.Fatalf("expected ErrEmptyField, got %v", err)
	}
}

func TestNew_UnknownLevel(t *testing.T) {
	_, err := levelfilter.New("verbose", "level")
	if err != levelfilter.ErrUnknownLevel {
		t.Fatalf("expected ErrUnknownLevel, got %v", err)
	}
}

func TestAllow_BelowMin(t *testing.T) {
	f, _ := levelfilter.New("warn", "level")
	if f.Allow(`{"level":"debug","msg":"hi"}`) {
		t.Fatal("expected debug to be dropped when min=warn")
	}
}

func TestAllow_AtMin(t *testing.T) {
	f, _ := levelfilter.New("warn", "level")
	if !f.Allow(`{"level":"warn","msg":"hi"}`) {
		t.Fatal("expected warn to pass when min=warn")
	}
}

func TestAllow_AboveMin(t *testing.T) {
	f, _ := levelfilter.New("warn", "level")
	if !f.Allow(`{"level":"error","msg":"boom"}`) {
		t.Fatal("expected error to pass when min=warn")
	}
}

func TestAllow_NonJSON(t *testing.T) {
	f, _ := levelfilter.New("error", "level")
	if !f.Allow("plain text line") {
		t.Fatal("expected non-JSON to pass through")
	}
}

func TestAllow_MissingField(t *testing.T) {
	f, _ := levelfilter.New("error", "level")
	if !f.Allow(`{"msg":"no level field"}`) {
		t.Fatal("expected missing field to pass through")
	}
}

func TestAllow_CaseInsensitive(t *testing.T) {
	f, _ := levelfilter.New("info", "level")
	if !f.Allow(`{"level":"INFO","msg":"hi"}`) {
		t.Fatal("expected INFO to match info case-insensitively")
	}
}
