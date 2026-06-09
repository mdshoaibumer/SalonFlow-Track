package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/salonflow/salonflow-track/internal/core/domain"
	"github.com/salonflow/salonflow-track/internal/core/usecase"
)

// CloudBackupHandler handles cloud backup API endpoints.
type CloudBackupHandler struct {
	uc *usecase.CloudBackupUseCase
}

// NewCloudBackupHandler creates a new CloudBackupHandler.
func NewCloudBackupHandler(uc *usecase.CloudBackupUseCase) *CloudBackupHandler {
	return &CloudBackupHandler{uc: uc}
}

// Routes returns cloud backup routes.
func (h *CloudBackupHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/config", h.GetConfig)
	r.Post("/config", h.SaveConfig)
	r.Post("/test", h.TestConnection)
	r.Post("/backup", h.BackupNow)
	r.Post("/restore/{id}", h.Restore)
	r.Get("/history", h.ListHistory)
	r.Get("/stats", h.GetStats)
	return r
}

// GetConfig handles GET /cloud-backup/config.
func (h *CloudBackupHandler) GetConfig(w http.ResponseWriter, r *http.Request) {
	cfg, err := h.uc.GetConfig(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, apiResponse{Success: true, Data: cfg})
}

// SaveConfig handles POST /cloud-backup/config.
func (h *CloudBackupHandler) SaveConfig(w http.ResponseWriter, r *http.Request) {
	var cfg domain.CloudBackupConfig
	if err := json.NewDecoder(r.Body).Decode(&cfg); err != nil {
		writeJSON(w, http.StatusBadRequest, apiResponse{Success: false, Error: &apiError{Code: "INVALID_JSON", Message: "Invalid request body"}})
		return
	}
	if err := h.uc.SaveConfig(r.Context(), &cfg); err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, apiResponse{Success: true, Data: cfg})
}

// TestConnection handles POST /cloud-backup/test.
func (h *CloudBackupHandler) TestConnection(w http.ResponseWriter, r *http.Request) {
	if err := h.uc.TestConnection(r.Context()); err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, apiResponse{Success: true, Data: map[string]string{"message": "Connection successful"}})
}

// BackupNow handles POST /cloud-backup/backup.
func (h *CloudBackupHandler) BackupNow(w http.ResponseWriter, r *http.Request) {
	history, err := h.uc.BackupNow(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, apiResponse{Success: true, Data: history})
}

// Restore handles POST /cloud-backup/restore/{id}.
func (h *CloudBackupHandler) Restore(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, apiResponse{Success: false, Error: &apiError{Code: "INVALID_ID", Message: "Invalid ID"}})
		return
	}
	if err := h.uc.Restore(r.Context(), id); err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, apiResponse{Success: true, Data: map[string]string{"message": "Restore initiated"}})
}

// ListHistory handles GET /cloud-backup/history.
func (h *CloudBackupHandler) ListHistory(w http.ResponseWriter, r *http.Request) {
	limit := 20
	offset := 0
	if l := r.URL.Query().Get("per_page"); l != "" {
		if v, _ := strconv.Atoi(l); v > 0 {
			limit = v
		}
	}
	if p := r.URL.Query().Get("page"); p != "" {
		if v, _ := strconv.Atoi(p); v > 1 {
			offset = (v - 1) * limit
		}
	}

	history, total, err := h.uc.ListHistory(r.Context(), limit, offset)
	if err != nil {
		writeError(w, err)
		return
	}
	page := offset/limit + 1
	totalPages := (total + limit - 1) / limit
	writeJSON(w, http.StatusOK, apiResponse{
		Success: true,
		Data:    history,
		Meta:    &apiMeta{Page: page, PerPage: limit, Total: total, TotalPages: totalPages},
	})
}

// GetStats handles GET /cloud-backup/stats.
func (h *CloudBackupHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.uc.GetStats(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, apiResponse{Success: true, Data: stats})
}
