package main

import (
	"context"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"

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

	var cfg *config.Config
	if configPath == "" {
		// No external config file — use built-in production defaults
		cfg = config.Default()
	} else {
		var err error
		cfg, err = config.Load(configPath)
		if err != nil {
			return nil, fmt.Errorf("load config: %w", err)
		}
	}
	cfg.App.Version = version

	// Resolve relative DB path to be next to the executable
	// so the app works regardless of the user's working directory
	if !filepath.IsAbs(cfg.Database.Path) && cfg.Database.Path != ":memory:" {
		exe, _ := os.Executable()
		cfg.Database.Path = filepath.Join(filepath.Dir(exe), cfg.Database.Path)
	}
	// Same for log file
	if cfg.Log.FilePath != "" && !filepath.IsAbs(cfg.Log.FilePath) {
		exe, _ := os.Executable()
		cfg.Log.FilePath = filepath.Join(filepath.Dir(exe), cfg.Log.FilePath)
	}

	log := logger.New(cfg.Log)
	log.Info("initializing SalonFlow Track desktop app",
		"version", version,
	)

	db, err := database.New(cfg.Database, log)
	if err != nil {
		return nil, fmt.Errorf("init database: %w", err)
	}

	// Use embedded migrations (compiled into the binary)
	migFS, err := fs.Sub(migrationsFS, "database/migrations")
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("embedded migrations: %w", err)
	}

	applied, err := db.MigrateUpFromFS(migFS)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("run migrations: %w", err)
	}
	if applied > 0 {
		log.Info("migrations applied", "count", applied)
	}

	container := NewContainer(cfg, log, db)

	// Ensure default admin user exists
	if err := container.authUC.EnsureDefaultAdmin(context.Background()); err != nil {
		log.Error("failed to ensure default admin", "error", err)
	}

	// Cleanup expired sessions on startup
	if err := container.authUC.CleanupExpiredSessions(context.Background()); err != nil {
		log.Warn("failed to cleanup expired sessions", "error", err)
	}

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
