package jsonstrip_test

import (
	"encoding/json"
	"testing"

	"github.com/yourorg/logpilot/internal/jsonstrip"
)

func TestStripper_ChainedApply(t *testing.T) {
	s1, _ := jsonstrip.New([]string{"token"})
	s2, _ := jsonstrip.New([]string{"password"})

	line := `{"user":"bob","token":"xyz","password":"hunter2","level":"info"}`
	after1 := s1.Apply(line)
	after2 := s2.Apply(after1)

	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(after2), &obj); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	for _, removed := range []string{"token", "password"} {
		if _, ok := obj[removed]; ok {
			t.Errorf("field %q should have been removed", removed)
		}
	}
	if obj["user"] != "bob" {
		t.Errorf("user field should be preserved")
	}
}
