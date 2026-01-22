# SimpleAI v1.0.1

## What's New

### Added

- **Linux Prerelease Builds** - Automated Linux builds now available via GitHub Actions workflow
- **Multi-Distribution Support** - Installation instructions for multiple Linux distributions:
  - Debian/Ubuntu (apt)
  - openSUSE/SUSE (zypper)
  - Fedora/RHEL (dnf)
  - Arch Linux (pacman)
- **Automated Prereleases Documentation** - Added README to `automated-prereleases/` folder

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

## Installation

### Windows

Download `SimpleAI 1.0.1.PRE.exe` from the [automated-prereleases](https://github.com/chrilep/SimpleAI/tree/main/automated-prereleases) folder.

### Linux

Download `SimpleAI-1.0.1.PRE` from the [automated-prereleases](https://github.com/chrilep/SimpleAI/tree/main/automated-prereleases) folder.

**Linux Requirements**: Install xdotool for window position memory:

- Debian/Ubuntu: `sudo apt-get install xdotool`
- openSUSE/SUSE: `sudo zypper install xdotool`
- Fedora/RHEL: `sudo dnf install xdotool`
- Arch Linux: `sudo pacman -S xdotool`

## Full Changelog

See [CHANGELOG.md](https://github.com/chrilep/SimpleAI/blob/main/CHANGELOG.md) for complete details.
