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
func getLinuxWindowGeometry(dbg bool) (x, y, width, height int, ok bool) {
	if dbg {
		println("[WindowPos][DEBUG] getLinuxWindowGeometry() called")
	}
	// Search for window by title pattern - get only the first/active window
	cmd := exec.Command("xdotool", "search", "--name", "^SimpleAI", "getwindowgeometry", "--shell")
	output, err := cmd.Output()
	if err != nil {
		// xdotool not installed or window not found
		if dbg {
			println("[WindowPos][DEBUG] xdotool command failed:", err.Error())
		}
		return 0, 0, 0, 0, false
	}

	outputStr := string(output)
	if dbg {
		println("[WindowPos][DEBUG] xdotool raw output:", outputStr)
	}

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
			if dbg {
				println("[WindowPos][DEBUG] Skipping malformed line:", line)
			}
			continue
		}

		key := strings.TrimSpace(parts[0])
		valueStr := strings.TrimSpace(parts[1])
		val, err := strconv.Atoi(valueStr)
		if err != nil {
			if dbg {
				println("[WindowPos][DEBUG] Failed to parse value for", key, ":", valueStr)
			}
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

	if dbg {
		println("[WindowPos][DEBUG] Parsed values - X:", x, "Y:", y, "W:", width, "H:", height)
		println("[WindowPos][DEBUG] Found flags - X:", foundX, "Y:", foundY, "W:", foundWidth, "H:", foundHeight)
	}

	// Validate that we got reasonable values
	// Reject suspiciously small windows (10x10 is typically a destroyed/closing window)
	if foundWidth && foundHeight && width > 50 && height > 50 {
		// Accept even if X/Y are 0 - could be valid screen position
		if dbg {
			println("[WindowPos][DEBUG] Valid geometry found, returning values")
		}
		return x, y, width, height, true
	}

	if dbg {
		println("[WindowPos][DEBUG] Invalid or incomplete geometry data (too small or missing)")
	}
	return 0, 0, 0, 0, false
}

// RestorePosition restores window position (Linux implementation)
func (wpm *WindowPositionManager) RestorePosition(ctx context.Context, windowID string) {
	const dbg = true // Set to false to disable debug logging

	if dbg {
		println("[WindowPos][DEBUG] RestorePosition() called for", windowID)
	}

	wpm.mu.RLock()
	pos, exists := wpm.positions[windowID]
	wpm.mu.RUnlock()

	if !exists || pos == nil || pos.Width == 0 || pos.Height == 0 {
		if dbg {
			println("[WindowPos][DEBUG] No saved position for", windowID)
		}
		return
	}

	if dbg {
		println("[WindowPos][DEBUG] Restoring position for", windowID, "- X:", pos.X, "Y:", pos.Y, "W:", pos.Width, "H:", pos.Height)
	}

	// Get screen dimensions to validate position
	screens, err := runtime.ScreenGetAll(ctx)
	if err == nil && len(screens) > 0 {
		// Use primary screen dimensions
		screenWidth := screens[0].Width
		screenHeight := screens[0].Height

		if dbg {
			println("[WindowPos][DEBUG] Screen dimensions:", screenWidth, "x", screenHeight)
			println("[WindowPos][DEBUG] Before validation - X:", pos.X, "Y:", pos.Y, "W:", pos.Width, "H:", pos.Height)
		}

		// Validate and correct position to stay within screen bounds
		pos.X, pos.Y, pos.Width, pos.Height = validateAndCorrectPosition(pos.X, pos.Y, pos.Width, pos.Height, screenWidth, screenHeight)

		if dbg {
			println("[WindowPos][DEBUG] After validation - X:", pos.X, "Y:", pos.Y, "W:", pos.Width, "H:", pos.Height)
		}
	} else {
		if dbg {
			println("[WindowPos][DEBUG] Warning: Could not get screen dimensions, skipping bounds validation")
		}
	}

	// Try Wails runtime methods first (may work on some GTK versions)
	// For GTK, set size before position (order matters)
	if dbg {
		println("[WindowPos][DEBUG] Calling runtime.WindowSetSize(", pos.Width, ",", pos.Height, ")")
	}
	runtime.WindowSetSize(ctx, pos.Width, pos.Height)

	if dbg {
		println("[WindowPos][DEBUG] Calling runtime.WindowSetPosition(", pos.X, ",", pos.Y, ")")
	}
	runtime.WindowSetPosition(ctx, pos.X, pos.Y)

	// Apply position using xdotool with polling and timeout
	go func() {
		const maxAttempts = 50   // 50 attempts at 100ms = 5 seconds timeout
		const pollInterval = 100 // milliseconds

		if dbg {
			println("[WindowPos][DEBUG] Starting xdotool goroutine")
		}

		// Wait for window to be ready
		windowFound := false
		for attempt := 0; attempt < maxAttempts; attempt++ {
			if dbg && attempt%10 == 0 {
				println("[WindowPos][DEBUG] Polling for window, attempt", attempt)
			}
			// Check if window exists and is ready
			cmd := exec.Command("xdotool", "search", "--name", "^SimpleAI")
			if output, err := cmd.Output(); err == nil && len(output) > 0 {
				windowFound = true
				if dbg {
					println("[WindowPos][DEBUG] Window found after", attempt, "attempts")
				}
				break
			}
			// Poll every 100ms
			exec.Command("sleep", "0.1").Run()
		}

		if !windowFound {
			if dbg {
				println("[WindowPos][DEBUG] ERROR: Window not found after 5 seconds timeout")
			}
			return
		}

		// Apply position now that window is ready
		if dbg {
			println("[WindowPos][DEBUG] Applying position with xdotool windowmove", pos.X, pos.Y)
		}

		cmd := exec.Command("xdotool", "search", "--name", "^SimpleAI", "windowmove",
			strconv.Itoa(pos.X), strconv.Itoa(pos.Y))
		if err := cmd.Run(); err != nil {
			if dbg {
				println("[WindowPos][DEBUG] ERROR: Failed to set position:", err.Error())
			}
			return
		}
		if dbg {
			println("[WindowPos][DEBUG] Applied position using xdotool")
		}

		// Apply size
		if dbg {
			println("[WindowPos][DEBUG] Applying size with xdotool windowsize", pos.Width, pos.Height)
		}

		cmd = exec.Command("xdotool", "search", "--name", "^SimpleAI", "windowsize",
			strconv.Itoa(pos.Width), strconv.Itoa(pos.Height))
		if err := cmd.Run(); err != nil {
			if dbg {
				println("[WindowPos][DEBUG] ERROR: Failed to set size:", err.Error())
			}
			return
		}
		if dbg {
			println("[WindowPos][DEBUG] Applied size using xdotool")
		}

		// Monitor and re-apply position if window manager moves it
		// GTK/WM may reposition the window after initial placement
		const monitorAttempts = 20  // Monitor for 2 seconds (20 x 100ms)
		const monitorInterval = 100 // milliseconds

		if dbg {
			println("[WindowPos][DEBUG] Starting position monitoring for 2 seconds...")
		}

		for i := 0; i < monitorAttempts; i++ {
			exec.Command("sleep", "0.1").Run()

			actualX, actualY, actualWidth, actualHeight, ok := getLinuxWindowGeometry(dbg)
			if !ok {
				if dbg && i == 0 {
					println("[WindowPos][DEBUG] Could not verify position (getLinuxWindowGeometry failed)")
				}
				continue
			}

			if i == 0 && dbg {
				println("[WindowPos][DEBUG] Initial verification - X:", actualX, "Y:", actualY, "W:", actualWidth, "H:", actualHeight)
			}

			// Check if position drifted
			if actualX != pos.X || actualY != pos.Y {
				if dbg {
					println("[WindowPos][DEBUG] ⚠ Position drift detected at", i*monitorInterval, "ms - Expected:", pos.X, pos.Y, "Got:", actualX, actualY, "ΔX:", actualX-pos.X, "ΔY:", actualY-pos.Y)
					println("[WindowPos][DEBUG] Re-applying position...")
				}

				// Re-apply position
				cmd := exec.Command("xdotool", "search", "--name", "^SimpleAI", "windowmove",
					strconv.Itoa(pos.X), strconv.Itoa(pos.Y))
				if err := cmd.Run(); err != nil {
					if dbg {
						println("[WindowPos][DEBUG] Failed to re-apply position:", err.Error())
					}
				} else if dbg {
					println("[WindowPos][DEBUG] Position re-applied")
				}
			} else if i == monitorAttempts-1 && dbg {
				// Last check - position is stable
				println("[WindowPos][DEBUG] ✓ Position stable at X:", actualX, "Y:", actualY, "after", i*monitorInterval, "ms")
			}
		}

		if dbg {
			println("[WindowPos][DEBUG] Position monitoring complete")
		}
	}()
}

