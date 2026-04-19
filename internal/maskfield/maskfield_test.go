package maskfield

import (
	"strings"
	"testing"
)

func TestNew_Valid(t *testing.T) {
	_, err := New([]string{"password"}, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNew_NoFields(t *testing.T) {
	_, err := New([]string{}, 2)
	if err != ErrNoFields {
		t.Fatalf("expected ErrNoFields, got %v", err)
	}
}

func TestNew_BlankField(t *testing.T) {
	_, err := New([]string{"  "}, 2)
	if err != ErrBlankField {
		t.Fatalf("expected ErrBlankField, got %v", err)
	}
}

func TestNew_NegativeVisible(t *testing.T) {
	_, err := New([]string{"token"}, -1)
	if err != ErrNegativeVisible {
		t.Fatalf("expected ErrNegativeVisible, got %v", err)
	}
}

func TestApply_MasksField(t *testing.T) {
	m, _ := New([]string{"password"}, 2)
	out := m.Apply(`{"password":"secret123","user":"alice"}`)
	if strings.Contains(out, "secret123") {
		t.Errorf("expected value to be masked, got: %s", out)
	}
	if !strings.Contains(out, "se*******") {
		t.Errorf("expected masked prefix 'se*******', got: %s", out)
	}
	if !strings.Contains(out, "alice") {
		t.Errorf("expected unmasked user field, got: %s", out)
	}
}

func TestApply_NonJSONPassThrough(t *testing.T) {
	m, _ := New([]string{"password"}, 2)
	line := "not json at all"
	if got := m.Apply(line); got != line {
		t.Errorf("expected passthrough, got: %s", got)
	}
}

func TestApply_VisibleExceedsLength(t *testing.T) {
	m, _ := New([]string{"token"}, 100)
	out := m.Apply(`{"token":"abc"}`)
	if !strings.Contains(out, "abc") {
		t.Errorf("expected full value when visible >= length, got: %s", out)
	}
}

func TestApply_ZeroVisible(t *testing.T) {
	m, _ := New([]string{"secret"}, 0)
	out := m.Apply(`{"secret":"hello"}`)
	if strings.Contains(out, "hello") {
		t.Errorf("expected fully masked value, got: %s", out)
	}
	if !strings.Contains(out, "*****") {
		t.Errorf("expected all asterisks, got: %s", out)
	}
}
