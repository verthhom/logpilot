package enrich_test

import (
	"encoding/json"
	"testing"

	"github.com/logpilot/internal/enrich"
)

func TestNew_Valid(t *testing.T) {
	e, err := enrich.New([]string{"env=prod", "region=us-east-1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if e == nil {
		t.Fatal("expected non-nil Enricher")
	}
}

func TestNew_NoRules(t *testing.T) {
	_, err := enrich.New(nil)
	if err != enrich.ErrNoRules {
		t.Fatalf("expected ErrNoRules, got %v", err)
	}
}

func TestNew_BadRule(t *testing.T) {
	_, err := enrich.New([]string{"nodivider"})
	if err == nil {
		t.Fatal("expected error for bad rule")
	}
}

func TestNew_BlankKey(t *testing.T) {
	_, err := enrich.New([]string{"=value"})
	if err == nil {
		t.Fatal("expected error for blank key")
	}
}

func TestApply_InjectsFields(t *testing.T) {
	e, _ := enrich.New([]string{"env=prod"})
	out := e.Apply(`{"msg":"hello"}`)
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(out), &obj); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if obj["env"] != "prod" {
		t.Errorf("expected env=prod, got %v", obj["env"])
	}
}

func TestApply_DoesNotOverwriteExisting(t *testing.T) {
	e, _ := enrich.New([]string{"env=prod"})
	out := e.Apply(`{"env":"dev"}`)
	var obj map[string]interface{}
	json.Unmarshal([]byte(out), &obj)
	if obj["env"] != "dev" {
		t.Errorf("existing field should not be overwritten, got %v", obj["env"])
	}
}

func TestApply_NonJSONPassThrough(t *testing.T) {
	e, _ := enrich.New([]string{"env=prod"})
	raw := "not json at all"
	if got := e.Apply(raw); got != raw {
		t.Errorf("expected pass-through, got %q", got)
	}
}

func TestApply_MultipleFields(t *testing.T) {
	e, _ := enrich.New([]string{"env=prod", "region=eu"})
	out := e.Apply(`{"msg":"ok"}`)
	var obj map[string]interface{}
	json.Unmarshal([]byte(out), &obj)
	if obj["region"] != "eu" {
		t.Errorf("expected region=eu, got %v", obj["region"])
	}
}
