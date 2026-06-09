package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"

	"github.com/salonflow/salonflow-track/internal/config"
	"github.com/salonflow/salonflow-track/internal/database"
	"github.com/salonflow/salonflow-track/internal/logger"
)

// App struct holds the Wails application context and backend services.
type App struct {
	ctx       context.Context
	cfg       *config.Config
	log       *slog.Logger
	db        *database.DB
	container *Container
	server    *http.Server
}

// NewApp creates a new App instance with all backend services initialized.
func NewApp() (*App, error) {
	configPath := resolveConfigPath()

	cfg, err := config.Load(configPath)
	if err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}
	cfg.App.Version = version

	log := logger.New(cfg.Log)
	log.Info("initializing SalonFlow Track desktop app",
		"version", version,
	)

	db, err := database.New(cfg.Database, log)
	if err != nil {
		return nil, fmt.Errorf("init database: %w", err)
	}

	applied, err := db.MigrateUp()
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("run migrations: %w", err)
	}
	if applied > 0 {
		log.Info("migrations applied", "count", applied)
	}

	container := NewContainer(cfg, log, db)

	// Start HTTP API server immediately (before Wails window opens)
	app := &App{
		cfg:       cfg,
		log:       log,
		db:        db,
		container: container,
	}
	app.startHTTPServer()

	return app, nil
}

// startHTTPServer starts the HTTP API server on the configured port.
func (a *App) startHTTPServer() {
	a.server = a.container.HTTPServer()

	// Find a free port if default is busy
	listener, err := net.Listen("tcp", a.server.Addr)
	if err != nil {
		// Try any available port
		listener, err = net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			a.log.Error("failed to start HTTP server", "error", err)
			return
		}
	}

	actualAddr := listener.Addr().String()
	a.log.Info("API server listening", "addr", actualAddr)

	go func() {
		if err := a.server.Serve(listener); err != nil && err != http.ErrServerClosed {
			a.log.Error("HTTP server error", "error", err)
		}
	}()
}

// startup is called when the Wails app starts. It saves the context.
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// shutdown is called when the Wails app is closing.
func (a *App) shutdown(ctx context.Context) {
	a.log.Info("shutting down desktop app...")

	if a.server != nil {
		if err := a.server.Shutdown(ctx); err != nil {
			a.log.Error("server shutdown error", "error", err)
		}
	}

	if a.db != nil {
		if err := a.db.Close(); err != nil {
			a.log.Error("database close error", "error", err)
		}
	}

	a.log.Info("shutdown complete")
}

// GetVersion returns the application version (exposed to frontend via Wails bindings).
func (a *App) GetVersion() string {
	return version
}
