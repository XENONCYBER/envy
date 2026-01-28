//go:build windows
// +build windows

package storage

import (
	"fmt"
	"os"
	"sync"
)

type FileLock struct {
	file *os.File
	path string
	mu   sync.Mutex
}

var globalLockMu sync.Mutex

func AcquireLock(lockPath string) (*FileLock, error) {
	globalLockMu.Lock()

	file, err := os.OpenFile(lockPath, os.O_CREATE|os.O_RDWR, 0o600)
	if err != nil {
		globalLockMu.Unlock()
		return nil, fmt.Errorf("failed to open lock file: %w", err)
	}

	return &FileLock{file: file, path: lockPath}, nil
}

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

func (l *FileLock) Release() error {
	if l == nil || l.file == nil {
		return nil
	}

	err := l.file.Close()
	globalLockMu.Unlock()
	return err
}
