# Automated Prerelease Builds

This folder contains **automatically generated prerelease builds** of SimpleAI that are created on every commit/push to the repository.

## ⚠️ Important Notice

These are **development builds** and may contain:

- Untested features
- Breaking changes
- Bugs and stability issues
- Incomplete functionality

**For stable releases**, please use the official releases from the [Releases page](../../../releases).

## Platform Support

Currently, automated prerelease builds are available for:

- ✅ **Windows** (64-bit)
- ✅ **Linux** (64-bit)

**Note for Linux users**: Install `xdotool` for window position persistence:

- Debian/Ubuntu: `sudo apt-get install xdotool`
- openSUSE/SUSE: `sudo zypper install xdotool`
- Fedora/RHEL: `sudo dnf install xdotool`
- Arch Linux: `sudo pacman -S xdotool`

Additional platforms (macOS) may be added in the future.

## Build Frequency

- **Trigger**: Every commit/push to the repository
- **Automation**: Builds are generated automatically via CI/CD pipeline
- **Version**: Each build includes the commit hash for traceability

## Versioning

Prerelease builds are named with the format: `SimpleAI <version>.PRE.exe` (Windows) or `SimpleAI-<version>.PRE` (Linux)

- Version number comes from `wails.json` → `info.productVersion`
- `.PRE` suffix indicates prerelease/development build
- Each build corresponds to a specific commit in the repository

## Usage

**Windows:**

1. Download the latest `SimpleAI <version>.PRE.exe` from this folder
2. Run the executable directly
3. Note: Builds are unsigned and may trigger Windows SmartScreen warnings

**Linux:**

1. Download the latest `SimpleAI-<version>.PRE` from this folder
2. Make executable: `chmod +x SimpleAI-<version>.PRE`
3. Run: `./SimpleAI-<version>.PRE`
4. Install xdotool for window position memory (see distribution-specific commands above)

## Feedback

If you encounter issues with prerelease builds, please report them with:

- The specific commit hash or build date
- Steps to reproduce the issue
- Your Windows version and system specifications

---

**Last Updated**: January 2026
