package snapshot_test

import (
	"path/filepath"
	"testing"

	"github.com/user/logpilot/internal/snapshot"
)

func TestStore_MultipleWrites(t *testing.T) {
	p := filepath.Join(t.TempDir(), "state.json")
	s, err := snapshot.New(p)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	for i := int64(1); i <= 5; i++ {
		err := s.Save(snapshot.Snapshot{
			Source:  "file.log",
			Read:    i * 10,
			Matched: i * 8,
			Dropped: i * 2,
		})
		if err != nil {
			t.Fatalf("Save iteration %d: %v", i, err)
		}
	}

	got := s.Latest()
	if got == nil {
		t.Fatal("expected non-nil snapshot")
	}
	if got.Read != 50 {
		t.Errorf("expected Read=50, got %d", got.Read)
	}
	if got.Matched != 40 {
		t.Errorf("expected Matched=40, got %d", got.Matched)
	}
}
