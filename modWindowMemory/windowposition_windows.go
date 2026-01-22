//go:build windows
// +build windows

package modWindowMemory

import (
	"context"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// Windows-specific window position management
//
// Windows has specific behavior with window decorations:
// - WindowGetPosition returns coordinates INCLUDING titlebar and borders
// - WindowSetPosition expects coordinates EXCLUDING decorations
// This creates an offset that must be detected and compensated for.
//
// This implementation:
// 1. Detects the offset on first restore by comparing set vs. get coordinates
// 2. Stores the offset in the manager for subsequent operations
// 3. Compensates automatically on all position sets

// RestorePosition restores window position (Windows implementation)
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

// SavePosition saves current window position (Windows implementation)
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
