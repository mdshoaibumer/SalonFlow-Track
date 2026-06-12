package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/salonflow/salonflow-track/internal/config"
	"gopkg.in/natefinch/lumberjack.v2"
)

// contextKey is used for storing values in context.
type contextKey string

const (
	CorrelationIDKey contextKey = "correlation_id"
	RequestIDKey     contextKey = "request_id"
	UserIDKey        contextKey = "user_id"
	UsernameKey      contextKey = "username"
	ModuleKey        contextKey = "module"
)

// MultiLogger manages multiple log outputs for different purposes.
type MultiLogger struct {
	App         *slog.Logger
	Error       *slog.Logger
	Audit       *slog.Logger
	Security    *slog.Logger
	Performance *slog.Logger
	logDir      string
}

// NewMultiLogger creates a structured multi-file logger system.
func NewMultiLogger(cfg config.LogConfig) *MultiLogger {
	logDir := cfg.FilePath
	if logDir == "" {
		logDir = "logs"
	}

	// Ensure log directory exists
	if err := os.MkdirAll(logDir, 0750); err != nil {
		// Fallback to stdout only
		defaultLogger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
		return &MultiLogger{
			App:         defaultLogger,
			Error:       defaultLogger,
			Audit:       defaultLogger,
			Security:    defaultLogger,
			Performance: defaultLogger,
			logDir:      logDir,
		}
	}

	level := parseLogLevel(cfg.Level)

	ml := &MultiLogger{logDir: logDir}
	ml.App = createLogger(filepath.Join(logDir, "app.log"), level, cfg)
	ml.Error = createLogger(filepath.Join(logDir, "error.log"), slog.LevelError, cfg)
	ml.Audit = createLogger(filepath.Join(logDir, "audit.log"), slog.LevelInfo, cfg)
	ml.Security = createLogger(filepath.Join(logDir, "security.log"), slog.LevelInfo, cfg)
	ml.Performance = createLogger(filepath.Join(logDir, "performance.log"), slog.LevelInfo, cfg)

	return ml
}

// LogDir returns the log directory path.
func (ml *MultiLogger) LogDir() string {
	return ml.logDir
}

func createLogger(path string, level slog.Level, cfg config.LogConfig) *slog.Logger {
	writer := &lumberjack.Logger{
		Filename:   path,
		MaxSize:    cfg.MaxSizeMB,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAgeDays,
		Compress:   cfg.Compress,
	}

	// Also write to stdout in development
	var w io.Writer
	if strings.ToLower(cfg.Output) == "both" {
		w = io.MultiWriter(writer, os.Stdout)
	} else {
		w = writer
	}

	opts := &slog.HandlerOptions{
		Level:     level,
		AddSource: level == slog.LevelDebug,
	}

	return slog.New(slog.NewJSONHandler(w, opts))
}

func parseLogLevel(s string) slog.Level {
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

// WithCorrelationID adds a correlation ID to the context.
func WithCorrelationID(ctx context.Context) context.Context {
	return context.WithValue(ctx, CorrelationIDKey, uuid.New().String())
}

// WithRequestID adds a request ID to the context.
func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, RequestIDKey, id)
}

// WithUser adds user info to the context.
func WithUser(ctx context.Context, userID, username string) context.Context {
	ctx = context.WithValue(ctx, UserIDKey, userID)
	ctx = context.WithValue(ctx, UsernameKey, username)
	return ctx
}

// WithModule adds module info to the context.
func WithModule(ctx context.Context, module string) context.Context {
	return context.WithValue(ctx, ModuleKey, module)
}

// GetCorrelationID retrieves the correlation ID from context.
func GetCorrelationID(ctx context.Context) string {
	if v, ok := ctx.Value(CorrelationIDKey).(string); ok {
		return v
	}
	return ""
}

// LogPerformance records a performance measurement.
func (ml *MultiLogger) LogPerformance(ctx context.Context, operation string, duration time.Duration, attrs ...slog.Attr) {
	baseAttrs := []slog.Attr{
		slog.String("operation", operation),
		slog.Duration("duration", duration),
		slog.String("duration_ms", fmt.Sprintf("%.2f", float64(duration.Nanoseconds())/1e6)),
	}

	if cid := GetCorrelationID(ctx); cid != "" {
		baseAttrs = append(baseAttrs, slog.String("correlation_id", cid))
	}
	if uid, ok := ctx.Value(UserIDKey).(string); ok {
		baseAttrs = append(baseAttrs, slog.String("user_id", uid))
	}

	baseAttrs = append(baseAttrs, attrs...)

	args := make([]any, 0, len(baseAttrs)*2)
	for _, attr := range baseAttrs {
		args = append(args, attr.Key, attr.Value.Any())
	}

	ml.Performance.Info("performance_metric", args...)
}

// LogSecurity records a security-related event.
func (ml *MultiLogger) LogSecurity(ctx context.Context, event string, level slog.Level, attrs ...slog.Attr) {
	baseAttrs := []slog.Attr{
		slog.String("event", event),
		slog.String("timestamp", time.Now().Format(time.RFC3339)),
	}

	if uid, ok := ctx.Value(UserIDKey).(string); ok {
		baseAttrs = append(baseAttrs, slog.String("user_id", uid))
	}
	if uname, ok := ctx.Value(UsernameKey).(string); ok {
		baseAttrs = append(baseAttrs, slog.String("username", uname))
	}

	baseAttrs = append(baseAttrs, attrs...)

	args := make([]any, 0, len(baseAttrs)*2)
	for _, attr := range baseAttrs {
		args = append(args, attr.Key, attr.Value.Any())
	}

	switch level {
	case slog.LevelWarn:
		ml.Security.Warn("security_event", args...)
	case slog.LevelError:
		ml.Security.Error("security_event", args...)
	default:
		ml.Security.Info("security_event", args...)
	}
}

// LogError records an error with stack trace information.
func (ml *MultiLogger) LogError(ctx context.Context, err error, module string, attrs ...slog.Attr) {
	baseAttrs := []slog.Attr{
		slog.String("module", module),
		slog.String("error", err.Error()),
		slog.String("timestamp", time.Now().Format(time.RFC3339)),
	}

	// Capture caller info
	_, file, line, ok := runtime.Caller(1)
	if ok {
		baseAttrs = append(baseAttrs, slog.String("source", fmt.Sprintf("%s:%d", file, line)))
	}

	if cid := GetCorrelationID(ctx); cid != "" {
		baseAttrs = append(baseAttrs, slog.String("correlation_id", cid))
	}
	if uid, ok := ctx.Value(UserIDKey).(string); ok {
		baseAttrs = append(baseAttrs, slog.String("user_id", uid))
	}

	baseAttrs = append(baseAttrs, attrs...)

	args := make([]any, 0, len(baseAttrs)*2)
	for _, attr := range baseAttrs {
		args = append(args, attr.Key, attr.Value.Any())
	}

	ml.Error.Error("application_error", args...)
}

// RecoverPanic catches and logs panics with full stack trace.
func (ml *MultiLogger) RecoverPanic(ctx context.Context, module string) {
	if r := recover(); r != nil {
		// Get stack trace
		buf := make([]byte, 4096)
		n := runtime.Stack(buf, false)
		stackTrace := string(buf[:n])

		ml.Error.Error("panic_recovered",
			"module", module,
			"panic", fmt.Sprintf("%v", r),
			"stack_trace", stackTrace,
			"timestamp", time.Now().Format(time.RFC3339),
		)

		ml.Security.Error("panic_detected",
			"module", module,
			"panic", fmt.Sprintf("%v", r),
		)
	}
}
