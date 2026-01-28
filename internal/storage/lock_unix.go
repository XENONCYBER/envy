//go:build unix
// +build unix

package storage

import (
	"fmt"
	"os"
	"syscall"
)

type FileLock struct {
	file *os.File
	path string
}

func AcquireLock(lockPath string) (*FileLock, error) {
	file, err := os.OpenFile(lockPath, os.O_CREATE|os.O_RDWR, 0o600)
	if err != nil {
		return nil, fmt.Errorf("failed to open lock file: %w", err)
	}

	if err := syscall.Flock(int(file.Fd()), syscall.LOCK_EX); err != nil {
		file.Close()
		return nil, fmt.Errorf("failed to acquire lock: %w", err)
	}

	return &FileLock{file: file, path: lockPath}, nil
}

func TryAcquireLock(lockPath string) (*FileLock, error) {
	file, err := os.OpenFile(lockPath, os.O_CREATE|os.O_RDWR, 0o600)
	if err != nil {
		return nil, fmt.Errorf("failed to open lock file: %w", err)
	}

	if err := syscall.Flock(int(file.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
		file.Close()
		if err == syscall.EWOULDBLOCK {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to acquire lock: %w", err)
	}

	return &FileLock{file: file, path: lockPath}, nil
}

func (l *FileLock) Release() error {
	if l == nil || l.file == nil {
		return nil
	}

	if err := syscall.Flock(int(l.file.Fd()), syscall.LOCK_UN); err != nil {
		l.file.Close()
		return fmt.Errorf("failed to release lock: %w", err)
	}

	return l.file.Close()
}
