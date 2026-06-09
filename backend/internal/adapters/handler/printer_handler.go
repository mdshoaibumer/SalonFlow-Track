package handler

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/salonflow/salonflow-track/internal/core/domain"
	"github.com/salonflow/salonflow-track/internal/core/usecase"
)

// PrinterHandler handles print API endpoints.
type PrinterHandler struct {
	uc *usecase.PrinterUseCase
}

// NewPrinterHandler creates a new PrinterHandler.
func NewPrinterHandler(uc *usecase.PrinterUseCase) *PrinterHandler {
	return &PrinterHandler{uc: uc}
}

// Routes returns printer routes.
func (h *PrinterHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/settings", h.GetSettings)
	r.Post("/settings", h.SaveSettings)
	r.Post("/invoice", h.PrintInvoice)
	r.Post("/receipt", h.PrintReceipt)
	r.Post("/test", h.PrintTest)
	r.Get("/history", h.History)
	r.Get("/{id}", h.GetJob)
	return r
}

// GetSettings handles GET /print/settings.
func (h *PrinterHandler) GetSettings(w http.ResponseWriter, r *http.Request) {
	settings, err := h.uc.GetSettings(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, apiResponse{Success: true, Data: settings})
}

// SaveSettings handles POST /print/settings.
func (h *PrinterHandler) SaveSettings(w http.ResponseWriter, r *http.Request) {
	var settings domain.PrinterSettings
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

// PrintInvoice handles POST /print/invoice.
func (h *PrinterHandler) PrintInvoice(w http.ResponseWriter, r *http.Request) {
	var data domain.ReceiptData
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		writeJSON(w, http.StatusBadRequest, apiResponse{Success: false, Error: &apiError{Code: "INVALID_JSON", Message: "Invalid request body"}})
		return
	}

	job, receipt, err := h.uc.PrintInvoice(r.Context(), &data)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, apiResponse{Success: true, Data: map[string]interface{}{
		"job":     job,
		"receipt": receipt,
	}})
}

// PrintReceipt handles POST /print/receipt.
func (h *PrinterHandler) PrintReceipt(w http.ResponseWriter, r *http.Request) {
	var data domain.ReceiptData
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		writeJSON(w, http.StatusBadRequest, apiResponse{Success: false, Error: &apiError{Code: "INVALID_JSON", Message: "Invalid request body"}})
		return
	}

	job, escpos, err := h.uc.PrintReceipt(r.Context(), &data)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, apiResponse{Success: true, Data: map[string]interface{}{
		"job":    job,
		"escpos": base64.StdEncoding.EncodeToString(escpos),
	}})
}

// PrintTest handles POST /print/test.
func (h *PrinterHandler) PrintTest(w http.ResponseWriter, r *http.Request) {
	job, escpos, err := h.uc.PrintTest(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, apiResponse{Success: true, Data: map[string]interface{}{
		"job":    job,
		"escpos": base64.StdEncoding.EncodeToString(escpos),
	}})
}

// History handles GET /print/history.
func (h *PrinterHandler) History(w http.ResponseWriter, r *http.Request) {
	limit := 20
	offset := 0
	if l := r.URL.Query().Get("per_page"); l != "" {
		if v, err := strconv.Atoi(l); err == nil && v > 0 {
			limit = v
		}
	}
	if p := r.URL.Query().Get("page"); p != "" {
		if v, err := strconv.Atoi(p); err == nil && v > 1 {
			offset = (v - 1) * limit
		}
	}

	jobs, total, err := h.uc.ListPrintJobs(r.Context(), limit, offset)
	if err != nil {
		writeError(w, err)
		return
	}

	page := offset/limit + 1
	totalPages := (total + limit - 1) / limit
	writeJSON(w, http.StatusOK, apiResponse{
		Success: true,
		Data:    jobs,
		Meta:    &apiMeta{Page: page, PerPage: limit, Total: total, TotalPages: totalPages},
	})
}

// GetJob handles GET /print/{id}.
func (h *PrinterHandler) GetJob(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, apiResponse{Success: false, Error: &apiError{Code: "INVALID_ID", Message: "Invalid print job ID"}})
		return
	}

	job, err := h.uc.GetPrintJob(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, apiResponse{Success: true, Data: job})
}
