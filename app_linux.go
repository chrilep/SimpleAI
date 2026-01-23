//go:build linux
// +build linux

package main

import (
	"os/exec"
)

// findAndActivateWindow searches for a window with the given title and activates it
// Returns true if window was found and activated, false otherwise
func findAndActivateWindow(windowTitle string) (bool, error) {
	// Try wmctrl first
	cmd := exec.Command("wmctrl", "-a", windowTitle)
	err := cmd.Run()
	if err == nil {
		return true, nil
	}

	// Fallback to xdotool
	cmd = exec.Command("xdotool", "search", "--name", windowTitle, "windowactivate")
	err = cmd.Run()
	if err == nil {
		return true, nil
	}

	// Neither tool worked or window not found
	return false, nil
}
