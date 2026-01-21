# SimpleAI Copilot Instructions

## Project Overview

SimpleAI is a **Wails-based desktop application** that provides a simple AI chatbot interface. It uses a Go backend with automatic method binding to JavaScript frontend, and Vite for frontend bundling.

**Key Stack:**

- Backend: Go 1.23 + Wails v2.10.2
- Frontend: Vanilla JavaScript + Vite + CSS
- Distribution: Cross-platform builds (Windows, Linux, macOS)

## Architecture

### Application Flow

SimpleAI works as a **smart launcher** for AI services:

1. **Launcher Page** - Shows buttons with custom title bar for navigation
2. **AI Service Navigation** - Clicking a service navigates the entire window to that AI site
3. **Window Position Persistence** - Saves and restores window positions across sessions
4. **Return Home** - Use window controls to reload and return to launcher

**Important Limitation**: When navigating to an AI service, the external site replaces ALL window content (including custom UI). This is unavoidable in Wails without using native Windows WebView2 COM API to create multiple webview instances.

### Why Not Iframes?

All AI services (ChatGPT, Claude, Copilot, Deepseek, Gemini, Grok, Meta AI, Perplexity) block iframe embedding via Content Security Policy (CSP) `frame-ancestors` directives. Attempting to use iframes results in CSP violations.

### Frontend-Backend Communication Pattern

The Wails framework **auto-generates TypeScript/JavaScript bindings** from Go methods:

```
Go Method (App struct) → Wails Binding → wailsjs/go/main/App.js → JavaScript
```

**Example:** `func (a *App) Greet(name string) string` in `app.go` automatically creates `Greet()` in `frontend/src/main.js`.

### Binding Conventions

1. **Only public Go methods are exposed** - Start method names with capital letters to bind them
2. **Receivers must use `*App` pointer** - Ensures app state access via `a.ctx`
3. **Context parameter** - Always receive `context.Context` in `startup(ctx context.Context)` and store in `a.ctx`
4. **Async calls in JS** - All bindings return Promises: `Greet(name).then(...).catch(...)`

### Directory Structure

```
├── main.go              # Wails app initialization & entry point
├── app.go               # App struct with bindable methods
├── wails.json           # Wails configuration
├── frontend/            # Vite frontend
│   ├── src/
│   │   ├── main.js      # Entry point, calls Go methods
│   │   ├── app.css      # Component styles
│   │   └── assets/      # Static images & fonts
│   ├── wailsjs/         # AUTO-GENERATED - Go bindings (don't edit)
│   └── package.json
└── build/               # Build outputs & platform configs
```

## Critical Development Workflows

### Local Development (Hot Reload)

**ALWAYS use the dev.ps1 script for development:**

```powershell
.\dev.ps1
```

- Automatically injects version from wails.json into binary
- Runs Vite dev server on port 34115 for browser access
- Provides fast hot reload for frontend changes
- Go code changes require manual restart

**Do not call `wails dev` directly** - the script ensures version is properly injected.

### Build Commands

**ALWAYS use the build.ps1 script for production builds:**

```powershell
.\build.ps1
```

This script:

- Reads version from wails.json
- Injects it into the binary via ldflags (`-X main.Version=x.x.x`)
- Builds for current platform (Windows by default)

For cross-platform builds, modify build.ps1 or use:

```powershell
$version = (Get-Content wails.json | ConvertFrom-Json).info.productVersion
wails build -platform linux/amd64 -ldflags "-X main.Version=$version"
```

**Platform Options:**

- `windows/amd64` - Windows 64-bit (default)
- `linux/amd64` - Linux 64-bit
- `darwin/arm64` - macOS ARM64 (Apple Silicon)
- `darwin/amd64` - macOS Intel 64-bit

Cross-compilation requires platform-specific toolchains:

- Windows on Linux/macOS: MinGW-w64
- Linux on Windows/macOS: GCC cross-compiler
- Check setup: `wails doctor`

Output goes to `build/bin/` (bundled with frontend from `frontend/dist/`)

## Key Patterns & Conventions

### Adding New Go Methods

1. Create public method in `app.go` struct with `*App` receiver
2. Return serializable types (strings, structs, ints, etc.)
3. Method auto-binds to `wailsjs/go/main/App.js` on rebuild
4. Call from JavaScript: `import {NewMethod} from '../wailsjs/go/main/App'`

```go
func (a *App) NewMethod(input string) (string, error) {
    // Use a.ctx for context operations
    return result, nil
}
```

### Frontend-Backend Error Handling

Use try-catch with Promise chains:

```javascript
try {
  MethodName(param)
    .then((result) => {
      /* handle result */
    })
    .catch((err) => console.error(err));
} catch (err) {
  console.error(err);
}
```

### Startup & Initialization

- Go: `startup(ctx context.Context)` called at app launch, store ctx in App struct
- JavaScript: `import` statements execute at module load
- Timing: Go methods are ready immediately in wailsjs bindings

## Build & Deployment

1. Frontend: `npm run build` (Vite compiles to `frontend/dist/`)
2. Backend: `wails build` embeds `frontend/dist/` via `//go:embed`
3. Result: Single executable with zero external dependencies

## Common Pitfalls to Avoid

- ❌ Don't export unexported Go methods (lowercase names won't bind)
- ❌ Don't mutate JavaScript bindings in `wailsjs/` - they regenerate
- ❌ Don't use `localStorage` for sensitive data - consider Go backend storage
- ❌ Don't build without running `npm run build` first (old frontend included)

## Testing & Debugging

- **Browser Dev Tools:** Run dev server, connect to `http://localhost:34115`, inspect Network/Console
- **Go Debug:** Use standard Go debugging with breakpoints in `app.go` and `main.go`
- **Frontend Build Errors:** Check `frontend/src/` and `wails.json` config
