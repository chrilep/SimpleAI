//go:build darwin
// +build darwin

package modWindowMemory

import (
	"context"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// macOS-specific window position management
//
// macOS coordinate system characteristics:
// - Origin (0,0) is at bottom-left of primary screen (unlike Windows/Linux top-left)
// - Wails abstracts this, but multi-monitor setups may have quirks
// - Window decorations are handled by the OS consistently
//
// This implementation:
// 1. Uses Wails runtime methods directly (they work reliably on macOS)
// 2. No offset compensation needed (macOS is consistent)
// 3. No special fallback mechanisms required

// RestorePosition restores window position (macOS implementation)
func (wpm *WindowPositionManager) RestorePosition(ctx context.Context, windowID string) {
	pos, exists := wpm.positions[windowID]
	if !exists || pos == nil || pos.Width == 0 || pos.Height == 0 {
		println("[WindowPos] No saved position for", windowID)
		return
	}

	println("[WindowPos] Restoring position for", windowID, "- X:", pos.X, "Y:", pos.Y, "W:", pos.Width, "H:", pos.Height)

	// macOS: Simple and reliable
	runtime.WindowSetPosition(ctx, pos.X, pos.Y)
	runtime.WindowSetSize(ctx, pos.Width, pos.Height)
}

// SavePosition saves current window position (macOS implementation)
func (wpm *WindowPositionManager) SavePosition(ctx context.Context, windowID string, storagePath string) {
	defer func() {
		if r := recover(); r != nil {
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
