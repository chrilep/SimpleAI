package main

import (
	"context"
	"embed"
	"os"
	"path/filepath"
	"strings"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/linux"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Check for command-line arguments
	var startupService string
	if len(os.Args) > 1 {
		startupService = strings.ToLower(os.Args[1])
	}

	// Create an instance of the app structure
	app := NewApp()
	app.startupService = startupService

	// Launcher gets frameless window for custom title bar
	frameless := startupService == ""

	// Get user data directory for WebView storage (cookies, sessions, cache)
	userDataDir, err := os.UserCacheDir()
	if err != nil {
		// Fallback to config dir if cache dir fails
		userDataDir, _ = os.UserConfigDir()
	}
	webviewDataPath := filepath.Join(userDataDir, "SimpleAI", "webview")

	// Create application with options
	err = wails.Run(&options.App{
		Title:     "SimpleAI",
		Width:     1024,
		Height:    768,
		MinWidth:  160,
		MinHeight: 50,
		Frameless: frameless,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		OnShutdown:       app.shutdown,
		OnBeforeClose: func(ctx context.Context) bool {
			// Save window position
			app.windowPosMgr.SavePosition(ctx, app.GetWindowTitle(), app.windowPosPath)

			// Save current URL if on an AI service page
			// Note: We can't execute JavaScript in external sites due to CSP,
			// but the URL will be saved on the next launch when user returns
			return false
		},
		Bind: []interface{}{
			app,
		},
		Windows: &windows.Options{
			WebviewUserDataPath: webviewDataPath,
			WebviewBrowserPath:  "",
		},
		Linux: &linux.Options{
			WebviewGpuPolicy: linux.WebviewGpuPolicyOnDemand,
			ProgramName:      "SimpleAI",
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
