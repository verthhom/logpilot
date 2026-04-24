package jsonprefix_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/yourorg/logpilot/internal/jsonprefix"
)

func TestNew_Valid(t *testing.T) {
	p, err := jsonprefix.New("app_")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p == nil {
		t.Fatal("expected non-nil Prefixer")
	}
}

func TestNew_EmptyPrefix(t *testing.T) {
	_, err := jsonprefix.New("")
	if err == nil {
		t.Fatal("expected error for empty prefix")
	}
	if err != jsonprefix.ErrEmptyPrefix {
		t.Fatalf("expected ErrEmptyPrefix, got %v", err)
	}
}

func TestNew_WhitespacePrefix(t *testing.T) {
	_, err := jsonprefix.New("   ")
	if err != jsonprefix.ErrEmptyPrefix {
		t.Fatalf("expected ErrEmptyPrefix, got %v", err)
	}
}

func TestApply_PrefixesKeys(t *testing.T) {
	p, _ := jsonprefix.New("app_")
	line := `{"level":"info","msg":"hello"}`
	result := p.Apply(line)

	var obj map[string]json.RawMessage
	if err := json.Unmarshal([]byte(result), &obj); err != nil {
		t.Fatalf("result is not valid JSON: %v", err)
	}
	for k := range obj {
		if !strings.HasPrefix(k, "app_") {
			t.Errorf("key %q missing prefix", k)
		}
	}
}

func TestApply_DoesNotDoublePrefixKeys(t *testing.T) {
	p, _ := jsonprefix.New("app_")
	line := `{"app_level":"info"}`
	result := p.Apply(line)

	var obj map[string]json.RawMessage
	if err := json.Unmarshal([]byte(result), &obj); err != nil {
		t.Fatalf("result is not valid JSON: %v", err)
	}
	if _, ok := obj["app_level"]; !ok {
		t.Error("expected key app_level to be present unchanged")
	}
	if _, ok := obj["app_app_level"]; ok {
		t.Error("key was double-prefixed")
	}
}

func TestApply_NonJSONPassThrough(t *testing.T) {
	p, _ := jsonprefix.New("app_")
	line := "plain text log line"
	if got := p.Apply(line); got != line {
		t.Errorf("expected pass-through, got %q", got)
	}
}

func TestApply_EmptyObject(t *testing.T) {
	p, _ := jsonprefix.New("svc_")
	result := p.Apply(`{}`)
	if result != `{}` {
		t.Errorf("unexpected result for empty object: %q", result)
	}
}
