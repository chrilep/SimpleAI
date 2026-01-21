package main

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// WindowPositionManager handles persistent window positioning with automatic
// compensation for platform-specific titlebar/border offsets (e.g., Windows).
//
// Usage:
//  1. Create manager: wpm := NewWindowPositionManager()
//  2. Load state: wpm.Load(storagePath)
//  3. On startup: wpm.RestorePosition(ctx)
//  4. On shutdown: wpm.SavePosition(ctx, storagePath)
//
// The manager automatically detects and compensates for coordinate system
// differences between WindowSetPosition and WindowGetPosition on Windows.
// Window titles are used as unique identifiers for storing positions.
type WindowPositionManager struct {
	positions map[string]*WindowPosition
	xOffset   int // Platform offset X (e.g., Windows border)
	yOffset   int // Platform offset Y (e.g., Windows titlebar)
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

// RestorePosition restores window position for a given window ID
// windowID: unique identifier (e.g., window title)
// Automatically detects and compensates for platform-specific offsets
func (wpm *WindowPositionManager) RestorePosition(ctx context.Context, windowID string) {
	pos, exists := wpm.positions[windowID]
	if !exists || pos == nil || pos.Width == 0 || pos.Height == 0 {
		println("[WindowPos] No saved position for", windowID)
		return
	}

	println("[WindowPos] Restoring position for", windowID, "- X:", pos.X, "Y:", pos.Y, "W:", pos.Width, "H:", pos.Height)

	// Apply offset compensation (discovered on previous run)
	targetX := pos.X - wpm.xOffset
	targetY := pos.Y - wpm.yOffset
	runtime.WindowSetPosition(ctx, targetX, targetY)
	runtime.WindowSetSize(ctx, pos.Width, pos.Height)

	// Measure actual offset and re-apply if needed (first run or offset changed)
	actualX, actualY := runtime.WindowGetPosition(ctx)
	offsetX := actualX - targetX
	offsetY := actualY - targetY

	if offsetX != 0 || offsetY != 0 {
		println("[WindowPos] Detected offset - X:", offsetX, "Y:", offsetY, "- Compensating immediately")
		wpm.xOffset = offsetX
		wpm.yOffset = offsetY
		// Re-apply with compensation
		runtime.WindowSetPosition(ctx, pos.X-wpm.xOffset, pos.Y-wpm.yOffset)
	}
}

// SavePosition saves current window position for a given window ID
// windowID: unique identifier (e.g., window title)
// storagePath: full path to JSON file (e.g., "path/to/windows.json")
// Returns early if window is being destroyed (width/height = 0)
func (wpm *WindowPositionManager) SavePosition(ctx context.Context, windowID string, storagePath string) {
	defer func() {
		if r := recover(); r != nil {
			// Ignore panics during shutdown - window may be destroyed
			println("[WindowPos] Recovered from panic during save:", r)
		}
	}()

	x, y := runtime.WindowGetPosition(ctx)
	width, height := runtime.WindowGetSize(ctx)

	// Don't save invalid dimensions (happens during shutdown)
	if width == 0 || height == 0 {
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
