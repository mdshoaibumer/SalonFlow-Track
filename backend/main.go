package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

const version = "0.2.0"

//go:embed all:dist
var assets embed.FS

func main() {
	if err := run(); err != nil {
		slog.Error("application failed", "error", err)
		os.Exit(1)
	}
}

func run() error {
	app, err := NewApp()
	if err != nil {
		return fmt.Errorf("init app: %w", err)
	}

	// Strip the "dist" prefix so index.html is at the root of the FS
	distFS, err := fs.Sub(assets, "dist")
	if err != nil {
		return fmt.Errorf("fs.Sub: %w", err)
	}

	err = wails.Run(&options.App{
		Title:     "SalonFlow Track",
		Width:     1400,
		Height:    900,
		MinWidth:  1024,
		MinHeight: 700,
		AssetServer: &assetserver.Options{
			Assets: distFS,
		},
		OnStartup:  app.startup,
		OnShutdown: app.shutdown,
		Bind:       append([]interface{}{app}, app.container.Bindings()...),
		Windows: &windows.Options{
			WebviewIsTransparent:              false,
			WindowIsTranslucent:               false,
			DisableFramelessWindowDecorations: false,
			Theme:                             windows.SystemDefault,
		},
	})

	if err != nil {
		return fmt.Errorf("wails run: %w", err)
	}
	return nil
}

func resolveConfigPath() string {
	// Check environment variable first
	if p := os.Getenv("SALONFLOW_CONFIG"); p != "" {
		return p
	}
	// Check current directory
	if _, err := os.Stat("config.yaml"); err == nil {
		return "config.yaml"
	}
	// Check next to executable
	exe, _ := os.Executable()
	exeDir := filepath.Join(filepath.Dir(exe), "config.yaml")
	if _, err := os.Stat(exeDir); err == nil {
		return exeDir
	}
	// No config file found — will use built-in defaults
	return ""
}
