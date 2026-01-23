//go:build darwin
// +build darwin

package main

import (
	"fmt"
	"os/exec"
	"strings"
)

// findAndActivateWindow searches for a window with the given title and activates it
// Returns true if window was found and activated, false otherwise
func findAndActivateWindow(windowTitle string) (bool, error) {
	// AppleScript to find and activate window by title
	script := fmt.Sprintf(`
		tell application "System Events"
			set foundWindow to false
			repeat with proc in (every process whose visible is true)
				tell proc
					if exists (windows whose name is "%s") then
						set frontmost to true
						set foundWindow to true
						exit repeat
					end if
				end tell
			end repeat
			return foundWindow
		end tell
	`, windowTitle)

	cmd := exec.Command("osascript", "-e", script)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return false, err
	}

	return strings.Contains(string(output), "true"), nil
}
