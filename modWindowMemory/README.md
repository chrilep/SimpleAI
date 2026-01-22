# Window Position Manager

Universal window position persistence module for Wails applications.

## Overview

This module provides platform-independent window position saving and restoring across Windows, Linux, and macOS. It handles platform-specific quirks transparently, making it easy to add window position memory to any Wails project.

## Architecture

```
windowposition.go          → Core logic (storage, JSON, manager struct)
windowposition_windows.go  → Windows-specific geometry handling
windowposition_linux.go    → Linux/GTK-specific geometry handling
windowposition_darwin.go   → macOS-specific geometry handling
```

Build tags (`//go:build`) ensure only the relevant platform file is compiled.

## Platform-Specific Behavior

### Windows

- **Issue**: WindowGetPosition returns coordinates including titlebar/borders, but WindowSetPosition expects coordinates excluding decorations
- **Solution**: Automatic offset detection and compensation on first restore
- **Status**: Fully functional

### Linux/GTK

- **Issue**: WindowGetPosition often returns (0,0) regardless of actual position
- **Solution**: Falls back to `xdotool` to query X11 directly when Wails methods fail
- **Requirements**: Install xdotool via your package manager:
  - Debian/Ubuntu: `sudo apt-get install xdotool`
  - openSUSE/SUSE: `sudo zypper install xdotool`
  - Fedora/RHEL: `sudo dnf install xdotool`
  - Arch Linux: `sudo pacman -S xdotool`
- **Status**: Functional with xdotool; graceful degradation without it

### macOS

- **Issue**: None (macOS APIs are consistent)
- **Solution**: Direct use of Wails runtime methods
- **Status**: Fully functional

## Usage

### Basic Integration

```go
package main

import (
    "context"
    "os"
    "path/filepath"
)

type App struct {
    ctx           context.Context
    windowPosMgr  *WindowPositionManager
    windowPosPath string
}

func NewApp() *App {
    configDir, _ := os.UserConfigDir()
    return &App{
        windowPosMgr:  NewWindowPositionManager(),
        windowPosPath: filepath.Join(configDir, "YourApp", "windows.json"),
    }
}

func (a *App) startup(ctx context.Context) {
    a.ctx = ctx

    // Load saved positions
    a.windowPosMgr.Load(a.windowPosPath)

    // Restore position for this window
    windowTitle := "Your App Window"
    a.windowPosMgr.RestorePosition(ctx, windowTitle)
}

func (a *App) shutdown(ctx context.Context) {
    // Save current position
    windowTitle := "Your App Window"
    a.windowPosMgr.SavePosition(ctx, windowTitle, a.windowPosPath)
}
```

### Multiple Windows

Use different window IDs (typically titles) to track multiple windows:

```go
// Main window
a.windowPosMgr.RestorePosition(ctx, "MyApp - Main")

// Settings window
a.windowPosMgr.RestorePosition(ctx, "MyApp - Settings")
```

### Storage Format

Positions are stored as JSON:

```json
{
  "MyApp - Main": {
    "x": 100,
    "y": 200,
    "width": 1024,
    "height": 768
  },
  "MyApp - Settings": {
    "x": 400,
    "y": 300,
    "width": 600,
    "height": 400
  }
}
```

## API Reference

### Types

```go
type WindowPositionManager struct {
    positions map[string]*WindowPosition
    xOffset   int // Platform-specific offset X
    yOffset   int // Platform-specific offset Y
}

type WindowPosition struct {
    X      int `json:"x"`
    Y      int `json:"y"`
    Width  int `json:"width"`
    Height int `json:"height"`
}
```

### Methods

#### `NewWindowPositionManager() *WindowPositionManager`

Creates a new window position manager instance.

#### `Load(storagePath string) error`

Loads saved window positions from disk. Creates directory if needed.

- Returns: `nil` if file doesn't exist (not an error)
- Returns: `error` on read/parse failures

#### `Save(storagePath string) error`

Saves current window positions to disk as JSON.

- Returns: `error` on write failures

#### `RestorePosition(ctx context.Context, windowID string)`

Restores window position and size for the given window ID.

- Platform-specific implementation in `windowposition_*.go`
- Gracefully handles missing positions (no-op)

#### `SavePosition(ctx context.Context, windowID string, storagePath string)`

Saves current window geometry for the given window ID.

- Platform-specific implementation in `windowposition_*.go`
- Reloads from disk first to preserve other windows' positions
- Skips save if dimensions are invalid (e.g., during shutdown)

#### `GetPosition(windowID string) *WindowPosition`

Returns saved position for a window ID, or `nil` if not found.

#### `SetPosition(windowID string, x, y, width, height int)`

Manually sets a position (doesn't save to disk).

## Linux Requirements

For reliable window position tracking on Linux, install xdotool:

**Debian/Ubuntu:**

```bash
sudo apt-get install xdotool
```

**openSUSE/SUSE:**

```bash
sudo zypper install xdotool
```

**Fedora/RHEL:**

```bash
sudo dnf install xdotool
```

**Arch Linux:**

```bash
sudo pacman -S xdotool
```

Without xdotool, the module will attempt to use Wails methods but may not be able to save positions if the window has been moved.

## Reusability

This module is designed to be copied into any Wails project:

1. Copy all `windowposition*.go` files to your project
2. Update `package main` to your package name if needed
3. Follow the usage pattern above
4. No modifications needed - platform detection is automatic via build tags

## Debugging

All operations log to stdout with `[WindowPos]` prefix:

```
[WindowPos] Restoring position for MyApp - X: 100 Y: 200 W: 1024 H: 768
[WindowPos] Detected offset - X: 8 Y: 31 - Compensating immediately
[WindowPos] Saving position for MyApp - X: 108 Y: 231 W: 1024 H: 768
```

On Linux without xdotool:

```
[WindowPos] Wails returned (0,0), trying xdotool fallback
[WindowPos] xdotool failed - install xdotool for Linux position tracking
[WindowPos] Run: sudo apt-get install xdotool
```

## License

Free to use in any project. No attribution required.
