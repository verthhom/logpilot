package snapshot

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

// Snapshot captures a point-in-time summary of pipeline metrics.
type Snapshot struct {
	Timestamp time.Time `json:"timestamp"`
	Read      int64     `json:"read"`
	Matched   int64     `json:"matched"`
	Dropped   int64     `json:"dropped"`
	Source    string    `json:"source"`
}

// Store holds the latest snapshot and persists it to disk.
type Store struct {
	mu   sync.RWMutex
	path string
	last *Snapshot
}

// New creates a new Store that persists snapshots to path.
func New(path string) (*Store, error) {
	if path == "" {
		return nil, ErrEmptyPath
	}
	return &Store{path: path}, nil
}

// Save stores the snapshot in memory and writes it to disk.
func (s *Store) Save(snap Snapshot) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	snap.Timestamp = time.Now().UTC()
	s.last = &snap
	return s.flush(snap)
}

// Latest returns the most recently saved snapshot, or nil.
func (s *Store) Latest() *Snapshot {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.last == nil {
		return nil
	}
	copy := *s.last
	return &copy
}

func (s *Store) flush(snap Snapshot) error {
	f, err := os.CreateTemp("", "snapshot-*.json")
	if err != nil {
		return err
	}
	tmpName := f.Name()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(snap); err != nil {
		f.Close()
		os.Remove(tmpName)
		return err
	}
	f.Close()
	return os.Rename(tmpName, s.path)
}
