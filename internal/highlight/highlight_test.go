package highlight

import (
	"strings"
	"testing"
)

func TestNew_Enabled(t *testing.T) {
	h := New(true)
	if !h.enabled {
		t.Fatal("expected enabled=true")
	}
}

func TestNew_Disabled(t *testing.T) {
	h := New(false)
	if h.enabled {
		t.Fatal("expected enabled=false")
	}
}

func TestLevel_KnownLevel_Colored(t *testing.T) {
	h := New(true)
	result := h.Level("error")
	if !strings.Contains(result, Red) {
		t.Errorf("expected red color for error level, got %q", result)
	}
	if !strings.Contains(result, Reset) {
		t.Errorf("expected reset sequence, got %q", result)
	}
}

func TestLevel_UnknownLevel_PlainText(t *testing.T) {
	h := New(true)
	result := h.Level("trace")
	if result != "trace" {
		t.Errorf("expected plain text for unknown level, got %q", result)
	}
}

func TestLevel_Disabled(t *testing.T) {
	h := New(false)
	result := h.Level("error")
	if result != "error" {
		t.Errorf("expected plain text when disabled, got %q", result)
	}
}

func TestLevel_CaseInsensitive(t *testing.T) {
	h := New(true)
	result := h.Level("WARN")
	if !strings.Contains(result, Yellow) {
		t.Errorf("expected yellow for WARN, got %q", result)
	}
}

func TestLevel_AllKnownLevels(t *testing.T) {
	tests := []struct {
		level    string
		expColor string
	}{
		{"error", Red},
		{"warn", Yellow},
		{"info", Green},
		{"debug", Blue},
	}
	h := New(true)
	for _, tt := range tests {
		t.Run(tt.level, func(t *testing.T) {
			result := h.Level(tt.level)
			if !strings.Contains(result, tt.expColor) {
				t.Errorf("Level(%q): expected color %q, got %q", tt.level, tt.expColor, result)
			}
		})
	}
}

func TestKey_Enabled(t *testing.T) {
	h := New(true)
	result := h.Key("message")
	if !strings.Contains(result, Cyan) {
		t.Errorf("expected cyan for key, got %q", result)
	}
}

func TestValue_Enabled(t *testing.T) {
	h := New(true)
	result := h.Value("hello")
	if !strings.Contains(result, Bold) {
		t.Errorf("expected bold for value, got %q", result)
	}
}

func TestStrip_RemovesEscapes(t *testing.T) {
	input := "\033[31merror\033[0m"
	got := Strip(input)
	if got != "error" {
		t.Errorf("Strip() = %q, want %q", got, "error")
	}
}

func TestStrip_PlainText_Unchanged(t *testing.T) {
	input := "plain text"
	got := Strip(input)
	if got != input {
		t.Errorf("Strip() = %q, want %q", got, input)
	}
}
