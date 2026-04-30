package jsonlookup

import (
	"encoding/json"
	"testing"
)

var sampleTable = map[string]string{
	"us-east-1": "US East",
	"eu-west-1": "EU West",
}

func TestNew_Valid(t *testing.T) {
	_, err := New("region", sampleTable, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNew_BlankField(t *testing.T) {
	_, err := New("  ", sampleTable, "")
	if err != ErrBlankField {
		t.Fatalf("expected ErrBlankField, got %v", err)
	}
}

func TestNew_EmptyTable(t *testing.T) {
	_, err := New("region", map[string]string{}, "")
	if err != ErrEmptyTable {
		t.Fatalf("expected ErrEmptyTable, got %v", err)
	}
}

func TestApply_ReplacesField(t *testing.T) {
	l, _ := New("region", sampleTable, "")
	out := l.Apply(`{"region":"us-east-1","host":"web1"}`)

	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(out), &obj); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if obj["region"] != "US East" {
		t.Errorf("expected 'US East', got %v", obj["region"])
	}
}

func TestApply_WritesToDestField(t *testing.T) {
	l, _ := New("region", sampleTable, "region_label")
	out := l.Apply(`{"region":"eu-west-1"}`)

	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(out), &obj); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if obj["region_label"] != "EU West" {
		t.Errorf("expected 'EU West', got %v", obj["region_label"])
	}
	if obj["region"] != "eu-west-1" {
		t.Errorf("source field should be preserved, got %v", obj["region"])
	}
}

func TestApply_MissingKey_PassThrough(t *testing.T) {
	l, _ := New("region", sampleTable, "")
	input := `{"region":"ap-southeast-1"}`
	out := l.Apply(input)
	if out != input {
		t.Errorf("expected passthrough, got %s", out)
	}
}

func TestApply_NonJSONPassThrough(t *testing.T) {
	l, _ := New("region", sampleTable, "")
	input := "not json at all"
	if out := l.Apply(input); out != input {
		t.Errorf("expected passthrough, got %s", out)
	}
}

func TestApply_NonStringFieldValue_PassThrough(t *testing.T) {
	l, _ := New("code", map[string]string{"42": "ok"}, "")
	input := `{"code":42}`
	if out := l.Apply(input); out != input {
		t.Errorf("expected passthrough for numeric value, got %s", out)
	}
}
