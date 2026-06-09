package database

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/salonflow/salonflow-track/internal/config"
)

func TestNew_InMemory(t *testing.T) {
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelWarn}))

	cfg := config.DatabaseConfig{
		Path:         ":memory:",
		MigrationDir: "testdata/migrations",
		MaxOpenConns: 1,
	}

	db, err := New(cfg, log)
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}
	defer db.Close()

	if err := db.Ping(context.Background()); err != nil {
		t.Fatalf("Ping() error: %v", err)
	}
}

func TestNew_ForeignKeysEnabled(t *testing.T) {
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelWarn}))

	cfg := config.DatabaseConfig{
		Path:         ":memory:",
		MigrationDir: "testdata/migrations",
		MaxOpenConns: 1,
	}

	db, err := New(cfg, log)
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}
	defer db.Close()

	var fkEnabled int
	err = db.Conn().QueryRow("PRAGMA foreign_keys").Scan(&fkEnabled)
	if err != nil {
		t.Fatalf("PRAGMA foreign_keys error: %v", err)
	}
	if fkEnabled != 1 {
		t.Error("expected foreign_keys to be enabled")
	}
}
