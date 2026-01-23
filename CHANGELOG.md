# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.2.0] - 2026-01-23

### Added

- **Window Activation & Deduplication** - Prevents duplicate AI service windows
  - Automatically detects and activates existing windows when opening the same service
  - Platform-specific implementation: PowerShell (Windows), wmctrl/xdotool (Linux), AppleScript (macOS)
  - Asynchronous window search with 2-second timeout to prevent UI blocking
  - Restores minimized windows and brings them to foreground
  - Falls back to opening new instance if existing window not found or on timeout
- **Screen Bounds Validation** - Windows are automatically kept within visible screen boundaries
  - Prevents windows from being positioned outside the visible display area
  - Ensures at least 20 pixels of the window remain visible on screen
  - Automatically corrects saved positions when restoring windows
  - Protects against lost windows after monitor disconnection or resolution changes
  - Cross-platform support via `modWindowMemory` module (Windows, Linux, macOS)
  - Intelligent validation using `runtime.ScreenGetAll()` for primary screen dimensions
- **File Locking for Window Positions** - Prevents race conditions in multi-instance setups
  - Platform-specific file locking: `LockFileEx` (Windows), `flock` (Linux), `flock` (macOS)
  - Retry mechanism with 50ms delays for concurrent write attempts
  - Protects `windows.json` from corruption when multiple instances save simultaneously
  - Separate platform files: `filelock_windows.go`, `filelock_linux.go`, `filelock_darwin.go`
- **Service Navigation Architecture** - Enhanced launcher with full-window service integration
  - Entire window navigates to selected AI service for seamless experience
  - Replaces iframe approach (all AI services block iframes via CSP `frame-ancestors` directives)
  - Direct service URLs for ChatGPT, Claude, Copilot, Deepseek, Gemini, Grok, Meta AI, and Perplexity
  - No custom UI overlay possible due to external site control, providing native AI service experience
  - Return to launcher via window reload/navigation controls
- **Updated Application Icons** - Refreshed visual identity
  - New optimized PNG icons (397KB → 376KB)
  - New optimized ICO icons (340KB → 281KB)
  - Better compression while maintaining visual quality

### Changed

- **Debug Output Control** - Configurable debug logging
  - Debug prints now controlled by `dbg` flag in each function
  - Set to `false` by default for cleaner production logs
  - Easy to enable per-function for troubleshooting
- **Linux Prerelease Build Workflow** - Enhanced automation
  - AppImage format for better compatibility across distributions
  - Improved build process with proper error handling
- **modWindowMemory Module Documentation** - Enhanced README
  - Added file locking documentation
  - Updated usage examples with screen bounds validation
  - Better platform-specific implementation notes

## [1.1.0] - 2026-01-22

### Added

- **Persistent WebView Storage** - Cookies, sessions, and cache now persist across application restarts
  - WebView data stored in system cache directory: `%LOCALAPPDATA%\SimpleAI\webview` (Windows) or `~/.cache/SimpleAI/webview` (Linux/macOS)
  - Maintains login sessions for all AI services between launches
  - Separate storage per AI service via command-line arguments

### Fixed

- **Linux Window Position Bug** - Fixed 10x10 pixel window position issue on SUSE Linux and other distributions
  - Enhanced debugging for window position detection on Linux platforms
  - Improved xdotool fallback mechanism for reliable position tracking

## [1.0.1] - 2026-01-22

### Added

- **Linux Prerelease Builds** - Automated Linux builds now available via GitHub Actions workflow
- **Automated Prerelease Workflow** - Both Windows and Linux binaries are automatically built and committed on every push
- **Automated Prereleases Documentation** - Added README to `automated-prereleases/` folder explaining prerelease build system

### Changed

- **Window Position Module Refactored** - Extracted window position persistence into standalone `modWindowMemory` module
  - Platform-specific implementations for Windows, Linux (with xdotool fallback), and macOS
  - Full documentation in `modWindowMemory/README.md`
  - Module is now portable and reusable across Wails projects
  - Linux support improved with xdotool X11 fallback for reliable position tracking

### Fixed

- **Linux Window Position Tracking** - Implemented xdotool-based fallback for GTK window geometry issues
  - Wails `WindowGetPosition()` on Linux/GTK often returns (0,0) - now handled gracefully
  - Requires `xdotool` package for full functionality on Linux

## [1.0.0] - 2026-01-21

### Added

- **Multi-Service Launcher** - Quick access to 8 major AI chatbot services:
  - ChatGPT (OpenAI GPT-4/5)
  - Claude (Anthropic Sonnet)
  - Copilot (Microsoft)
  - Deepseek (V3.2)
  - Gemini (Google 2.0/2.5)
  - Grok (X/Twitter)
  - Meta AI (LLaMA)
  - Perplexity (AI Research)
- **Frameless Custom Title Bar** - Modern UI for launcher window with minimize/close controls
- **Window Position Persistence** - Automatically saves and restores window positions and sizes across sessions
- **Multi-Instance Support** - Open multiple AI services simultaneously in separate windows
- **Service Information Modals** - Info button (?) on each service showing detailed descriptions and use cases
- **Responsive Layout** - Grid layout adapts to window size with fixed service button dimensions
- **Drag-to-Move** - Entire launcher window draggable via custom title bar
- **Cross-Platform Support** - Runs on Windows, macOS, and Linux
- **Auto-Generated Bindings** - Wails framework automatically generates JavaScript bindings from Go methods
- **Development Scripts** - `dev.ps1` for hot reload development, `build.ps1` for production builds
- **Window Position Manager** - Reusable module with automatic platform offset compensation for Windows titlebar/borders
- **Race Condition Fix** - Multi-instance safe window position saving with file reload before write

### Technical Details

- Built with Wails v2.10.2 + Go 1.23
- Frontend: Vanilla JavaScript + Vite 3.x + CSS
- Rendering: WebView2 (Windows), WebKit (macOS/Linux)
- Storage: JSON file for window positions
- License: GNU AGPL-3.0

[1.0.0]: https://github.com/chrilep/SimpleAI/releases/tag/v1.0.0
