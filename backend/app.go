package main

import (
	"context"
	"fmt"
	"log/slog"

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

	app := &App{
		cfg:       cfg,
		log:       log,
		db:        db,
		container: container,
	}

	return app, nil
}

// startup is called when the Wails app starts. It saves the context and propagates to all services.
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.container.SetContext(ctx)
}

// shutdown is called when the Wails app is closing.
func (a *App) shutdown(ctx context.Context) {
	a.log.Info("shutting down desktop app...")

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

// GetEnvironment returns the current environment.
func (a *App) GetEnvironment() string {
	return a.cfg.App.Environment
}
