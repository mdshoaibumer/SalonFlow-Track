package config

import (
	"os"
	"testing"
)

func TestLoad_ValidConfig(t *testing.T) {
	content := `
app:
  name: "Test App"
  environment: "test"
server:
  host: "127.0.0.1"
  port: "9090"
database:
  path: ":memory:"
  migration_dir: "testdata/migrations"
log:
  level: "debug"
  format: "text"
  output: "stdout"
`
	f, err := os.CreateTemp(t.TempDir(), "config-*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()

	cfg, err := Load(f.Name())
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	if cfg.App.Name != "Test App" {
		t.Errorf("expected app.name = 'Test App', got %q", cfg.App.Name)
	}
	if cfg.Server.Port != "9090" {
		t.Errorf("expected server.port = '9090', got %q", cfg.Server.Port)
	}
	if cfg.Database.Path != ":memory:" {
		t.Errorf("expected database.path = ':memory:', got %q", cfg.Database.Path)
	}
}

func TestLoad_Defaults(t *testing.T) {
	content := `
database:
  path: ":memory:"
`
	f, err := os.CreateTemp(t.TempDir(), "config-*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	f.WriteString(content)
	f.Close()

	cfg, err := Load(f.Name())
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	if cfg.Server.Port != "8080" {
		t.Errorf("expected default port 8080, got %q", cfg.Server.Port)
	}
	if cfg.Log.Level != "info" {
		t.Errorf("expected default log level 'info', got %q", cfg.Log.Level)
	}
	if cfg.App.Environment != "development" {
		t.Errorf("expected default environment 'development', got %q", cfg.App.Environment)
	}
}

func TestLoad_InvalidEnvironment(t *testing.T) {
	content := `
app:
  environment: "staging"
database:
  path: ":memory:"
`
	f, err := os.CreateTemp(t.TempDir(), "config-*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	f.WriteString(content)
	f.Close()

	_, err = Load(f.Name())
	if err == nil {
		t.Fatal("expected error for invalid environment")
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := Load("/nonexistent/config.yaml")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}
