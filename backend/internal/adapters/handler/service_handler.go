package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/salonflow/salonflow-track/internal/core/usecase"
	"github.com/salonflow/salonflow-track/pkg/apperror"
)

// ServiceHandler handles service HTTP endpoints.
type ServiceHandler struct {
	uc *usecase.ServiceUseCase
}

// NewServiceHandler creates a new ServiceHandler.
func NewServiceHandler(uc *usecase.ServiceUseCase) *ServiceHandler {
	return &ServiceHandler{uc: uc}
}

// Routes registers service routes on the provided router.
func (h *ServiceHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Post("/", h.Create)
	r.Get("/", h.List)
	r.Get("/stats", h.Stats)
	r.Get("/{id}", h.GetByID)
	r.Put("/{id}", h.Update)
	r.Delete("/{id}", h.Delete)
	return r
}

// Create handles POST /api/v1/services
func (h *ServiceHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input usecase.CreateServiceInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, &apperror.Error{Kind: apperror.KindValidation, Message: "invalid request body"})
		return
	}

	svc, err := h.uc.Create(r.Context(), input)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, apiResponse{
		Success: true,
		Data:    svc,
	})
}

// List handles GET /api/v1/services
func (h *ServiceHandler) List(w http.ResponseWriter, r *http.Request) {
	input := usecase.ListServiceInput{
		Search:   r.URL.Query().Get("search"),
		Status:   r.URL.Query().Get("status"),
		Category: r.URL.Query().Get("category"),
		Page:     queryInt(r, "page", 1),
		PerPage:  queryInt(r, "per_page", 20),
	}

	output, err := h.uc.List(r.Context(), input)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, apiResponse{
		Success: true,
		Data:    output.Services,
		Meta: &apiMeta{
			Page:       output.Page,
			PerPage:    output.PerPage,
			Total:      output.Total,
			TotalPages: output.TotalPages,
		},
	})
}

// GetByID handles GET /api/v1/services/{id}
func (h *ServiceHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, &apperror.Error{Kind: apperror.KindValidation, Message: "invalid service ID"})
		return
	}

	svc, err := h.uc.GetByID(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, apiResponse{
		Success: true,
		Data:    svc,
	})
}

// Update handles PUT /api/v1/services/{id}
func (h *ServiceHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, &apperror.Error{Kind: apperror.KindValidation, Message: "invalid service ID"})
		return
	}

	var input usecase.UpdateServiceInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, &apperror.Error{Kind: apperror.KindValidation, Message: "invalid request body"})
		return
	}

	svc, err := h.uc.Update(r.Context(), id, input)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, apiResponse{
		Success: true,
		Data:    svc,
	})
}

// Delete handles DELETE /api/v1/services/{id}
func (h *ServiceHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, &apperror.Error{Kind: apperror.KindValidation, Message: "invalid service ID"})
		return
	}

	if err := h.uc.Delete(r.Context(), id); err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, apiResponse{
		Success: true,
		Data:    map[string]string{"message": "service deleted successfully"},
	})
}

// Stats handles GET /api/v1/services/stats
func (h *ServiceHandler) Stats(w http.ResponseWriter, r *http.Request) {
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
