package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/salonflow/salonflow-track/internal/core/usecase"
	"github.com/salonflow/salonflow-track/pkg/apperror"
)

// BackupHandler handles backup and restore API endpoints.
type BackupHandler struct {
	uc *usecase.BackupUseCase
}

// NewBackupHandler creates a new BackupHandler.
func NewBackupHandler(uc *usecase.BackupUseCase) *BackupHandler {
	return &BackupHandler{uc: uc}
}

// Routes returns backup routes.
func (h *BackupHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Post("/", h.Create)
	r.Get("/", h.List)
	r.Get("/stats", h.Stats)
	r.Get("/restores", h.ListRestores)
	r.Post("/{id}/verify", h.Verify)
	r.Post("/{id}/restore", h.Restore)
	r.Delete("/{id}", h.Delete)
	return r
}

// Create handles POST /backups - creates a new backup.
func (h *BackupHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		BackupType string `json:"backup_type"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	record, err := h.uc.CreateBackup(r.Context(), req.BackupType)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, apiResponse{Success: true, Data: record})
}

// List handles GET /backups - lists backup history.
func (h *BackupHandler) List(w http.ResponseWriter, r *http.Request) {
	page := queryInt(r, "page", 1)
	perPage := queryInt(r, "per_page", 20)

	backups, total, err := h.uc.ListBackups(r.Context(), page, perPage)
	if err != nil {
		writeError(w, err)
		return
	}

	totalPages := (total + perPage - 1) / perPage
	writeJSON(w, http.StatusOK, apiResponse{
		Success: true,
		Data:    backups,
		Meta:    &apiMeta{Page: page, PerPage: perPage, Total: total, TotalPages: totalPages},
	})
}

// Stats handles GET /backups/stats - returns backup stats.
func (h *BackupHandler) Stats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.uc.GetStats(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, apiResponse{Success: true, Data: stats})
}

// ListRestores handles GET /backups/restores - lists restore history.
func (h *BackupHandler) ListRestores(w http.ResponseWriter, r *http.Request) {
	page := queryInt(r, "page", 1)
	perPage := queryInt(r, "per_page", 20)

	restores, total, err := h.uc.ListRestores(r.Context(), page, perPage)
	if err != nil {
		writeError(w, err)
		return
	}

	totalPages := (total + perPage - 1) / perPage
	writeJSON(w, http.StatusOK, apiResponse{
		Success: true,
		Data:    restores,
		Meta:    &apiMeta{Page: page, PerPage: perPage, Total: total, TotalPages: totalPages},
	})
}

// Verify handles POST /backups/{id}/verify - verifies a backup.
func (h *BackupHandler) Verify(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		writeError(w, apperror.Validation("id", "invalid UUID"))
		return
	}

	result, err := h.uc.VerifyBackup(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, apiResponse{Success: true, Data: result})
}

// Restore handles POST /backups/{id}/restore - restores from a backup.
func (h *BackupHandler) Restore(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		writeError(w, apperror.Validation("id", "invalid UUID"))
		return
	}

	var req struct {
		Notes string `json:"notes"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	record, err := h.uc.RestoreBackup(r.Context(), id, req.Notes)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, apiResponse{Success: true, Data: record})
}

// Delete handles DELETE /backups/{id} - deletes a backup.
func (h *BackupHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		writeError(w, apperror.Validation("id", "invalid UUID"))
		return
	}

	if err := h.uc.DeleteBackup(r.Context(), id); err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, apiResponse{Success: true})
}
