package jsoncast

import (
	"encoding/json"
	"testing"
)

func TestNew_Valid(t *testing.T) {
	_, err := New([]string{"latency:float", "status:int"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNew_NoRules(t *testing.T) {
	_, err := New(nil)
	if err == nil {
		t.Fatal("expected error for empty rules")
	}
}

func TestNew_BadRule(t *testing.T) {
	_, err := New([]string{"nodivider"})
	if err == nil {
		t.Fatal("expected error for malformed rule")
	}
}

func TestNew_UnknownType(t *testing.T) {
	_, err := New([]string{"field:xml"})
	if err == nil {
		t.Fatal("expected error for unknown type")
	}
}

func TestApply_CastsToFloat(t *testing.T) {
	c, _ := New([]string{"latency:float"})
	out := c.Apply(`{"latency":"3.14"}`)
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(out), &obj); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if _, ok := obj["latency"].(float64); !ok {
		t.Fatalf("expected float64, got %T", obj["latency"])
	}
}

func TestApply_CastsToInt(t *testing.T) {
	c, _ := New([]string{"status:int"})
	out := c.Apply(`{"status":"200"}`)
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(out), &obj); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if _, ok := obj["status"].(float64); !ok {
		t.Fatalf("expected numeric, got %T", obj["status"])
	}
}

func TestApply_CastsToBool(t *testing.T) {
	c, _ := New([]string{"retried:bool"})
	out := c.Apply(`{"retried":"true"}`)
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(out), &obj); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if v, ok := obj["retried"].(bool); !ok || !v {
		t.Fatalf("expected bool true, got %v (%T)", obj["retried"], obj["retried"])
	}
}

func TestApply_CastsToString(t *testing.T) {
	c, _ := New([]string{"code:string"})
	out := c.Apply(`{"code":404}`)
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(out), &obj); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if _, ok := obj["code"].(string); !ok {
		t.Fatalf("expected string, got %T", obj["code"])
	}
}

func TestApply_NonJSONPassThrough(t *testing.T) {
	c, _ := New([]string{"x:int"})
	line := "not json at all"
	if got := c.Apply(line); got != line {
		t.Fatalf("expected pass-through, got %q", got)
	}
}

func TestApply_MissingFieldSkipped(t *testing.T) {
	c, _ := New([]string{"missing:int"})
	line := `{"present":"hello"}`
	out := c.Apply(line)
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(out), &obj); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if _, ok := obj["missing"]; ok {
		t.Fatal("expected missing field to remain absent")
	}
}
