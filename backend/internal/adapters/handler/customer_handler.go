package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/salonflow/salonflow-track/internal/core/usecase"
	"github.com/salonflow/salonflow-track/pkg/apperror"
)

// CustomerHandler handles customer HTTP endpoints.
type CustomerHandler struct {
	uc *usecase.CustomerUseCase
}

// NewCustomerHandler creates a new CustomerHandler.
func NewCustomerHandler(uc *usecase.CustomerUseCase) *CustomerHandler {
	return &CustomerHandler{uc: uc}
}

// Routes registers customer routes on the provided router.
func (h *CustomerHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Post("/", h.Create)
	r.Get("/", h.List)
	r.Get("/stats", h.Stats)
	r.Get("/{id}", h.GetByID)
	r.Put("/{id}", h.Update)
	r.Delete("/{id}", h.Delete)
	return r
}

// Create handles POST /api/v1/customers
func (h *CustomerHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input usecase.CreateCustomerInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, &apperror.Error{Kind: apperror.KindValidation, Message: "invalid request body"})
		return
	}

	customer, err := h.uc.Create(r.Context(), input)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, apiResponse{
		Success: true,
		Data:    customer,
	})
}

// List handles GET /api/v1/customers
func (h *CustomerHandler) List(w http.ResponseWriter, r *http.Request) {
	input := usecase.ListCustomerInput{
		Search:  r.URL.Query().Get("search"),
		Status:  r.URL.Query().Get("status"),
		Page:    queryInt(r, "page", 1),
		PerPage: queryInt(r, "per_page", 20),
	}

	output, err := h.uc.List(r.Context(), input)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, apiResponse{
		Success: true,
		Data:    output.Customers,
		Meta: &apiMeta{
			Page:       output.Page,
			PerPage:    output.PerPage,
			Total:      output.Total,
			TotalPages: output.TotalPages,
		},
	})
}

// GetByID handles GET /api/v1/customers/{id}
func (h *CustomerHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, &apperror.Error{Kind: apperror.KindValidation, Message: "invalid customer ID"})
		return
	}

	customer, err := h.uc.GetByID(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, apiResponse{
		Success: true,
		Data:    customer,
	})
}

// Update handles PUT /api/v1/customers/{id}
func (h *CustomerHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, &apperror.Error{Kind: apperror.KindValidation, Message: "invalid customer ID"})
		return
	}

	var input usecase.UpdateCustomerInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, &apperror.Error{Kind: apperror.KindValidation, Message: "invalid request body"})
		return
	}

	customer, err := h.uc.Update(r.Context(), id, input)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, apiResponse{
		Success: true,
		Data:    customer,
	})
}

// Delete handles DELETE /api/v1/customers/{id}
func (h *CustomerHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, &apperror.Error{Kind: apperror.KindValidation, Message: "invalid customer ID"})
		return
	}

	if err := h.uc.Delete(r.Context(), id); err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, apiResponse{
		Success: true,
		Data:    map[string]string{"message": "customer deleted successfully"},
	})
}

// Stats handles GET /api/v1/customers/stats
func (h *CustomerHandler) Stats(w http.ResponseWriter, r *http.Request) {
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
