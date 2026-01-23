package main

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"SimpleAI/modWindowMemory"

	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

// Version is set during build from wails.json
var Version = "ersion dev"

// App struct
type App struct {
	ctx            context.Context
	startupService string
	windowPosMgr   *modWindowMemory.WindowPositionManager
	windowPosPath  string // Path to windows.json
}

// NewApp creates a new App application struct
func NewApp() *App {
	configDir, _ := os.UserConfigDir()

	return &App{
		windowPosMgr:  modWindowMemory.NewWindowPositionManager(),
		windowPosPath: filepath.Join(configDir, "SimpleAI", "windows.json"),
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	dbg := false // Set to true for debug output
	a.ctx = ctx
	if dbg {
		println("[DEBUG] Startup - Service:", a.startupService, "PID:", os.Getpid())
	}

	a.windowPosMgr.Load(a.windowPosPath)

	// Get window title for position restore
	windowTitle := "SimpleAI"
	if a.startupService != "" {
		// Find service label
		serviceLabels := map[string]string{
			"chatgpt":    "ChatGPT",
			"claude":     "Claude (Sonnet)",
			"copilot":    "Copilot",
			"deepseek":   "Deepseek",
			"gemini":     "Gemini",
			"grok":       "Grok",
			"meta":       "Meta AI",
			"perplexity": "Perplexity",
		}
		if label, ok := serviceLabels[a.startupService]; ok {
			windowTitle = "SimpleAI - " + label
		}
	}
	wailsRuntime.WindowSetTitle(ctx, windowTitle)
	a.windowPosMgr.RestorePosition(ctx, windowTitle)
}

// shutdown is called when the app is about to quit
func (a *App) shutdown(ctx context.Context) {
	dbg := false // Set to true for debug output
	if dbg {
		println("[DEBUG] Shutdown - Service:", a.startupService, "PID:", os.Getpid())
	}
	// Note: Window position is already saved in OnBeforeClose hook (main.go)
	// Don't save here as window may already be destroyed
}

// GetStartupService returns the service name to navigate to on startup
func (a *App) GetStartupService() string {
	return a.startupService
}

// GetVersion returns the application version
func (a *App) GetVersion() string {
	return Version
}

// GetWindowTitle returns the current window title based on startup service
func (a *App) GetWindowTitle() string {
	windowTitle := "SimpleAI"
	if a.startupService != "" {
		serviceLabels := map[string]string{
			"chatgpt":    "ChatGPT",
			"claude":     "Claude (Sonnet)",
			"copilot":    "Copilot",
			"deepseek":   "Deepseek",
			"gemini":     "Gemini",
			"grok":       "Grok",
			"meta":       "Meta AI",
			"perplexity": "Perplexity",
		}
		if label, ok := serviceLabels[a.startupService]; ok {
			windowTitle = "SimpleAI - " + label
		}
	}
	return windowTitle
}

// GoHome navigates back to the launcher page
func (a *App) GoHome() {
	wailsRuntime.WindowReload(a.ctx)
}

// OpenNewInstance opens a new instance of the app with the specified service
// or activates an existing window if one is already open
func (a *App) OpenNewInstance(serviceName string) error {
	dbg := false // Set to true for debug output
	if dbg {
		println("[DEBUG] OpenNewInstance called for:", serviceName)
	}

	// Get window title for this service
	serviceLabels := map[string]string{
		"chatgpt":    "ChatGPT",
		"claude":     "Claude (Sonnet)",
		"copilot":    "Copilot",
		"deepseek":   "Deepseek",
		"gemini":     "Gemini",
		"grok":       "Grok",
		"meta":       "Meta AI",
		"perplexity": "Perplexity",
	}

	windowTitle := "SimpleAI - " + serviceLabels[serviceName]
	if dbg {
		println("[DEBUG] Trying to find an existing window with title:", windowTitle)
	}

	// Try to find and activate existing window asynchronously
	// This prevents blocking the UI while searching for windows
	resultChan := make(chan struct {
		found bool
		err   error
	}, 1)

	go func() {
		found, err := findAndActivateWindow(windowTitle)
		resultChan <- struct {
			found bool
			err   error
		}{found, err}
	}()

	// Wait for result with timeout
	select {
	case result := <-resultChan:
		if result.err != nil {
			if dbg {
				println("[DEBUG] Error searching for window:", result.err)
			}
			// Continue to open new instance on error
		}
		if result.found {
			if dbg {
				println("[DEBUG] Found and activated existing window:", windowTitle)
			}
			return nil
		}
	case <-time.After(2 * time.Second):
		// Timeout - proceed to open new instance
		if dbg {
			println("[DEBUG] Window search timed out, opening new instance")
		}
	}

	// No existing window found, start new instance
	if dbg {
		println("[DEBUG] No existing window found, starting new instance")
	}
	exePath, err := os.Executable()
	if err != nil {
		return err
	}

	cmd := exec.Command(exePath, serviceName)
	err = cmd.Start()
	if err != nil {
		return err
	}

	if dbg {
		println("[DEBUG] Started new instance with PID:", cmd.Process.Pid)
	}
	return nil
}

// SaveWindowPositionManual allows manual saving of window position from frontend
func (a *App) SaveWindowPositionManual() error {
	a.windowPosMgr.SavePosition(a.ctx, a.GetWindowTitle(), a.windowPosPath)
	return nil
}

// findAndActivateWindow is implemented in platform-specific files:
// - app_windows.go: Uses native Windows API for fast window search
// - app_linux.go: Uses wmctrl or xdotool
// - app_darwin.go: Uses osascript with AppleScript
