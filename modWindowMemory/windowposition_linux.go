//go:build linux
// +build linux

package modWindowMemory

import (
	"context"
	"os/exec"
	"strconv"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// Linux-specific window position management
//
// Linux/GTK has known issues with Wails window geometry APIs:
// - WindowGetPosition often returns (0, 0) even when window is moved
// - WindowGetSize may return default values (1024x768) instead of actual size
// - GTK event loop timing can cause stale values to be read
//
// This implementation provides fallback mechanisms:
// 1. Try Wails runtime methods first (may work on some GTK versions)
// 2. Fall back to xdotool if available (reads from X11 directly)
// 3. Gracefully handle missing xdotool by skipping save
//
// Requirements for full functionality:
// - xdotool must be installed via package manager:
//   Debian/Ubuntu: sudo apt-get install xdotool
//   openSUSE/SUSE: sudo zypper install xdotool
//   Fedora/RHEL: sudo dnf install xdotool
//   Arch Linux: sudo pacman -S xdotool
// - X11 display server (Wayland support may vary)

// getLinuxWindowGeometry attempts to get window geometry using xdotool
// This bypasses GTK/Wails issues by querying X11 directly
func getLinuxWindowGeometry() (x, y, width, height int, ok bool) {
	// Search for window by title pattern - get only the first/active window
	cmd := exec.Command("xdotool", "search", "--name", "^SimpleAI", "getwindowgeometry", "--shell")
	output, err := cmd.Output()
	if err != nil {
		// xdotool not installed or window not found
		println("[WindowPos] xdotool command failed:", err.Error())
		return 0, 0, 0, 0, false
	}

	outputStr := string(output)
	println("[WindowPos] xdotool raw output:", outputStr)

	// Parse xdotool output format:
	// WINDOW=123456
	// X=100
	// Y=200
	// WIDTH=800
	// HEIGHT=600
	lines := strings.Split(outputStr, "\n")

	// Variables to track if we found all required values
	foundX, foundY, foundWidth, foundHeight := false, false, false, false

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.Split(line, "=")
		if len(parts) != 2 {
			println("[WindowPos] Skipping malformed line:", line)
			continue
		}

		key := strings.TrimSpace(parts[0])
		valueStr := strings.TrimSpace(parts[1])
		val, err := strconv.Atoi(valueStr)
		if err != nil {
			println("[WindowPos] Failed to parse value for", key, ":", valueStr)
			continue
		}

		switch key {
		case "X":
			x = val
			foundX = true
		case "Y":
			y = val
			foundY = true
		case "WIDTH":
			width = val
			foundWidth = true
		case "HEIGHT":
			height = val
			foundHeight = true
		}
	}

	println("[WindowPos] Parsed values - X:", x, "Y:", y, "W:", width, "H:", height)
	println("[WindowPos] Found flags - X:", foundX, "Y:", foundY, "W:", foundWidth, "H:", foundHeight)

	// Validate that we got reasonable values
	// Reject suspiciously small windows (10x10 is typically a destroyed/closing window)
	if foundWidth && foundHeight && width > 50 && height > 50 {
		// Accept even if X/Y are 0 - could be valid screen position
		return x, y, width, height, true
	}

	println("[WindowPos] Invalid or incomplete geometry data (too small or missing)")
	return 0, 0, 0, 0, false
}

// RestorePosition restores window position (Linux implementation)
func (wpm *WindowPositionManager) RestorePosition(ctx context.Context, windowID string) {
	pos, exists := wpm.positions[windowID]
	if !exists || pos == nil || pos.Width == 0 || pos.Height == 0 {
		println("[WindowPos] No saved position for", windowID)
		return
	}

	println("[WindowPos] Restoring position for", windowID, "- X:", pos.X, "Y:", pos.Y, "W:", pos.Width, "H:", pos.Height)

	// For GTK, set size before position (order matters)
	runtime.WindowSetSize(ctx, pos.Width, pos.Height)
	runtime.WindowSetPosition(ctx, pos.X, pos.Y)

	// Note: We don't verify the position was set correctly on Linux
	// because WindowGetPosition is unreliable. The position will be
	// verified on next save attempt.
}

// SavePosition saves current window position (Linux implementation)
func (wpm *WindowPositionManager) SavePosition(ctx context.Context, windowID string, storagePath string) {
	defer func() {
		if r := recover(); r != nil {
			println("[WindowPos] Recovered from panic during save:", r)
		}
	}()

	// First, try Wails runtime methods
	x, y := runtime.WindowGetPosition(ctx)
	width, height := runtime.WindowGetSize(ctx)

	// Check if we got default/invalid values (common GTK issue)
	if x == 0 && y == 0 {
		println("[WindowPos] Wails returned (0,0), trying xdotool fallback")
		// Try xdotool as fallback
		xX, xY, xWidth, xHeight, ok := getLinuxWindowGeometry()
		if ok {
			println("[WindowPos] xdotool success - X:", xX, "Y:", xY, "W:", xWidth, "H:", xHeight)
			x, y, width, height = xX, xY, xWidth, xHeight
		} else {
			println("[WindowPos] xdotool failed - install xdotool for Linux position tracking")
			println("[WindowPos] Debian/Ubuntu: sudo apt-get install xdotool")
			println("[WindowPos] openSUSE/SUSE: sudo zypper install xdotool")
			println("[WindowPos] Fedora/RHEL: sudo dnf install xdotool")
			println("[WindowPos] Arch Linux: sudo pacman -S xdotool")
			return
		}
	}

	// Don't save invalid dimensions
	if width == 0 || height == 0 {
		println("[WindowPos] Skipping save - invalid dimensions")
		return
	}

	println("[WindowPos] Saving position for", windowID, "- X:", x, "Y:", y, "W:", width, "H:", height)

	// Reload from disk to preserve positions of other running instances
	wpm.Load(storagePath)

	wpm.positions[windowID] = &WindowPosition{
		X:      x,
		Y:      y,
		Width:  width,
		Height: height,
	}

	wpm.Save(storagePath)
}
