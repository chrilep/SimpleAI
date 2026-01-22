package modWindowMemory

import (
	"encoding/json"
	"os"
	"path/filepath"
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
	xOffset   int // Platform-specific offset X (e.g., Windows border)
	yOffset   int // Platform-specific offset Y (e.g., Windows titlebar)
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
	// Ensure directory exists
	dir := filepath.Dir(storagePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := os.ReadFile(storagePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // File doesn't exist yet, not an error
		}
		return err
	}

	return json.Unmarshal(data, &wpm.positions)
}

// Save writes current window positions to disk
// storagePath: full path to JSON file (e.g., "path/to/windows.json")
func (wpm *WindowPositionManager) Save(storagePath string) error {
	data, err := json.Marshal(wpm.positions)
	if err != nil {
		return err
	}

	return os.WriteFile(storagePath, data, 0644)
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
	return wpm.positions[pageID]
}

// SetPosition sets the position for a page ID (doesn't save to disk)
func (wpm *WindowPositionManager) SetPosition(pageID string, x, y, width, height int) {
	wpm.positions[pageID] = &WindowPosition{
		X:      x,
		Y:      y,
		Width:  width,
		Height: height,
	}
}
