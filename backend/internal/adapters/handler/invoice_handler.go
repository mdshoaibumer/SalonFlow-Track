package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/salonflow/salonflow-track/internal/core/usecase"
	"github.com/salonflow/salonflow-track/pkg/apperror"
)

// InvoiceHandler handles invoice HTTP endpoints.
type InvoiceHandler struct {
	uc *usecase.InvoiceUseCase
}

// NewInvoiceHandler creates a new InvoiceHandler.
func NewInvoiceHandler(uc *usecase.InvoiceUseCase) *InvoiceHandler {
	return &InvoiceHandler{uc: uc}
}

// Routes registers invoice routes on the provided router.
func (h *InvoiceHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Post("/", h.Create)
	r.Get("/", h.List)
	r.Get("/stats", h.Stats)
	r.Get("/{id}", h.GetByID)
	r.Post("/{id}/payment", h.RecordPayment)
	return r
}

// Create handles POST /api/v1/invoices
func (h *InvoiceHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input usecase.CreateInvoiceInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, &apperror.Error{Kind: apperror.KindValidation, Message: "invalid request body"})
		return
	}

	invoice, err := h.uc.Create(r.Context(), input)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, apiResponse{
		Success: true,
		Data:    invoice,
	})
}

// List handles GET /api/v1/invoices
func (h *InvoiceHandler) List(w http.ResponseWriter, r *http.Request) {
	input := usecase.ListInvoiceInput{
		CustomerID:    r.URL.Query().Get("customer_id"),
		StaffID:       r.URL.Query().Get("staff_id"),
		PaymentStatus: r.URL.Query().Get("payment_status"),
		DateFrom:      r.URL.Query().Get("date_from"),
		DateTo:        r.URL.Query().Get("date_to"),
		Search:        r.URL.Query().Get("search"),
		Page:          queryInt(r, "page", 1),
		PerPage:       queryInt(r, "per_page", 20),
	}

	output, err := h.uc.List(r.Context(), input)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, apiResponse{
		Success: true,
		Data:    output.Invoices,
		Meta: &apiMeta{
			Page:       output.Page,
			PerPage:    output.PerPage,
			Total:      output.Total,
			TotalPages: output.TotalPages,
		},
	})
}

// GetByID handles GET /api/v1/invoices/{id}
func (h *InvoiceHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, &apperror.Error{Kind: apperror.KindValidation, Message: "invalid invoice ID"})
		return
	}

	invoice, err := h.uc.GetByID(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, apiResponse{
		Success: true,
		Data:    invoice,
	})
}

// RecordPayment handles POST /api/v1/invoices/{id}/payment
func (h *InvoiceHandler) RecordPayment(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, &apperror.Error{Kind: apperror.KindValidation, Message: "invalid invoice ID"})
		return
	}

	var input usecase.RecordPaymentInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, &apperror.Error{Kind: apperror.KindValidation, Message: "invalid request body"})
		return
	}

	payment, err := h.uc.RecordPayment(r.Context(), id, input)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, apiResponse{
		Success: true,
		Data:    payment,
	})
}

// Stats handles GET /api/v1/invoices/stats
func (h *InvoiceHandler) Stats(w http.ResponseWriter, r *http.Request) {
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
