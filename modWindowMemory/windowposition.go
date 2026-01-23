package modWindowMemory

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// WindowPositionManager handles persistent window positioning across platforms.
//
// This module provides a universal solution for saving and restoring window positions
// in Wails applications. It abstracts platform-specific differences in window geometry
// handling (Windows decorations, Linux/GTK quirks, macOS coordinate systems).
//
// Architecture:
//   - windowposition.go: Platform-independent storage/retrieval logic
//   - windowposition_windows.go: Windows-specific geometry handling
//   - windowposition_linux.go: Linux/GTK-specific geometry handling
//   - windowposition_darwin.go: macOS-specific geometry handling
//
// Usage in any Wails project:
//  1. Create manager: wpm := NewWindowPositionManager()
//  2. Load saved state: wpm.Load(storagePath)
//  3. On app startup: wpm.RestorePosition(ctx, windowID)
//  4. On app shutdown: wpm.SavePosition(ctx, windowID, storagePath)
//
// Window IDs are typically the window title, allowing multiple windows to be tracked.
type WindowPositionManager struct {
	positions map[string]*WindowPosition
	xOffset   int          // Platform-specific offset X (e.g., Windows border)
	yOffset   int          // Platform-specific offset Y (e.g., Windows titlebar)
	mu        sync.RWMutex // Protects positions map from concurrent access
}

// WindowPosition stores position and size for a single window
type WindowPosition struct {
	X      int `json:"x"`
	Y      int `json:"y"`
	Width  int `json:"width"`
	Height int `json:"height"`
}

// NewWindowPositionManager creates a new window position manager
func NewWindowPositionManager() *WindowPositionManager {
	return &WindowPositionManager{
		positions: make(map[string]*WindowPosition),
	}
}

// Load reads saved window positions from disk
// storagePath: full path to JSON file (e.g., "path/to/windows.json")
func (wpm *WindowPositionManager) Load(storagePath string) error {
	wpm.mu.Lock()
	defer wpm.mu.Unlock()

	// Ensure directory exists
	dir := filepath.Dir(storagePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Use file locking to prevent race conditions with other instances
	file, err := openWithLock(storagePath, os.O_RDONLY, false)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // File doesn't exist yet, not an error
		}
		return err
	}
	defer file.Close()

	data, err := os.ReadFile(storagePath)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &wpm.positions)
}

// Save writes current window positions to disk
// storagePath: full path to JSON file (e.g., "path/to/windows.json")
func (wpm *WindowPositionManager) Save(storagePath string) error {
	wpm.mu.RLock()
	data, err := json.Marshal(wpm.positions)
	wpm.mu.RUnlock()

	if err != nil {
		return err
	}

	// Use file locking with retry to handle concurrent writes from multiple instances
	for attempts := 0; attempts < 5; attempts++ {
		file, err := openWithLock(storagePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, true)
		if err != nil {
			if attempts < 4 {
				time.Sleep(time.Millisecond * 50) // Wait before retry
				continue
			}
			return err
		}

		_, err = file.Write(data)
		file.Close()
		return err
	}

	return nil
}

// RestorePosition restores window position for a given window ID.
// windowID: unique identifier (typically window title)
//
// Platform-specific implementation in windowposition_*.go files.
// Each platform handles coordinate systems and decorations differently.

// SavePosition saves current window position for a given window ID.
// windowID: unique identifier (typically window title)
// storagePath: full path to JSON file (e.g., "path/to/windows.json")
//
// Platform-specific implementation in windowposition_*.go files.
// Each platform may require different methods to reliably read window geometry.

// GetPosition returns the saved position for a page ID, or nil if not found
func (wpm *WindowPositionManager) GetPosition(pageID string) *WindowPosition {
	wpm.mu.RLock()
	defer wpm.mu.RUnlock()
	return wpm.positions[pageID]
}

// SetPosition sets the position for a page ID (doesn't save to disk)
func (wpm *WindowPositionManager) SetPosition(pageID string, x, y, width, height int) {
	wpm.mu.Lock()
	defer wpm.mu.Unlock()
	wpm.positions[pageID] = &WindowPosition{
		X:      x,
		Y:      y,
		Width:  width,
		Height: height,
	}
}

// validateAndCorrectPosition ensures window position is within visible screen bounds.
// Returns corrected position coordinates that keep the window fully visible.
//
// Parameters:
//   - x, y: requested window position
//   - width, height: window dimensions
//   - screenWidth, screenHeight: primary screen dimensions
//
// Returns corrected x, y coordinates.
func validateAndCorrectPosition(x, y, width, height, screenWidth, screenHeight int) (int, int) {
	const minVisibleOffset = 20 // Minimum pixels that must remain visible

	correctedX := x
	correctedY := y

	// Ensure window is not too far left
	if correctedX < -width+minVisibleOffset {
		correctedX = -width + minVisibleOffset
	}

	// Ensure window is not too far right
	if correctedX > screenWidth-minVisibleOffset {
		correctedX = screenWidth - minVisibleOffset
	}

	// Ensure window is not too far up (negative Y = above screen)
	if correctedY < 0 {
		correctedY = 0
	}

	// Ensure window is not too far down
	if correctedY > screenHeight-minVisibleOffset {
		correctedY = screenHeight - minVisibleOffset
	}

	if correctedX != x || correctedY != y {
		println("[WindowPos] Position corrected from (", x, ",", y, ") to (", correctedX, ",", correctedY, ") to stay within screen bounds")
	}

	return correctedX, correctedY
}
