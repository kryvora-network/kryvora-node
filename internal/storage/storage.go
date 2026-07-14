package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

type Store struct {
	dir string
	mu  sync.Mutex
}

func Open(dir string) (*Store, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create data dir: %w", err)
	}
	return &Store{dir: dir}, nil
}

func (s *Store) Sync() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return nil
}

func (s *Store) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return nil
}

func (s *Store) Path(name string) string {
	return filepath.Join(s.dir, name)
}
