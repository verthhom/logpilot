package jsonxform_test

import (
	"encoding/json"
	"testing"

	"github.com/yourorg/logpilot/internal/jsonxform"
)

// TestXformer_ChainedApplication verifies that two Transformers can be
// applied sequentially, each building on the previous result.
func TestXformer_ChainedApplication(t *testing.T) {
	xf1, err := jsonxform.New(`level_bracket=level:[{{.Value}}]`)
	if err != nil {
		t.Fatalf("New xf1: %v", err)
	}
	xf2, err := jsonxform.New(`summary=msg:event={{.Value}}`)
	if err != nil {
		t.Fatalf("New xf2: %v", err)
	}

	input := `{"level":"warn","msg":"disk full"}`
	out := xf2.Apply(xf1.Apply(input))

	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(out), &obj); err != nil {
		t.Fatalf("output not valid JSON: %v", err)
	}

	if obj["level_bracket"] != "[warn]" {
		t.Errorf("level_bracket = %q, want %q", obj["level_bracket"], "[warn]")
	}
	if obj["summary"] != "event=disk full" {
		t.Errorf("summary = %q, want %q", obj["summary"], "event=disk full")
	}
}
