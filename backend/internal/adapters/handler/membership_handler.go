package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/salonflow/salonflow-track/internal/core/domain"
	"github.com/salonflow/salonflow-track/internal/core/usecase"
)

// MembershipHandler handles membership API endpoints.
type MembershipHandler struct {
	uc *usecase.MembershipUseCase
}

// NewMembershipHandler creates a new MembershipHandler.
func NewMembershipHandler(uc *usecase.MembershipUseCase) *MembershipHandler {
	return &MembershipHandler{uc: uc}
}

// Routes returns membership routes.
func (h *MembershipHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Post("/plans", h.CreatePlan)
	r.Get("/plans", h.ListPlans)
	r.Get("/plans/{id}", h.GetPlan)
	r.Put("/plans/{id}", h.UpdatePlan)
	r.Delete("/plans/{id}", h.DeletePlan)
	r.Post("/sell", h.SellPlan)
	r.Post("/use-session", h.UseSession)
	r.Get("/subscriptions", h.ListSubscriptions)
	r.Get("/stats", h.GetStats)
	return r
}

type createPlanRequest struct {
	Name               string                  `json:"name"`
	Description        string                  `json:"description"`
	PlanType           string                  `json:"plan_type"`
	Price              float64                 `json:"price"`
	DurationDays       int                     `json:"duration_days"`
	MaxSessions        int                     `json:"max_sessions"`
	DiscountPercentage float64                 `json:"discount_percentage"`
	PriorityBooking    bool                    `json:"priority_booking"`
	Services           []packageServiceRequest `json:"services"`
}

type packageServiceRequest struct {
	ServiceID        string `json:"service_id"`
	ServiceName      string `json:"service_name"`
	SessionsIncluded int    `json:"sessions_included"`
}

// CreatePlan handles POST /memberships/plans.
func (h *MembershipHandler) CreatePlan(w http.ResponseWriter, r *http.Request) {
	var req createPlanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, apiResponse{Success: false, Error: &apiError{Code: "INVALID_JSON", Message: "Invalid request body"}})
		return
	}

	plan := domain.NewMembershipPlan(req.Name, req.PlanType, req.Price, req.DurationDays, req.MaxSessions)
	plan.Description = req.Description
	plan.DiscountPercentage = req.DiscountPercentage
	plan.PriorityBooking = req.PriorityBooking

	var services []domain.PackageService
	for _, s := range req.Services {
		services = append(services, *domain.NewPackageService(plan.ID, s.ServiceID, s.ServiceName, s.SessionsIncluded))
	}

	if err := h.uc.CreatePlan(r.Context(), plan, services); err != nil {
		writeError(w, err)
		return
	}
	plan.Services = services
	writeJSON(w, http.StatusCreated, apiResponse{Success: true, Data: plan})
}

// ListPlans handles GET /memberships/plans.
func (h *MembershipHandler) ListPlans(w http.ResponseWriter, r *http.Request) {
	planType := r.URL.Query().Get("type")
	plans, err := h.uc.ListPlans(r.Context(), planType)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, apiResponse{Success: true, Data: plans})
}

// GetPlan handles GET /memberships/plans/{id}.
func (h *MembershipHandler) GetPlan(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, apiResponse{Success: false, Error: &apiError{Code: "INVALID_ID", Message: "Invalid ID"}})
		return
	}
	plan, err := h.uc.GetPlan(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, apiResponse{Success: true, Data: plan})
}

// UpdatePlan handles PUT /memberships/plans/{id}.
func (h *MembershipHandler) UpdatePlan(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, apiResponse{Success: false, Error: &apiError{Code: "INVALID_ID", Message: "Invalid ID"}})
		return
	}
	var req createPlanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, apiResponse{Success: false, Error: &apiError{Code: "INVALID_JSON", Message: "Invalid request body"}})
		return
	}

	plan := &domain.MembershipPlan{
		ID:                 id,
		Name:               req.Name,
		Description:        req.Description,
		PlanType:           req.PlanType,
		Price:              req.Price,
		DurationDays:       req.DurationDays,
		MaxSessions:        req.MaxSessions,
		DiscountPercentage: req.DiscountPercentage,
		PriorityBooking:    req.PriorityBooking,
		IsActive:           true,
	}

	var services []domain.PackageService
	for _, s := range req.Services {
		services = append(services, *domain.NewPackageService(id, s.ServiceID, s.ServiceName, s.SessionsIncluded))
	}

	if err := h.uc.UpdatePlan(r.Context(), plan, services); err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, apiResponse{Success: true, Data: plan})
}

// DeletePlan handles DELETE /memberships/plans/{id}.
func (h *MembershipHandler) DeletePlan(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, apiResponse{Success: false, Error: &apiError{Code: "INVALID_ID", Message: "Invalid ID"}})
		return
	}
	if err := h.uc.DeletePlan(r.Context(), id); err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, apiResponse{Success: true, Data: map[string]string{"message": "Plan deleted"}})
}

// SellPlan handles POST /memberships/sell.
func (h *MembershipHandler) SellPlan(w http.ResponseWriter, r *http.Request) {
	var req struct {
		CustomerID string  `json:"customer_id"`
		PlanID     string  `json:"plan_id"`
		AmountPaid float64 `json:"amount_paid"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, apiResponse{Success: false, Error: &apiError{Code: "INVALID_JSON", Message: "Invalid request body"}})
		return
	}
	planID, err := uuid.Parse(req.PlanID)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, apiResponse{Success: false, Error: &apiError{Code: "INVALID_ID", Message: "Invalid plan ID"}})
		return
	}
	sub, err := h.uc.SellPlan(r.Context(), req.CustomerID, planID, req.AmountPaid)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, apiResponse{Success: true, Data: sub})
}

// UseSession handles POST /memberships/use-session.
func (h *MembershipHandler) UseSession(w http.ResponseWriter, r *http.Request) {
	var req struct {
		SubscriptionID string `json:"subscription_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, apiResponse{Success: false, Error: &apiError{Code: "INVALID_JSON", Message: "Invalid request body"}})
		return
	}
	id, err := uuid.Parse(req.SubscriptionID)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, apiResponse{Success: false, Error: &apiError{Code: "INVALID_ID", Message: "Invalid subscription ID"}})
		return
	}
	if err := h.uc.UseSession(r.Context(), id); err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, apiResponse{Success: true, Data: map[string]string{"message": "Session used"}})
}

// ListSubscriptions handles GET /memberships/subscriptions.
func (h *MembershipHandler) ListSubscriptions(w http.ResponseWriter, r *http.Request) {
	limit := 20
	offset := 0
	if l := r.URL.Query().Get("per_page"); l != "" {
		if v, _ := strconv.Atoi(l); v > 0 {
			limit = v
		}
	}
	if p := r.URL.Query().Get("page"); p != "" {
		if v, _ := strconv.Atoi(p); v > 1 {
			offset = (v - 1) * limit
		}
	}
	customerID := r.URL.Query().Get("customer_id")
	status := r.URL.Query().Get("status")

	subs, total, err := h.uc.ListSubscriptions(r.Context(), customerID, status, limit, offset)
	if err != nil {
		writeError(w, err)
		return
	}
	page := offset/limit + 1
	totalPages := (total + limit - 1) / limit
	writeJSON(w, http.StatusOK, apiResponse{
		Success: true,
		Data:    subs,
		Meta:    &apiMeta{Page: page, PerPage: limit, Total: total, TotalPages: totalPages},
	})
}

// GetStats handles GET /memberships/stats.
func (h *MembershipHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.uc.GetStats(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, apiResponse{Success: true, Data: stats})
}
