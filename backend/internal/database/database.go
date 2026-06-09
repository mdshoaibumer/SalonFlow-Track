package database

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"

	"github.com/salonflow/salonflow-track/internal/config"
)

// DB manages the SQLite database lifecycle.
type DB struct {
	conn         *sql.DB
	log          *slog.Logger
	migrationDir string
}

// New opens a SQLite database connection and applies PRAGMA settings.
func New(cfg config.DatabaseConfig, log *slog.Logger) (*DB, error) {
	// Ensure the directory exists (skip for in-memory)
	if cfg.Path != ":memory:" {
		dir := filepath.Dir(cfg.Path)
		if dir != "" && dir != "." {
			if err := os.MkdirAll(dir, 0750); err != nil {
				return nil, fmt.Errorf("create db directory: %w", err)
			}
		}
	}

	dsn := cfg.Path
	if cfg.Path != ":memory:" {
		dsn = fmt.Sprintf("file:%s?_journal_mode=WAL&_busy_timeout=5000&_synchronous=NORMAL&_foreign_keys=ON&_cache_size=-64000", cfg.Path)
	} else {
		dsn = ":memory:?_foreign_keys=ON"
	}

	conn, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}

	conn.SetMaxOpenConns(cfg.MaxOpenConns)
	conn.SetMaxIdleConns(cfg.MaxOpenConns)

	if err := conn.Ping(); err != nil {
		conn.Close()
		return nil, fmt.Errorf("ping sqlite: %w", err)
	}

	log.Info("database connected", "path", cfg.Path)

	return &DB{
		conn:         conn,
		log:          log,
		migrationDir: cfg.MigrationDir,
	}, nil
}

// Conn returns the underlying *sql.DB connection.
func (db *DB) Conn() *sql.DB {
	return db.conn
}

// Close closes the database connection.
func (db *DB) Close() error {
	db.log.Info("closing database connection")
	return db.conn.Close()
}

// Ping checks database connectivity.
func (db *DB) Ping(ctx context.Context) error {
	return db.conn.PingContext(ctx)
}

// MigrateUp applies all pending migrations.
// Returns the number of migrations applied.
func (db *DB) MigrateUp() (int, error) {
	m, err := db.newMigrate()
	if err != nil {
		return 0, err
	}
	// Note: We intentionally do NOT call m.Close() here because golang-migrate's
	// Close() also closes the underlying *sql.DB connection that we share with
	// the rest of the application. The source (file reader) will be GC'd.

	beforeVer, _, _ := m.Version()

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return 0, fmt.Errorf("migrate up: %w", err)
	}

	afterVer, _, _ := m.Version()
	applied := int(afterVer) - int(beforeVer)
	if applied < 0 {
		applied = 0
	}

	return applied, nil
}

// MigrateDown rolls back the last migration.
func (db *DB) MigrateDown() error {
	m, err := db.newMigrate()
	if err != nil {
		return err
	}
	defer m.Close()

	if err := m.Steps(-1); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migrate down: %w", err)
	}
	return nil
}

// MigrateVersion returns the current migration version.
func (db *DB) MigrateVersion() (uint, bool, error) {
	m, err := db.newMigrate()
	if err != nil {
		return 0, false, err
	}
	defer m.Close()
	return m.Version()
}

func (db *DB) newMigrate() (*migrate.Migrate, error) {
	driver, err := sqlite3.WithInstance(db.conn, &sqlite3.Config{})
	if err != nil {
		return nil, fmt.Errorf("create migrate driver: %w", err)
	}

	absDir, err := filepath.Abs(db.migrationDir)
	if err != nil {
		return nil, fmt.Errorf("resolve migration dir: %w", err)
	}

	sourceURL := fmt.Sprintf("file://%s", filepath.ToSlash(absDir))

	m, err := migrate.NewWithDatabaseInstance(sourceURL, "sqlite3", driver)
	if err != nil {
		return nil, fmt.Errorf("create migrator: %w", err)
	}

	m.Log = &migrateLogger{log: db.log}
	return m, nil
}

// migrateLogger adapts slog to the migrate.Logger interface.
type migrateLogger struct {
	log *slog.Logger
}

func (l *migrateLogger) Printf(format string, v ...interface{}) {
	l.log.Info(fmt.Sprintf(format, v...), "component", "migrator")
}

func (l *migrateLogger) Verbose() bool {
	return false
}
