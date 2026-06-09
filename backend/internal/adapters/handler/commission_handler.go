package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/salonflow/salonflow-track/internal/core/usecase"
	"github.com/salonflow/salonflow-track/pkg/apperror"
)

// CommissionHandler handles commission HTTP endpoints.
type CommissionHandler struct {
	uc *usecase.CommissionUseCase
}

// NewCommissionHandler creates a new CommissionHandler.
func NewCommissionHandler(uc *usecase.CommissionUseCase) *CommissionHandler {
	return &CommissionHandler{uc: uc}
}

// Routes registers commission routes.
func (h *CommissionHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Post("/rules", h.CreateRule)
	r.Get("/rules", h.ListRules)
	r.Get("/rules/{id}", h.GetRule)
	r.Put("/rules/{id}", h.UpdateRule)
	r.Delete("/rules/{id}", h.DeleteRule)
	r.Get("/staff/{id}", h.StaffCommission)
	r.Get("/monthly", h.Monthly)
	r.Get("/stats", h.Stats)
	return r
}

// CreateRule handles POST /api/v1/commissions/rules
func (h *CommissionHandler) CreateRule(w http.ResponseWriter, r *http.Request) {
	var input usecase.CreateRuleInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, &apperror.Error{Kind: apperror.KindValidation, Message: "invalid request body"})
		return
	}

	rule, err := h.uc.CreateRule(r.Context(), input)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, apiResponse{Success: true, Data: rule})
}

// ListRules handles GET /api/v1/commissions/rules
func (h *CommissionHandler) ListRules(w http.ResponseWriter, r *http.Request) {
	input := usecase.ListRulesInput{
		RuleType:   r.URL.Query().Get("rule_type"),
		TargetType: r.URL.Query().Get("target_type"),
		Page:       queryInt(r, "page", 1),
		PerPage:    queryInt(r, "per_page", 20),
	}

	if active := r.URL.Query().Get("is_active"); active != "" {
		val := active == "true" || active == "1"
		input.IsActive = &val
	}

	output, err := h.uc.ListRules(r.Context(), input)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, apiResponse{
		Success: true,
		Data:    output.Rules,
		Meta: &apiMeta{
			Page:       output.Page,
			PerPage:    output.PerPage,
			Total:      output.Total,
			TotalPages: output.TotalPages,
		},
	})
}

// GetRule handles GET /api/v1/commissions/rules/{id}
func (h *CommissionHandler) GetRule(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, &apperror.Error{Kind: apperror.KindValidation, Message: "invalid rule ID"})
		return
	}

	rule, err := h.uc.GetRuleByID(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, apiResponse{Success: true, Data: rule})
}

// UpdateRule handles PUT /api/v1/commissions/rules/{id}
func (h *CommissionHandler) UpdateRule(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, &apperror.Error{Kind: apperror.KindValidation, Message: "invalid rule ID"})
		return
	}

	var input usecase.UpdateRuleInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, &apperror.Error{Kind: apperror.KindValidation, Message: "invalid request body"})
		return
	}

	rule, err := h.uc.UpdateRule(r.Context(), id, input)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, apiResponse{Success: true, Data: rule})
}

// DeleteRule handles DELETE /api/v1/commissions/rules/{id}
func (h *CommissionHandler) DeleteRule(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, &apperror.Error{Kind: apperror.KindValidation, Message: "invalid rule ID"})
		return
	}

	if err := h.uc.DeleteRule(r.Context(), id); err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, apiResponse{Success: true, Data: map[string]string{"message": "rule deleted"}})
}

// StaffCommission handles GET /api/v1/commissions/staff/{id}
func (h *CommissionHandler) StaffCommission(w http.ResponseWriter, r *http.Request) {
	input := usecase.GetStaffCommissionInput{
		StaffID:  chi.URLParam(r, "id"),
		DateFrom: r.URL.Query().Get("date_from"),
		DateTo:   r.URL.Query().Get("date_to"),
	}

	data, err := h.uc.GetStaffCommission(r.Context(), input)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, apiResponse{Success: true, Data: data})
}

// Monthly handles GET /api/v1/commissions/monthly
func (h *CommissionHandler) Monthly(w http.ResponseWriter, r *http.Request) {
	input := usecase.MonthlyCommissionInput{
		Month: r.URL.Query().Get("month"),
	}

	data, err := h.uc.GetMonthlyCommission(r.Context(), input)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, apiResponse{Success: true, Data: data})
}

// Stats handles GET /api/v1/commissions/stats
func (h *CommissionHandler) Stats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.uc.GetStats(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, apiResponse{Success: true, Data: stats})
}
