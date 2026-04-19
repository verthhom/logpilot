package jsonpath

import (
	"testing"
)

func TestGet_TopLevelField(t *testing.T) {
	e := New()
	v, err := e.Get(`{"level":"info"}`, "level")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v != "info" {
		t.Fatalf("expected info, got %v", v)
	}
}

func TestGet_NestedField(t *testing.T) {
	e := New()
	v, err := e.Get(`{"meta":{"host":"srv1"}}`, "meta.host")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v != "srv1" {
		t.Fatalf("expected srv1, got %v", v)
	}
}

func TestGet_MissingField(t *testing.T) {
	e := New()
	_, err := e.Get(`{"a":1}`, "b")
	if err != ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestGet_InvalidJSON(t *testing.T) {
	e := New()
	_, err := e.Get(`not-json`, "a")
	if err == nil {
		t.Fatal("expected error for invalid json")
	}
}

func TestGetString_NumericField(t *testing.T) {
	e := New()
	s, err := e.GetString(`{"code":404}`, "code")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s != "404" {
		t.Fatalf("expected 404, got %s", s)
	}
}

func TestGet_DeepNesting(t *testing.T) {
	e := New()
	v, err := e.Get(`{"a":{"b":{"c":"deep"}}}`, "a.b.c")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v != "deep" {
		t.Fatalf("expected deep, got %v", v)
	}
}

func TestGet_IntermediateNotObject(t *testing.T) {
	e := New()
	_, err := e.Get(`{"a":"scalar"}`, "a.b")
	if err != ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}
