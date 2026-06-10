package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/salonflow/salonflow-track/internal/core/usecase"
)

// LicenseHandler handles license API endpoints.
type LicenseHandler struct {
	uc *usecase.LicenseUseCase
}

// NewLicenseHandler creates a new LicenseHandler.
func NewLicenseHandler(uc *usecase.LicenseUseCase) *LicenseHandler {
	return &LicenseHandler{uc: uc}
}

// Routes returns license routes.
func (h *LicenseHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", h.GetStatus)
	r.Post("/validate", h.Validate)
	r.Post("/activate", h.Activate)
	r.Post("/renew", h.Renew)
	r.Get("/events", h.ListEvents)
	return r
}

// GetStatus handles GET /license - returns current license status.
func (h *LicenseHandler) GetStatus(w http.ResponseWriter, r *http.Request) {
	status, err := h.uc.GetStatus(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, apiResponse{Success: true, Data: status})
}

// Validate handles POST /license/validate - validates the license.
func (h *LicenseHandler) Validate(w http.ResponseWriter, r *http.Request) {
	result, err := h.uc.Validate(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, apiResponse{Success: true, Data: result})
}

// Activate handles POST /license/activate - activates a new license.
func (h *LicenseHandler) Activate(w http.ResponseWriter, r *http.Request) {
	var req struct {
		LicenseKey   string `json:"license_key"`
		CustomerName string `json:"customer_name"`
		SalonName    string `json:"salon_name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, apiResponse{Success: false, Error: &apiError{Code: "INVALID_JSON", Message: "Invalid request body"}})
		return
	}

	lic, err := h.uc.Activate(r.Context(), req.LicenseKey, req.CustomerName, req.SalonName)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, apiResponse{Success: true, Data: lic})
}

// Renew handles POST /license/renew - renews the license.
func (h *LicenseHandler) Renew(w http.ResponseWriter, r *http.Request) {
	var req struct {
		LicenseKey string `json:"license_key"`
	}
	_ = json.NewDecoder(r.Body).Decode(&req)

	lic, err := h.uc.Renew(r.Context(), req.LicenseKey)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, apiResponse{Success: true, Data: lic})
}

// ListEvents handles GET /license/events - lists license audit events.
func (h *LicenseHandler) ListEvents(w http.ResponseWriter, r *http.Request) {
	page := queryInt(r, "page", 1)
	perPage := queryInt(r, "per_page", 20)

	events, total, err := h.uc.ListEvents(r.Context(), page, perPage)
	if err != nil {
		writeError(w, err)
		return
	}
	totalPages := (total + perPage - 1) / perPage
	writeJSON(w, http.StatusOK, apiResponse{
		Success: true,
		Data:    events,
		Meta:    &apiMeta{Page: page, PerPage: perPage, Total: total, TotalPages: totalPages},
	})
}
