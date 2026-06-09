package handler

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/salonflow/salonflow-track/internal/core/usecase"
)

// AnalyticsHandler handles analytics and reports API.
type AnalyticsHandler struct {
	uc *usecase.AnalyticsUseCase
}

// NewAnalyticsHandler creates a new AnalyticsHandler.
func NewAnalyticsHandler(uc *usecase.AnalyticsUseCase) *AnalyticsHandler {
	return &AnalyticsHandler{uc: uc}
}

// Routes returns analytics routes.
func (h *AnalyticsHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/dashboard", h.Dashboard)
	r.Get("/kpis", h.KPIs)
	r.Get("/revenue", h.Revenue)
	r.Get("/customers", h.Customers)
	r.Get("/staff", h.Staff)
	r.Get("/services", h.Services)
	r.Get("/expenses", h.Expenses)
	r.Get("/inventory", h.Inventory)
	r.Get("/profit-loss", h.ProfitLoss)
	return r
}

func (h *AnalyticsHandler) dateRange(r *http.Request) (string, string) {
	from := r.URL.Query().Get("date_from")
	to := r.URL.Query().Get("date_to")
	if from == "" {
		from = time.Now().Format("2006-01") + "-01"
	}
	if to == "" {
		to = time.Now().Format("2006-01-02")
	}
	return from, to
}

// Dashboard returns executive dashboard stats.
func (h *AnalyticsHandler) Dashboard(w http.ResponseWriter, r *http.Request) {
	stats, err := h.uc.GetDashboard(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, apiResponse{Success: true, Data: stats})
}

// KPIs returns KPI metrics.
func (h *AnalyticsHandler) KPIs(w http.ResponseWriter, r *http.Request) {
	from, to := h.dateRange(r)
	kpis, err := h.uc.GetKPIs(r.Context(), from, to)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, apiResponse{Success: true, Data: kpis})
}

// Revenue returns revenue analytics report.
func (h *AnalyticsHandler) Revenue(w http.ResponseWriter, r *http.Request) {
	from, to := h.dateRange(r)
	groupBy := r.URL.Query().Get("group_by")
	if groupBy == "" {
		groupBy = "day"
	}
	report, err := h.uc.GetRevenueReport(r.Context(), from, to, groupBy)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, apiResponse{Success: true, Data: report})
}

// Customers returns customer analytics report.
func (h *AnalyticsHandler) Customers(w http.ResponseWriter, r *http.Request) {
	from, to := h.dateRange(r)
	report, err := h.uc.GetCustomerReport(r.Context(), from, to)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, apiResponse{Success: true, Data: report})
}

// Staff returns staff analytics report.
func (h *AnalyticsHandler) Staff(w http.ResponseWriter, r *http.Request) {
	from, to := h.dateRange(r)
	report, err := h.uc.GetStaffReport(r.Context(), from, to)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, apiResponse{Success: true, Data: report})
}

// Services returns service analytics report.
func (h *AnalyticsHandler) Services(w http.ResponseWriter, r *http.Request) {
	from, to := h.dateRange(r)
	report, err := h.uc.GetServiceReport(r.Context(), from, to)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, apiResponse{Success: true, Data: report})
}

// Expenses returns expense analytics report.
func (h *AnalyticsHandler) Expenses(w http.ResponseWriter, r *http.Request) {
	from, to := h.dateRange(r)
	report, err := h.uc.GetExpenseReport(r.Context(), from, to)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, apiResponse{Success: true, Data: report})
}

// Inventory returns inventory analytics report.
func (h *AnalyticsHandler) Inventory(w http.ResponseWriter, r *http.Request) {
	report, err := h.uc.GetInventoryReport(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, apiResponse{Success: true, Data: report})
}

// ProfitLoss returns P&L report.
func (h *AnalyticsHandler) ProfitLoss(w http.ResponseWriter, r *http.Request) {
	from, to := h.dateRange(r)
	groupBy := r.URL.Query().Get("group_by")
	if groupBy == "" {
		groupBy = "month"
	}
	report, err := h.uc.GetProfitLossReport(r.Context(), from, to, groupBy)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, apiResponse{Success: true, Data: report})
}
