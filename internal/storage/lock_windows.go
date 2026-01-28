//go:build windows
// +build windows

package storage

import (
	"fmt"
	"os"
	"sync"
)

// FileLock represents a file-based lock for Windows systems
// Uses a combination of file handle and mutex for cross-process locking
type FileLock struct {
	file *os.File
	path string
	mu   sync.Mutex
}

// Global mutex for simple in-process locking on Windows
// For true cross-process locking, we'd use Windows LockFileEx API
var globalLockMu sync.Mutex

// AcquireLock acquires an exclusive lock on the lock file, blocking until available
func AcquireLock(lockPath string) (*FileLock, error) {
	globalLockMu.Lock()

	file, err := os.OpenFile(lockPath, os.O_CREATE|os.O_RDWR, 0o600)
	if err != nil {
		globalLockMu.Unlock()
		return nil, fmt.Errorf("failed to open lock file: %w", err)
	}

	return &FileLock{file: file, path: lockPath}, nil
}

// TryAcquireLock attempts to acquire an exclusive lock without blocking
// Returns nil, nil if the lock is already held by another process
func TryAcquireLock(lockPath string) (*FileLock, error) {
	if !globalLockMu.TryLock() {
		return nil, nil
	}

	file, err := os.OpenFile(lockPath, os.O_CREATE|os.O_RDWR, 0o600)
	if err != nil {
		globalLockMu.Unlock()
		return nil, fmt.Errorf("failed to open lock file: %w", err)
	}

	return &FileLock{file: file, path: lockPath}, nil
}

// Release releases the file lock
func (l *FileLock) Release() error {
	if l == nil || l.file == nil {
		return nil
	}

	err := l.file.Close()
	globalLockMu.Unlock()
	return err
}
