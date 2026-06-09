package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/salonflow/salonflow-track/internal/core/usecase"
	"github.com/salonflow/salonflow-track/pkg/apperror"
)

// PerformanceHandler handles performance HTTP endpoints.
type PerformanceHandler struct {
	uc *usecase.PerformanceUseCase
}

// NewPerformanceHandler creates a new PerformanceHandler.
func NewPerformanceHandler(uc *usecase.PerformanceUseCase) *PerformanceHandler {
	return &PerformanceHandler{uc: uc}
}

// Routes registers performance routes.
func (h *PerformanceHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/daily", h.Daily)
	r.Get("/weekly", h.Weekly)
	r.Get("/monthly", h.Monthly)
	r.Get("/top-performers", h.TopPerformers)
	r.Get("/revenue-trend", h.RevenueTrend)
	r.Get("/stats", h.Stats)
	r.Get("/staff/{id}", h.StaffDetail)
	r.Get("/staff/{id}/trend", h.StaffTrend)
	return r
}

// Daily handles GET /api/v1/performance/daily
func (h *PerformanceHandler) Daily(w http.ResponseWriter, r *http.Request) {
	input := usecase.DailyPerformanceInput{
		StaffID: r.URL.Query().Get("staff_id"),
		Date:    r.URL.Query().Get("date"),
	}

	data, err := h.uc.GetDailyPerformance(r.Context(), input)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, apiResponse{Success: true, Data: data})
}

// Weekly handles GET /api/v1/performance/weekly
func (h *PerformanceHandler) Weekly(w http.ResponseWriter, r *http.Request) {
	input := usecase.PeriodPerformanceInput{
		StaffID:  r.URL.Query().Get("staff_id"),
		DateFrom: r.URL.Query().Get("date_from"),
		DateTo:   r.URL.Query().Get("date_to"),
	}

	data, err := h.uc.GetWeeklyPerformance(r.Context(), input)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, apiResponse{Success: true, Data: data})
}

// Monthly handles GET /api/v1/performance/monthly
func (h *PerformanceHandler) Monthly(w http.ResponseWriter, r *http.Request) {
	input := usecase.PeriodPerformanceInput{
		StaffID:  r.URL.Query().Get("staff_id"),
		DateFrom: r.URL.Query().Get("date_from"),
		DateTo:   r.URL.Query().Get("date_to"),
	}

	data, err := h.uc.GetMonthlyPerformance(r.Context(), input)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, apiResponse{Success: true, Data: data})
}

// TopPerformers handles GET /api/v1/performance/top-performers
func (h *PerformanceHandler) TopPerformers(w http.ResponseWriter, r *http.Request) {
	input := usecase.TopPerformersInput{
		DateFrom: r.URL.Query().Get("date_from"),
		DateTo:   r.URL.Query().Get("date_to"),
		Limit:    queryInt(r, "limit", 10),
	}

	data, err := h.uc.GetTopPerformers(r.Context(), input)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, apiResponse{Success: true, Data: data})
}

// RevenueTrend handles GET /api/v1/performance/revenue-trend
func (h *PerformanceHandler) RevenueTrend(w http.ResponseWriter, r *http.Request) {
	dateFrom := r.URL.Query().Get("date_from")
	dateTo := r.URL.Query().Get("date_to")

	data, err := h.uc.GetRevenueTrend(r.Context(), dateFrom, dateTo)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, apiResponse{Success: true, Data: data})
}

// Stats handles GET /api/v1/performance/stats
func (h *PerformanceHandler) Stats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.uc.GetStats(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, apiResponse{Success: true, Data: stats})
}

// StaffDetail handles GET /api/v1/performance/staff/{id}
func (h *PerformanceHandler) StaffDetail(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, &apperror.Error{Kind: apperror.KindValidation, Message: "invalid staff ID"})
		return
	}

	dateFrom := r.URL.Query().Get("date_from")
	dateTo := r.URL.Query().Get("date_to")

	data, err := h.uc.GetStaffPerformance(r.Context(), id, dateFrom, dateTo)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, apiResponse{Success: true, Data: data})
}

// StaffTrend handles GET /api/v1/performance/staff/{id}/trend
func (h *PerformanceHandler) StaffTrend(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, &apperror.Error{Kind: apperror.KindValidation, Message: "invalid staff ID"})
		return
	}

	dateFrom := r.URL.Query().Get("date_from")
	dateTo := r.URL.Query().Get("date_to")

	data, err := h.uc.GetStaffRevenueTrend(r.Context(), id, dateFrom, dateTo)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, apiResponse{Success: true, Data: data})
}
