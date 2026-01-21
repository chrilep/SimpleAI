package main

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// Version is set during build from wails.json
var Version = "ersion dev"

// App struct
type App struct {
	ctx            context.Context
	startupService string
	windowPosMgr   *WindowPositionManager
	windowPosPath  string // Path to windows.json
}

// NewApp creates a new App application struct
func NewApp() *App {
	configDir, _ := os.UserConfigDir()

	return &App{
		windowPosMgr:  NewWindowPositionManager(),
		windowPosPath: filepath.Join(configDir, "SimpleAI", "windows.json"),
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	println("[DEBUG] Startup - Service:", a.startupService, "PID:", os.Getpid())
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
	runtime.WindowSetTitle(ctx, windowTitle)
	a.windowPosMgr.RestorePosition(ctx, windowTitle)
}

// shutdown is called when the app is about to quit
func (a *App) shutdown(ctx context.Context) {
	println("[DEBUG] Shutdown - Service:", a.startupService, "PID:", os.Getpid())
	a.windowPosMgr.SavePosition(ctx, a.GetWindowTitle(), a.windowPosPath)
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
	runtime.WindowReload(a.ctx)
}

// OpenNewInstance opens a new instance of the app with the specified service
func (a *App) OpenNewInstance(serviceName string) error {
	println("[DEBUG] OpenNewInstance called for:", serviceName)
	exePath, err := os.Executable()
	if err != nil {
		return err
	}

	cmd := exec.Command(exePath, serviceName)
	err = cmd.Start()
	if err != nil {
		return err
	}

	println("[DEBUG] New instance started with PID:", cmd.Process.Pid)
	return nil
}

// SaveWindowPositionManual allows manual saving of window position from frontend
func (a *App) SaveWindowPositionManual() error {
	a.windowPosMgr.SavePosition(a.ctx, a.GetWindowTitle(), a.windowPosPath)
	return nil
}
