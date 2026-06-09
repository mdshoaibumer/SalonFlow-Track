package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/salonflow/salonflow-track/internal/core/usecase"
)

// UpdateHandler handles update API endpoints.
type UpdateHandler struct {
	uc *usecase.UpdateUseCase
}

// NewUpdateHandler creates a new UpdateHandler.
func NewUpdateHandler(uc *usecase.UpdateUseCase) *UpdateHandler {
	return &UpdateHandler{uc: uc}
}

// Routes returns update routes.
func (h *UpdateHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/check", h.Check)
	r.Post("/download", h.Download)
	r.Post("/install", h.Install)
	r.Get("/history", h.History)
	r.Get("/status", h.Status)
	return r
}

// Check handles GET /update/check - checks for new versions.
func (h *UpdateHandler) Check(w http.ResponseWriter, r *http.Request) {
	status, err := h.uc.CheckForUpdate(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, apiResponse{Success: true, Data: status})
}

// Download handles POST /update/download - downloads the latest update.
func (h *UpdateHandler) Download(w http.ResponseWriter, r *http.Request) {
	record, err := h.uc.DownloadUpdate(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, apiResponse{Success: true, Data: record})
}

// Install handles POST /update/install - installs the downloaded update.
func (h *UpdateHandler) Install(w http.ResponseWriter, r *http.Request) {
	record, err := h.uc.InstallUpdate(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, apiResponse{Success: true, Data: record})
}

// History handles GET /update/history - returns update history.
func (h *UpdateHandler) History(w http.ResponseWriter, r *http.Request) {
	page := queryInt(r, "page", 1)
	perPage := queryInt(r, "per_page", 20)

	records, total, err := h.uc.ListHistory(r.Context(), page, perPage)
	if err != nil {
		writeError(w, err)
		return
	}
	totalPages := (total + perPage - 1) / perPage
	writeJSON(w, http.StatusOK, apiResponse{
		Success: true,
		Data:    records,
		Meta:    &apiMeta{Page: page, PerPage: perPage, Total: total, TotalPages: totalPages},
	})
}

// Status handles GET /update/status - returns current update status.
func (h *UpdateHandler) Status(w http.ResponseWriter, r *http.Request) {
	status, err := h.uc.GetUpdateStatus(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, apiResponse{Success: true, Data: status})
}
