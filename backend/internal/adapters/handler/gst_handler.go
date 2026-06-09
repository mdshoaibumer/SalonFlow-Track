package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/salonflow/salonflow-track/internal/core/domain"
	"github.com/salonflow/salonflow-track/internal/core/usecase"
)

// GSTHandler handles GST API endpoints.
type GSTHandler struct {
	uc *usecase.GSTUseCase
}

// NewGSTHandler creates a new GSTHandler.
func NewGSTHandler(uc *usecase.GSTUseCase) *GSTHandler {
	return &GSTHandler{uc: uc}
}

// Routes returns GST routes.
func (h *GSTHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/settings", h.GetSettings)
	r.Post("/settings", h.SaveSettings)
	r.Get("/tax-rates", h.ListTaxRates)
	r.Post("/tax-rates", h.CreateTaxRate)
	r.Put("/tax-rates/{id}", h.UpdateTaxRate)
	r.Delete("/tax-rates/{id}", h.DeleteTaxRate)
	r.Get("/invoice/{id}/tax", h.GetInvoiceTax)
	r.Get("/reports", h.GetReport)
	return r
}

// GetSettings handles GET /gst/settings.
func (h *GSTHandler) GetSettings(w http.ResponseWriter, r *http.Request) {
	settings, err := h.uc.GetSettings(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, apiResponse{Success: true, Data: settings})
}

// SaveSettings handles POST /gst/settings.
func (h *GSTHandler) SaveSettings(w http.ResponseWriter, r *http.Request) {
	var settings domain.GSTSettings
	if err := json.NewDecoder(r.Body).Decode(&settings); err != nil {
		writeJSON(w, http.StatusBadRequest, apiResponse{Success: false, Error: &apiError{Code: "INVALID_JSON", Message: "Invalid request body"}})
		return
	}

	if err := h.uc.SaveSettings(r.Context(), &settings); err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, apiResponse{Success: true, Data: settings})
}

// ListTaxRates handles GET /gst/tax-rates.
func (h *GSTHandler) ListTaxRates(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")
	rates, err := h.uc.ListTaxRates(r.Context(), category)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, apiResponse{Success: true, Data: rates})
}

// CreateTaxRate handles POST /gst/tax-rates.
func (h *GSTHandler) CreateTaxRate(w http.ResponseWriter, r *http.Request) {
	var rate domain.TaxRate
	if err := json.NewDecoder(r.Body).Decode(&rate); err != nil {
		writeJSON(w, http.StatusBadRequest, apiResponse{Success: false, Error: &apiError{Code: "INVALID_JSON", Message: "Invalid request body"}})
		return
	}

	newRate := domain.NewTaxRate(rate.Name, rate.HSNCode, rate.Category, rate.CGSTRate, rate.SGSTRate, rate.IGSTRate)
	if err := h.uc.CreateTaxRate(r.Context(), newRate); err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, apiResponse{Success: true, Data: newRate})
}

// UpdateTaxRate handles PUT /gst/tax-rates/{id}.
func (h *GSTHandler) UpdateTaxRate(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, apiResponse{Success: false, Error: &apiError{Code: "INVALID_ID", Message: "Invalid tax rate ID"}})
		return
	}

	var rate domain.TaxRate
	if err := json.NewDecoder(r.Body).Decode(&rate); err != nil {
		writeJSON(w, http.StatusBadRequest, apiResponse{Success: false, Error: &apiError{Code: "INVALID_JSON", Message: "Invalid request body"}})
		return
	}
	rate.ID = id

	if err := h.uc.UpdateTaxRate(r.Context(), &rate); err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, apiResponse{Success: true, Data: rate})
}

// DeleteTaxRate handles DELETE /gst/tax-rates/{id}.
func (h *GSTHandler) DeleteTaxRate(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, apiResponse{Success: false, Error: &apiError{Code: "INVALID_ID", Message: "Invalid tax rate ID"}})
		return
	}

	if err := h.uc.DeleteTaxRate(r.Context(), id); err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, apiResponse{Success: true, Data: map[string]string{"message": "Tax rate deleted"}})
}

// GetInvoiceTax handles GET /gst/invoice/{id}/tax.
func (h *GSTHandler) GetInvoiceTax(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, apiResponse{Success: false, Error: &apiError{Code: "INVALID_ID", Message: "Invalid invoice ID"}})
		return
	}

	lines, err := h.uc.GetInvoiceTaxLines(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, apiResponse{Success: true, Data: lines})
}

// GetReport handles GET /gst/reports.
func (h *GSTHandler) GetReport(w http.ResponseWriter, r *http.Request) {
	filter := domain.GSTReportFilter{
		Period:    r.URL.Query().Get("period"),
		StartDate: r.URL.Query().Get("start_date"),
		EndDate:   r.URL.Query().Get("end_date"),
	}

	report, err := h.uc.GetReport(r.Context(), filter)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, apiResponse{Success: true, Data: report})
}
