package jsonclean_test

import (
	"testing"

	"github.com/yourorg/logpilot/internal/jsonclean"
)

func TestApply_NonJSONPassThrough(t *testing.T) {
	c := jsonclean.New(jsonclean.Options{RemoveNull: true})
	input := "not json at all"
	if got := c.Apply(input); got != input {
		t.Fatalf("expected pass-through, got %q", got)
	}
}

func TestApply_RemovesNullFields(t *testing.T) {
	c := jsonclean.New(jsonclean.Options{RemoveNull: true})
	got := c.Apply(`{"a":"hello","b":null}`)
	if got == "" {
		t.Fatal("unexpected empty output")
	}
	if contains(got, `"b"`) {
		t.Fatalf("null field not removed: %s", got)
	}
	if !contains(got, `"a"`) {
		t.Fatalf("non-null field removed: %s", got)
	}
}

func TestApply_RemovesEmptyStrings(t *testing.T) {
	c := jsonclean.New(jsonclean.Options{RemoveEmptyString: true})
	got := c.Apply(`{"x":"","y":"keep"}`)
	if contains(got, `"x"`) {
		t.Fatalf("empty-string field not removed: %s", got)
	}
	if !contains(got, `"y"`) {
		t.Fatalf("non-empty field removed: %s", got)
	}
}

func TestApply_RemovesEmptyArray(t *testing.T) {
	c := jsonclean.New(jsonclean.Options{RemoveEmptyArray: true})
	got := c.Apply(`{"tags":[],"msg":"hi"}`)
	if contains(got, `"tags"`) {
		t.Fatalf("empty-array field not removed: %s", got)
	}
}

func TestApply_RemovesEmptyObject(t *testing.T) {
	c := jsonclean.New(jsonclean.Options{RemoveEmptyObject: true})
	got := c.Apply(`{"meta":{},"level":"info"}`)
	if contains(got, `"meta"`) {
		t.Fatalf("empty-object field not removed: %s", got)
	}
}

func TestApply_NoOptionsKeepsAll(t *testing.T) {
	c := jsonclean.New(jsonclean.Options{})
	input := `{"a":null,"b":"","c":[],"d":{}}`
	got := c.Apply(input)
	for _, key := range []string{`"a"`, `"b"`, `"c"`, `"d"`} {
		if !contains(got, key) {
			t.Fatalf("field %s unexpectedly removed from %s", key, got)
		}
	}
}

func TestApply_AllOptionsRemoveAll(t *testing.T) {
	c := jsonclean.New(jsonclean.Options{
		RemoveNull:        true,
		RemoveEmptyString: true,
		RemoveEmptyArray:  true,
		RemoveEmptyObject: true,
	})
	got := c.Apply(`{"a":null,"b":"","c":[],"d":{}}`)
	if got != "{}" {
		t.Fatalf("expected empty object, got %s", got)
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && containsStr(s, sub))
}

func containsStr(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
