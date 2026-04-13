package storage

import (
	"fmt"
	"sync"
)

var (
	mu    sync.RWMutex
	locks = make(map[string]*sync.Mutex)
)

// fileLock returns the per-file mutex for the given path, creating it if needed.
func fileLock(path string) *sync.Mutex {
	mu.Lock()
	defer mu.Unlock()
	if _, ok := locks[path]; !ok {
		locks[path] = &sync.Mutex{}
	}
	return locks[path]
}

// withFileLock executes fn while holding the exclusive lock for path.
func withFileLock(path string, fn func() error) error {
	l := fileLock(path)
	l.Lock()
	defer l.Unlock()
	if err := fn(); err != nil {
		return fmt.Errorf("locked op on %s: %w", path, err)
	}
	return nil
}
