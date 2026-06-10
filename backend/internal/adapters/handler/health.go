package handler

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"runtime"
	"time"

	"github.com/salonflow/salonflow-track/internal/config"
	"github.com/salonflow/salonflow-track/internal/database"
)

// HealthHandler serves the /health endpoint.
type HealthHandler struct {
	db        *database.DB
	cfg       *config.Config
	log       *slog.Logger
	startedAt time.Time
}

// NewHealthHandler creates a new HealthHandler.
func NewHealthHandler(db *database.DB, cfg *config.Config, log *slog.Logger) *HealthHandler {
	return &HealthHandler{
		db:        db,
		cfg:       cfg,
		log:       log,
		startedAt: time.Now(),
	}
}

// HealthResponse is the JSON payload returned by the health check.
type HealthResponse struct {
	Status         string `json:"status"`
	Version        string `json:"version"`
	Environment    string `json:"environment"`
	DatabaseStatus string `json:"database_status"`
	Uptime         string `json:"uptime"`
	GoVersion      string `json:"go_version"`
}

// Check handles GET /api/v1/health.
func (h *HealthHandler) Check(w http.ResponseWriter, r *http.Request) {
	dbStatus := "healthy"
	if err := h.db.Ping(context.Background()); err != nil {
		dbStatus = "unhealthy"
		h.log.Warn("health check: database unhealthy", "error", err)
	}

	status := "healthy"
	if dbStatus != "healthy" {
		status = "degraded"
	}

	resp := HealthResponse{
		Status:         status,
		Version:        h.cfg.App.Version,
		Environment:    h.cfg.App.Environment,
		DatabaseStatus: dbStatus,
		Uptime:         time.Since(h.startedAt).Round(time.Second).String(),
		GoVersion:      runtime.Version(),
	}

	httpStatus := http.StatusOK
	if status != "healthy" {
		httpStatus = http.StatusServiceUnavailable
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(httpStatus)
	_ = json.NewEncoder(w).Encode(resp)
}