// SavePosition saves current window position (Linux implementation)
func (wpm *WindowPositionManager) SavePosition(ctx context.Context, windowID string, storagePath string) {
	const dbg = true // Set to false to disable debug logging

	defer func() {
		if r := recover(); r != nil {
			if dbg {
				println("[WindowPos][DEBUG] Recovered from panic during save:", r)
			}
		}
	}()

	if dbg {
		println("[WindowPos][DEBUG] SavePosition() called for", windowID)
	}

	// First, try Wails runtime methods
	x, y := runtime.WindowGetPosition(ctx)
	width, height := runtime.WindowGetSize(ctx)

	if dbg {
		println("[WindowPos][DEBUG] Wails runtime returned - X:", x, "Y:", y, "W:", width, "H:", height)
	}

	// Check if we got default/invalid values (common GTK issue)
	if x == 0 && y == 0 {
		if dbg {
			println("[WindowPos][DEBUG] Wails returned (0,0), trying xdotool fallback")
		}
		// Try xdotool as fallback
		xX, xY, xWidth, xHeight, ok := getLinuxWindowGeometry(dbg)
		if ok {
			if dbg {
				println("[WindowPos][DEBUG] xdotool success - X:", xX, "Y:", xY, "W:", xWidth, "H:", xHeight)
			}
			x, y, width, height = xX, xY, xWidth, xHeight
		} else {
			if dbg {
				println("[WindowPos][DEBUG] xdotool failed - install xdotool for Linux position tracking")
				println("[WindowPos][DEBUG] Debian/Ubuntu: sudo apt-get install xdotool")
				println("[WindowPos][DEBUG] openSUSE/SUSE: sudo zypper install xdotool")
				println("[WindowPos][DEBUG] Fedora/RHEL: sudo dnf install xdotool")
				println("[WindowPos][DEBUG] Arch Linux: sudo pacman -S xdotool")
			}
			return
		}
	}

	// Don't save invalid dimensions
	if width == 0 || height == 0 {
		if dbg {
			println("[WindowPos][DEBUG] Skipping save - invalid dimensions")
		}
		return
	}

	if dbg {
		println("[WindowPos][DEBUG] Saving position for", windowID, "- X:", x, "Y:", y, "W:", width, "H:", height)
	}

	// Reload from disk to preserve positions of other running instances
	wpm.Load(storagePath)

	wpm.mu.Lock()
	wpm.positions[windowID] = &WindowPosition{
		X:      x,
		Y:      y,
		Width:  width,
		Height: height,
	}
	wpm.mu.Unlock()

	wpm.Save(storagePath)
}
