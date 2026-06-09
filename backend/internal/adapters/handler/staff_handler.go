package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/salonflow/salonflow-track/internal/core/usecase"
	"github.com/salonflow/salonflow-track/pkg/apperror"
)

// StaffHandler handles staff HTTP endpoints.
type StaffHandler struct {
	uc *usecase.StaffUseCase
}

// NewStaffHandler creates a new StaffHandler.
func NewStaffHandler(uc *usecase.StaffUseCase) *StaffHandler {
	return &StaffHandler{uc: uc}
}

// Routes registers staff routes on the provided router.
func (h *StaffHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Post("/", h.Create)
	r.Get("/", h.List)
	r.Get("/stats", h.Stats)
	r.Get("/{id}", h.GetByID)
	r.Put("/{id}", h.Update)
	r.Delete("/{id}", h.Delete)
	return r
}

// Create handles POST /api/v1/staff
func (h *StaffHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input usecase.CreateStaffInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, &apperror.Error{Kind: apperror.KindValidation, Message: "invalid request body"})
		return
	}

	staff, err := h.uc.Create(r.Context(), input)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, apiResponse{
		Success: true,
		Data:    staff,
	})
}

// List handles GET /api/v1/staff
func (h *StaffHandler) List(w http.ResponseWriter, r *http.Request) {
	input := usecase.ListStaffInput{
		Search:      r.URL.Query().Get("search"),
		Status:      r.URL.Query().Get("status"),
		Designation: r.URL.Query().Get("designation"),
		Page:        queryInt(r, "page", 1),
		PerPage:     queryInt(r, "per_page", 20),
	}

	output, err := h.uc.List(r.Context(), input)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, apiResponse{
		Success: true,
		Data:    output.Staff,
		Meta: &apiMeta{
			Page:       output.Page,
			PerPage:    output.PerPage,
			Total:      output.Total,
			TotalPages: output.TotalPages,
		},
	})
}

// GetByID handles GET /api/v1/staff/{id}
func (h *StaffHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, &apperror.Error{Kind: apperror.KindValidation, Message: "invalid staff ID"})
		return
	}

	staff, err := h.uc.GetByID(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, apiResponse{
		Success: true,
		Data:    staff,
	})
}

// Update handles PUT /api/v1/staff/{id}
func (h *StaffHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, &apperror.Error{Kind: apperror.KindValidation, Message: "invalid staff ID"})
		return
	}

	var input usecase.UpdateStaffInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, &apperror.Error{Kind: apperror.KindValidation, Message: "invalid request body"})
		return
	}

	staff, err := h.uc.Update(r.Context(), id, input)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, apiResponse{
		Success: true,
		Data:    staff,
	})
}

// Delete handles DELETE /api/v1/staff/{id}
func (h *StaffHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, &apperror.Error{Kind: apperror.KindValidation, Message: "invalid staff ID"})
		return
	}

	if err := h.uc.Delete(r.Context(), id); err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, apiResponse{
		Success: true,
		Data:    map[string]string{"message": "staff deleted successfully"},
	})
}

// Stats handles GET /api/v1/staff/stats
func (h *StaffHandler) Stats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.uc.Stats(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, apiResponse{
		Success: true,
		Data:    stats,
	})
}

// --- response helpers ---

type apiResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *apiError   `json:"error,omitempty"`
	Meta    *apiMeta    `json:"meta,omitempty"`
}

type apiError struct {
	Type    string `json:"type"`
	Message string `json:"message"`
	Code    string `json:"code,omitempty"`
	Field   string `json:"field,omitempty"`
}

type apiMeta struct {
	Page       int `json:"page"`
	PerPage    int `json:"per_page"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, err error) {
	appErr, ok := apperror.AsError(err)
	if !ok {
		appErr = &apperror.Error{Kind: apperror.KindInternal, Message: "internal server error"}
	}

	status := appErr.HTTPStatus()
	resp := apiResponse{
		Success: false,
		Error: &apiError{
			Type:    kindToString(appErr.Kind),
			Message: appErr.Message,
			Code:    appErr.Code,
			Field:   appErr.Field,
		},
	}

	writeJSON(w, status, resp)
}

func kindToString(k apperror.Kind) string {
	switch k {
	case apperror.KindValidation:
		return "validation_error"
	case apperror.KindNotFound:
		return "not_found"
	case apperror.KindConflict:
		return "conflict"
	case apperror.KindBusiness:
		return "business_error"
	default:
		return "internal_error"
	}
}

func queryInt(r *http.Request, key string, defaultVal int) int {
	s := r.URL.Query().Get(key)
	if s == "" {
		return defaultVal
	}
	v, err := strconv.Atoi(s)
	if err != nil || v < 1 {
		return defaultVal
	}
	return v
}
