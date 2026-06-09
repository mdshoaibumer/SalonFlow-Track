package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/salonflow/salonflow-track/internal/core/usecase"
	"github.com/salonflow/salonflow-track/pkg/apperror"
)

// SalaryHandler handles salary and advance HTTP endpoints.
type SalaryHandler struct {
	uc *usecase.SalaryUseCase
}

// NewSalaryHandler creates a new SalaryHandler.
func NewSalaryHandler(uc *usecase.SalaryUseCase) *SalaryHandler {
	return &SalaryHandler{uc: uc}
}

// Routes registers salary and advance routes.
func (h *SalaryHandler) Routes() chi.Router {
	r := chi.NewRouter()

	// Salary
	r.Post("/generate", h.GenerateSalary)
	r.Get("/", h.ListSalaries)
	r.Get("/cycles", h.ListCycles)
	r.Get("/stats", h.Stats)
	r.Get("/{id}", h.GetSalaryByID)
	r.Post("/{id}/pay", h.PaySalary)

	// Advances
	r.Post("/advances", h.CreateAdvance)
	r.Get("/advances", h.ListAdvances)
	r.Put("/advances/{id}/approve", h.ApproveAdvance)
	r.Put("/advances/{id}/reject", h.RejectAdvance)

	return r
}

// GenerateSalary handles POST /api/v1/salary/generate
func (h *SalaryHandler) GenerateSalary(w http.ResponseWriter, r *http.Request) {
	var input usecase.GenerateSalaryInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, &apperror.Error{Kind: apperror.KindValidation, Message: "invalid request body"})
		return
	}

	output, err := h.uc.GenerateMonthlySalary(r.Context(), input)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, apiResponse{Success: true, Data: output})
}

// ListSalaries handles GET /api/v1/salary?month=&year=
func (h *SalaryHandler) ListSalaries(w http.ResponseWriter, r *http.Request) {
	input := usecase.ListSalariesInput{
		Month: queryInt(r, "month", 0),
		Year:  queryInt(r, "year", 0),
	}

	if input.Month == 0 || input.Year == 0 {
		writeError(w, &apperror.Error{Kind: apperror.KindValidation, Message: "month and year are required"})
		return
	}

	records, err := h.uc.ListSalaries(r.Context(), input)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, apiResponse{Success: true, Data: records})
}

// ListCycles handles GET /api/v1/salary/cycles?year=
func (h *SalaryHandler) ListCycles(w http.ResponseWriter, r *http.Request) {
	year := queryInt(r, "year", 0)
	cycles, err := h.uc.ListCycles(r.Context(), year)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, apiResponse{Success: true, Data: cycles})
}

// Stats handles GET /api/v1/salary/stats
func (h *SalaryHandler) Stats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.uc.GetStats(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, apiResponse{Success: true, Data: stats})
}

// GetSalaryByID handles GET /api/v1/salary/{id}
func (h *SalaryHandler) GetSalaryByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, &apperror.Error{Kind: apperror.KindValidation, Message: "invalid salary ID"})
		return
	}

	record, err := h.uc.GetSalaryByID(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, apiResponse{Success: true, Data: record})
}

// PaySalary handles POST /api/v1/salary/{id}/pay
func (h *SalaryHandler) PaySalary(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, &apperror.Error{Kind: apperror.KindValidation, Message: "invalid salary ID"})
		return
	}

	if err := h.uc.PaySalary(r.Context(), id); err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, apiResponse{Success: true, Data: map[string]string{"message": "salary paid"}})
}

// CreateAdvance handles POST /api/v1/salary/advances
func (h *SalaryHandler) CreateAdvance(w http.ResponseWriter, r *http.Request) {
	var input usecase.CreateAdvanceInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, &apperror.Error{Kind: apperror.KindValidation, Message: "invalid request body"})
		return
	}

	advance, err := h.uc.CreateAdvance(r.Context(), input)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, apiResponse{Success: true, Data: advance})
}

// ListAdvances handles GET /api/v1/salary/advances
func (h *SalaryHandler) ListAdvances(w http.ResponseWriter, r *http.Request) {
	input := usecase.ListAdvancesInput{
		StaffID: r.URL.Query().Get("staff_id"),
		Status:  r.URL.Query().Get("status"),
		Page:    queryInt(r, "page", 1),
		PerPage: queryInt(r, "per_page", 20),
	}

	output, err := h.uc.ListAdvances(r.Context(), input)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, apiResponse{
		Success: true,
		Data:    output.Advances,
		Meta: &apiMeta{
			Page:       output.Page,
			PerPage:    output.PerPage,
			Total:      output.Total,
			TotalPages: output.TotalPages,
		},
	})
}

// ApproveAdvance handles PUT /api/v1/salary/advances/{id}/approve
func (h *SalaryHandler) ApproveAdvance(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, &apperror.Error{Kind: apperror.KindValidation, Message: "invalid advance ID"})
		return
	}

	advance, err := h.uc.ApproveAdvance(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, apiResponse{Success: true, Data: advance})
}

// RejectAdvance handles PUT /api/v1/salary/advances/{id}/reject
func (h *SalaryHandler) RejectAdvance(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, &apperror.Error{Kind: apperror.KindValidation, Message: "invalid advance ID"})
		return
	}

	advance, err := h.uc.RejectAdvance(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, apiResponse{Success: true, Data: advance})
}
