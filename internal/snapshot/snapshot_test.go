package snapshot

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func tempPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "snap.json")
}

func TestNew_EmptyPath(t *testing.T) {
	_, err := New("")
	if err != ErrEmptyPath {
		t.Fatalf("expected ErrEmptyPath, got %v", err)
	}
}

func TestNew_ValidPath(t *testing.T) {
	s, err := New(tempPath(t))
	if err != nil || s == nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestLatest_InitiallyNil(t *testing.T) {
	s, _ := New(tempPath(t))
	if s.Latest() != nil {
		t.Fatal("expected nil before any save")
	}
}

func TestSave_UpdatesLatest(t *testing.T) {
	s, _ := New(tempPath(t))
	snap := Snapshot{Source: "app.log", Read: 10, Matched: 8, Dropped: 2}
	if err := s.Save(snap); err != nil {
		t.Fatalf("Save error: %v", err)
	}
	got := s.Latest()
	if got == nil {
		t.Fatal("expected snapshot, got nil")
	}
	if got.Read != 10 || got.Matched != 8 || got.Dropped != 2 {
		t.Errorf("unexpected values: %+v", got)
	}
	if got.Timestamp.IsZero() {
		t.Error("timestamp should be set")
	}
}

func TestSave_WritesFile(t *testing.T) {
	p := tempPath(t)
	s, _ := New(p)
	snap := Snapshot{Source: "stdin", Read: 5, Matched: 3, Dropped: 1}
	if err := s.Save(snap); err != nil {
		t.Fatalf("Save error: %v", err)
	}
	data, err := os.ReadFile(p)
	if err != nil {
		t.Fatalf("ReadFile error: %v", err)
	}
	var decoded Snapshot
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}
	if decoded.Source != "stdin" {
		t.Errorf("expected source stdin, got %s", decoded.Source)
	}
}

func TestLatest_IsCopy(t *testing.T) {
	s, _ := New(tempPath(t))
	s.Save(Snapshot{Read: 1})
	a := s.Latest()
	a.Read = 999
	b := s.Latest()
	if b.Read == 999 {
		t.Error("Latest should return a copy, not a reference")
	}
}
