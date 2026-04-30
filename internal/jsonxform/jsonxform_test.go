package jsonxform

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestNew_Valid(t *testing.T) {
	_, err := New(`upper=msg:{{.Value}}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNew_EmptyRule(t *testing.T) {
	_, err := New("")
	if err == nil {
		t.Fatal("expected error for empty rule")
	}
}

func TestNew_BadRule_NoEquals(t *testing.T) {
	_, err := New("nodest")
	if err == nil {
		t.Fatal("expected error for missing '='")
	}
}

func TestNew_BadRule_NoColon(t *testing.T) {
	_, err := New("dest=srcNoColon")
	if err == nil {
		t.Fatal("expected error for missing ':'")
	}
}

func TestNew_BadTemplate(t *testing.T) {
	_, err := New("dest=src:{{.Value")
	if err == nil {
		t.Fatal("expected error for malformed template")
	}
}

func TestApply_TransformsField(t *testing.T) {
	xf, err := New(`tagged=level:[{{.Value}}]`)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	out := xf.Apply(`{"level":"info","msg":"hi"}`)

	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(out), &obj); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if obj["tagged"] != "[info]" {
		t.Errorf("tagged = %q, want %q", obj["tagged"], "[info]")
	}
}

func TestApply_NonJSONPassThrough(t *testing.T) {
	xf, _ := New(`out=field:{{.Value}}`)
	raw := "plain text line"
	if got := xf.Apply(raw); got != raw {
		t.Errorf("got %q, want %q", got, raw)
	}
}

func TestApply_MissingSourceField_PassThrough(t *testing.T) {
	xf, _ := New(`out=missing:{{.Value}}`)
	raw := `{"level":"debug"}`
	if got := xf.Apply(raw); got != raw {
		t.Errorf("got %q, want %q", got, raw)
	}
}

func TestApply_WritesToDifferentDest(t *testing.T) {
	xf, _ := New(`copy=msg:{{.Value}}`)
	out := xf.Apply(`{"msg":"hello"}`)

	if !strings.Contains(out, `"copy"`) {
		t.Errorf("expected 'copy' field in output: %s", out)
	}
	if !strings.Contains(out, `"msg"`) {
		t.Errorf("expected original 'msg' field preserved: %s", out)
	}
}
