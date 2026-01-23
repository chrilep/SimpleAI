//go:build windows
// +build windows

package main

import (
	"syscall"
	"unsafe"
)

var (
	user32                   = syscall.NewLazyDLL("user32.dll")
	procEnumWindows          = user32.NewProc("EnumWindows")
	procGetWindowTextW       = user32.NewProc("GetWindowTextW")
	procSetForegroundWindow  = user32.NewProc("SetForegroundWindow")
	procShowWindow           = user32.NewProc("ShowWindow")
	procIsIconic             = user32.NewProc("IsIconic")
	procGetWindowTextLengthW = user32.NewProc("GetWindowTextLengthW")
)

// findWindowByTitle searches for a window with the exact title
func findWindowByTitle(title string) (uintptr, error) {
	var foundHwnd uintptr
	targetTitle := title

	// Callback function for EnumWindows
	callback := syscall.NewCallback(func(hwnd uintptr, lparam uintptr) uintptr {
		// Get window title length
		length, _, _ := procGetWindowTextLengthW.Call(hwnd)
		if length == 0 {
			return 1 // Continue enumeration
		}

		// Get window title
		buf := make([]uint16, length+1)
		procGetWindowTextW.Call(hwnd, uintptr(unsafe.Pointer(&buf[0])), uintptr(length+1))
		windowTitle := syscall.UTF16ToString(buf)

		// Check if title matches
		if windowTitle == targetTitle {
			foundHwnd = hwnd
			return 0 // Stop enumeration
		}

		return 1 // Continue enumeration
	})

	// Enumerate all top-level windows
	procEnumWindows.Call(callback, 0)

	return foundHwnd, nil
}

// setForegroundWindow brings the specified window to the foreground
func setForegroundWindow(hwnd uintptr) bool {
	ret, _, _ := procSetForegroundWindow.Call(hwnd)
	return ret != 0
}

// showWindow shows the window in the specified state
func showWindow(hwnd uintptr, nCmdShow int) bool {
	ret, _, _ := procShowWindow.Call(hwnd, uintptr(nCmdShow))
	return ret != 0
}

// isIconic checks if the window is minimized
func isIconic(hwnd uintptr) bool {
	ret, _, _ := procIsIconic.Call(hwnd)
	return ret != 0
}

// findAndActivateWindow searches for a window with the given title and activates it
// Returns true if window was found and activated, false otherwise
func findAndActivateWindow(windowTitle string) (bool, error) {
	// Use native Go syscall for direct Windows API access - much faster than PowerShell
	hwnd, err := findWindowByTitle(windowTitle)
	if err != nil || hwnd == 0 {
		return false, nil
	}

	// Check if window is minimized (iconic)
	if isIconic(hwnd) {
		// SW_RESTORE = 9
		showWindow(hwnd, 9)
	}

	// Bring window to foreground
	setForegroundWindow(hwnd)
	return true, nil
}
