//go:build linux
// +build linux

package modWindowMemory

import (
	"os"
	"syscall"
	"time"
)

// openWithLock opens a file with exclusive locking on Linux
// Uses flock system call for advisory file locking
func openWithLock(path string, flags int, exclusive bool) (*os.File, error) {
	// Ensure file exists for read operations
	if flags&os.O_RDONLY != 0 {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return nil, err
		}
	}

	file, err := os.OpenFile(path, flags, 0644)
	if err != nil {
		return nil, err
	}

	// Try to lock the file with timeout
	lockType := syscall.LOCK_SH // Shared lock for reading
	if exclusive {
		lockType = syscall.LOCK_EX // Exclusive lock for writing
	}

	for attempts := 0; attempts < 10; attempts++ {
		err = syscall.Flock(int(file.Fd()), lockType|syscall.LOCK_NB)
		if err == nil {
			return file, nil
		}

		// If lock failed, wait and retry
		if attempts < 9 {
			time.Sleep(time.Millisecond * 20)
		}
	}

	file.Close()
	return nil, err
}
