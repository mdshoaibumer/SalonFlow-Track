package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config is the top-level application configuration.
type Config struct {
	App      AppConfig      `yaml:"app"`
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Log      LogConfig      `yaml:"log"`
}

// AppConfig holds application metadata.
type AppConfig struct {
	Name        string `yaml:"name"`
	Version     string `yaml:"-"` // Set at build time
	Environment string `yaml:"environment"`
}

// ServerConfig holds HTTP server settings.
type ServerConfig struct {
	Host               string `yaml:"host"`
	Port               string `yaml:"port"`
	ReadTimeoutSec     int    `yaml:"read_timeout_sec"`
	WriteTimeoutSec    int    `yaml:"write_timeout_sec"`
	ShutdownTimeoutSec int    `yaml:"shutdown_timeout_sec"`
}

// ShutdownTimeout returns the shutdown timeout as a Duration.
func (s ServerConfig) ShutdownTimeout() time.Duration {
	if s.ShutdownTimeoutSec <= 0 {
		return 10 * time.Second
	}
	return time.Duration(s.ShutdownTimeoutSec) * time.Second
}

// DatabaseConfig holds SQLite settings.
type DatabaseConfig struct {
	Path         string `yaml:"path"`
	MigrationDir string `yaml:"migration_dir"`
	MaxOpenConns int    `yaml:"max_open_conns"`
}

// LogConfig holds logging settings.
type LogConfig struct {
	Level      string `yaml:"level"`
	Format     string `yaml:"format"`
	Output     string `yaml:"output"`
	FilePath   string `yaml:"file_path"`
	MaxSizeMB  int    `yaml:"max_size_mb"`
	MaxBackups int    `yaml:"max_backups"`
	MaxAgeDays int    `yaml:"max_age_days"`
	Compress   bool   `yaml:"compress"`
}

// Load reads and validates configuration from a YAML file.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config %s: %w", path, err)
	}

	expanded := os.ExpandEnv(string(data))

	var cfg Config
	if err := yaml.Unmarshal([]byte(expanded), &cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	cfg.setDefaults()
	if err := cfg.validate(); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (c *Config) setDefaults() {
	if c.App.Name == "" {
		c.App.Name = "SalonFlow Track"
	}
	if c.App.Environment == "" {
		c.App.Environment = "development"
	}
	if c.Server.Host == "" {
		c.Server.Host = "127.0.0.1"
	}
	if c.Server.Port == "" {
		c.Server.Port = "8080"
	}
	if c.Server.ReadTimeoutSec == 0 {
		c.Server.ReadTimeoutSec = 30
	}
	if c.Server.WriteTimeoutSec == 0 {
		c.Server.WriteTimeoutSec = 30
	}
	if c.Server.ShutdownTimeoutSec == 0 {
		c.Server.ShutdownTimeoutSec = 10
	}
	if c.Database.Path == "" {
		c.Database.Path = "data/salonflow.db"
	}
	if c.Database.MigrationDir == "" {
		c.Database.MigrationDir = "migrations"
	}
	if c.Database.MaxOpenConns == 0 {
		c.Database.MaxOpenConns = 1
	}
	if c.Log.Level == "" {
		c.Log.Level = "info"
	}
	if c.Log.Format == "" {
		c.Log.Format = "json"
	}
	if c.Log.Output == "" {
		c.Log.Output = "stdout"
	}
	if c.Log.MaxSizeMB == 0 {
		c.Log.MaxSizeMB = 100
	}
	if c.Log.MaxBackups == 0 {
		c.Log.MaxBackups = 3
	}
	if c.Log.MaxAgeDays == 0 {
		c.Log.MaxAgeDays = 28
	}
}

func (c *Config) validate() error {
	validEnvs := map[string]bool{"development": true, "production": true, "test": true}
	if !validEnvs[c.App.Environment] {
		return fmt.Errorf("config: invalid environment %q (must be development, production, or test)", c.App.Environment)
	}
	if c.Database.Path == "" {
		return fmt.Errorf("config: database.path is required")
	}
	return nil
}
