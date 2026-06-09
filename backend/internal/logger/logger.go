package logger

import (
	"io"
	"log/slog"
	"os"
	"strings"

	"github.com/salonflow/salonflow-track/internal/config"
	"gopkg.in/natefinch/lumberjack.v2"
)

// New creates a configured *slog.Logger from the application config.
func New(cfg config.LogConfig) *slog.Logger {
	level := parseLevel(cfg.Level)
	writer := buildWriter(cfg)

	opts := &slog.HandlerOptions{
		Level:     level,
		AddSource: level == slog.LevelDebug,
	}

	var h slog.Handler
	switch strings.ToLower(cfg.Format) {
	case "text":
		h = slog.NewTextHandler(writer, opts)
	default:
		h = slog.NewJSONHandler(writer, opts)
	}

	logger := slog.New(h)
	slog.SetDefault(logger)
	return logger
}

func buildWriter(cfg config.LogConfig) io.Writer {
	switch strings.ToLower(cfg.Output) {
	case "file":
		if cfg.FilePath == "" {
			return os.Stdout
		}
		return &lumberjack.Logger{
			Filename:   cfg.FilePath,
			MaxSize:    cfg.MaxSizeMB,
			MaxBackups: cfg.MaxBackups,
			MaxAge:     cfg.MaxAgeDays,
			Compress:   cfg.Compress,
		}
	case "stderr":
		return os.Stderr
	default:
		return os.Stdout
	}
}

func parseLevel(s string) slog.Level {
	switch strings.ToLower(s) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
