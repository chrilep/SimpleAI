//go:build windows
// +build windows

package modWindowMemory

import (
	"os"
	"syscall"
	"time"
	"unsafe"
)

// Windows API constants for file locking
const (
	LOCKFILE_EXCLUSIVE_LOCK   = 0x00000002
	LOCKFILE_FAIL_IMMEDIATELY = 0x00000001
)

// Windows API function
var (
	kernel32       = syscall.NewLazyDLL("kernel32.dll")
	procLockFileEx = kernel32.NewProc("LockFileEx")
)

// lockFileEx wraps the Windows LockFileEx API
func lockFileEx(handle syscall.Handle, flags, reserved, lockLow, lockHigh uint32, overlapped *syscall.Overlapped) error {
	r1, _, err := procLockFileEx.Call(
		uintptr(handle),
		uintptr(flags),
		uintptr(reserved),
		uintptr(lockLow),
		uintptr(lockHigh),
		uintptr(unsafe.Pointer(overlapped)),
	)
	if r1 == 0 {
		return err
	}
	return nil
}

// openWithLock opens a file with exclusive locking on Windows
// Uses Windows LockFileEx API for proper file locking
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
	for attempts := 0; attempts < 10; attempts++ {
		// Windows file locking via LockFileEx
		overlapped := &syscall.Overlapped{}
		lockFlags := uint32(LOCKFILE_FAIL_IMMEDIATELY)

		if exclusive {
			lockFlags |= LOCKFILE_EXCLUSIVE_LOCK
		}

		handle := syscall.Handle(file.Fd())
		err = lockFileEx(handle, lockFlags, 0, 1, 0, overlapped)

		if err == nil {
			return file, nil
		}

		// If lock failed, wait and retry
		if attempts < 9 {
			time.Sleep(20 * time.Millisecond)
		}
	}

	file.Close()
	return nil, err
}
