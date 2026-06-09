package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/salonflow/salonflow-track/internal/core/domain"
	"github.com/salonflow/salonflow-track/internal/core/usecase"
	"github.com/salonflow/salonflow-track/pkg/apperror"
)

// ImportHandler handles import API endpoints.
type ImportHandler struct {
	uc        *usecase.ImportUseCase
	uploadDir string
}

// NewImportHandler creates a new ImportHandler.
func NewImportHandler(uc *usecase.ImportUseCase, uploadDir string) *ImportHandler {
	return &ImportHandler{uc: uc, uploadDir: uploadDir}
}

// Routes returns import routes.
func (h *ImportHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Post("/upload", h.Upload)
	r.Post("/validate", h.Validate)
	r.Post("/process", h.Process)
	r.Get("/history", h.History)
	r.Get("/{id}", h.GetJob)
	r.Get("/{id}/logs", h.Logs)
	return r
}

// Upload handles POST /import/upload - file upload and initial parsing.
func (h *ImportHandler) Upload(w http.ResponseWriter, r *http.Request) {
	// Max 50MB
	r.ParseMultipartForm(50 << 20)

	file, header, err := r.FormFile("file")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, apiResponse{Success: false, Error: &apiError{Code: "NO_FILE", Message: "File is required"}})
		return
	}
	defer file.Close()

	targetEntity := r.FormValue("target_entity")

	// Save file to upload dir
	if err := os.MkdirAll(h.uploadDir, 0750); err != nil {
		writeJSON(w, http.StatusInternalServerError, apiResponse{Success: false, Error: &apiError{Code: "STORAGE_ERROR", Message: "Failed to create upload directory"}})
		return
	}

	// Generate unique filename
	ext := filepath.Ext(header.Filename)
	destPath := filepath.Join(h.uploadDir, uuid.New().String()+ext)

	dst, err := os.Create(destPath)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, apiResponse{Success: false, Error: &apiError{Code: "STORAGE_ERROR", Message: "Failed to save file"}})
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		os.Remove(destPath)
		writeJSON(w, http.StatusInternalServerError, apiResponse{Success: false, Error: &apiError{Code: "STORAGE_ERROR", Message: "Failed to save file"}})
		return
	}

	job, headers, mappings, err := h.uc.Upload(r.Context(), header.Filename, destPath, targetEntity)
	if err != nil {
		os.Remove(destPath)
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, apiResponse{Success: true, Data: map[string]interface{}{
		"job":      job,
		"headers":  headers,
		"mappings": mappings,
	}})
}

// Validate handles POST /import/validate - validates data with column mapping.
func (h *ImportHandler) Validate(w http.ResponseWriter, r *http.Request) {
	var req struct {
		JobID    string                 `json:"job_id"`
		Mappings []domain.ColumnMapping `json:"mappings"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, apiResponse{Success: false, Error: &apiError{Code: "INVALID_JSON", Message: "Invalid request body"}})
		return
	}

	jobID, err := uuid.Parse(req.JobID)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, apiResponse{Success: false, Error: &apiError{Code: "INVALID_ID", Message: "Invalid job ID"}})
		return
	}

	preview, err := h.uc.Validate(r.Context(), jobID, req.Mappings)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, apiResponse{Success: true, Data: preview})
}

// Process handles POST /import/process - executes the import.
func (h *ImportHandler) Process(w http.ResponseWriter, r *http.Request) {
	var req struct {
		JobID string `json:"job_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, apiResponse{Success: false, Error: &apiError{Code: "INVALID_JSON", Message: "Invalid request body"}})
		return
	}

	jobID, err := uuid.Parse(req.JobID)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, apiResponse{Success: false, Error: &apiError{Code: "INVALID_ID", Message: "Invalid job ID"}})
		return
	}

	job, err := h.uc.Process(r.Context(), jobID)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, apiResponse{Success: true, Data: job})
}

// History handles GET /import/history - lists import jobs.
func (h *ImportHandler) History(w http.ResponseWriter, r *http.Request) {
	page := queryInt(r, "page", 1)
	perPage := queryInt(r, "per_page", 20)

	jobs, total, err := h.uc.ListJobs(r.Context(), page, perPage)
	if err != nil {
		writeError(w, err)
		return
	}

	totalPages := (total + perPage - 1) / perPage
	writeJSON(w, http.StatusOK, apiResponse{
		Success: true,
		Data:    jobs,
		Meta:    &apiMeta{Page: page, PerPage: perPage, Total: total, TotalPages: totalPages},
	})
}

// GetJob handles GET /import/{id} - returns a specific job.
func (h *ImportHandler) GetJob(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, apiResponse{Success: false, Error: &apiError{Code: "INVALID_ID", Message: "Invalid job ID"}})
		return
	}

	job, err := h.uc.GetJob(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, apiResponse{Success: true, Data: job})
}

// Logs handles GET /import/{id}/logs - returns logs for a job.
func (h *ImportHandler) Logs(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, apiResponse{Success: false, Error: &apiError{Code: "INVALID_ID", Message: "Invalid job ID"}})
		return
	}

	status := r.URL.Query().Get("status")
	page := queryInt(r, "page", 1)
	perPage := queryInt(r, "per_page", 50)

	logs, total, err := h.uc.ListLogs(r.Context(), id, status, page, perPage)
	if err != nil {
		writeError(w, err)
		return
	}

	totalPages := (total + perPage - 1) / perPage
	writeJSON(w, http.StatusOK, apiResponse{
		Success: true,
		Data:    logs,
		Meta:    &apiMeta{Page: page, PerPage: perPage, Total: total, TotalPages: totalPages},
	})
}

// Ensure the handler uses apperror (avoid unused import warning).
var _ = apperror.NotFound
