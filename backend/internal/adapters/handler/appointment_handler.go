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

// AppointmentHandler handles appointment API endpoints.
type AppointmentHandler struct {
	uc *usecase.AppointmentUseCase
}

// NewAppointmentHandler creates a new AppointmentHandler.
func NewAppointmentHandler(uc *usecase.AppointmentUseCase) *AppointmentHandler {
	return &AppointmentHandler{uc: uc}
}

// Routes returns appointment routes.
func (h *AppointmentHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Post("/", h.Create)
	r.Get("/", h.List)
	r.Get("/{id}", h.Get)
	r.Put("/{id}", h.Update)
	r.Put("/{id}/status", h.UpdateStatus)
	r.Delete("/{id}", h.Delete)
	r.Get("/{id}/history", h.History)
	return r
}

type createAppointmentRequest struct {
	CustomerID      string                      `json:"customer_id"`
	StaffID         string                      `json:"staff_id"`
	AppointmentDate string                      `json:"appointment_date"`
	StartTime       string                      `json:"start_time"`
	EndTime         string                      `json:"end_time"`
	Notes           string                      `json:"notes"`
	IsWalkin        bool                        `json:"is_walkin"`
	Services        []appointmentServiceRequest `json:"services"`
}

type appointmentServiceRequest struct {
	ServiceID       string  `json:"service_id"`
	ServiceName     string  `json:"service_name"`
	DurationMinutes int     `json:"duration_minutes"`
	Price           float64 `json:"price"`
}

// Create handles POST /appointments.
func (h *AppointmentHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createAppointmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, apiResponse{Success: false, Error: &apiError{Code: "INVALID_JSON", Message: "Invalid request body"}})
		return
	}

	appt := domain.NewAppointment(req.CustomerID, req.StaffID, req.AppointmentDate, req.StartTime, req.EndTime, req.IsWalkin)
	appt.Notes = req.Notes

	var services []domain.AppointmentService
	for _, s := range req.Services {
		services = append(services, *domain.NewAppointmentService(appt.ID, s.ServiceID, s.ServiceName, s.DurationMinutes, s.Price))
	}

	if err := h.uc.Create(r.Context(), appt, services); err != nil {
		writeError(w, err)
		return
	}

	appt.Services = services
	writeJSON(w, http.StatusCreated, apiResponse{Success: true, Data: appt})
}

// List handles GET /appointments.
func (h *AppointmentHandler) List(w http.ResponseWriter, r *http.Request) {
	limit := 50
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

	filter := domain.AppointmentFilter{
		Date:       r.URL.Query().Get("date"),
		StaffID:    r.URL.Query().Get("staff_id"),
		CustomerID: r.URL.Query().Get("customer_id"),
		Status:     r.URL.Query().Get("status"),
		StartDate:  r.URL.Query().Get("start_date"),
		EndDate:    r.URL.Query().Get("end_date"),
		Limit:      limit,
		Offset:     offset,
	}

	appts, total, err := h.uc.List(r.Context(), filter)
	if err != nil {
		writeError(w, err)
		return
	}

	page := offset/limit + 1
	totalPages := (total + limit - 1) / limit
	writeJSON(w, http.StatusOK, apiResponse{
		Success: true,
		Data:    appts,
		Meta:    &apiMeta{Page: page, PerPage: limit, Total: total, TotalPages: totalPages},
	})
}

// Get handles GET /appointments/{id}.
func (h *AppointmentHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, apiResponse{Success: false, Error: &apiError{Code: "INVALID_ID", Message: "Invalid appointment ID"}})
		return
	}

	appt, err := h.uc.GetByID(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, apiResponse{Success: true, Data: appt})
}

// Update handles PUT /appointments/{id}.
func (h *AppointmentHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, apiResponse{Success: false, Error: &apiError{Code: "INVALID_ID", Message: "Invalid appointment ID"}})
		return
	}

	var req createAppointmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, apiResponse{Success: false, Error: &apiError{Code: "INVALID_JSON", Message: "Invalid request body"}})
		return
	}

	appt := &domain.Appointment{
		ID:              id,
		CustomerID:      req.CustomerID,
		StaffID:         req.StaffID,
		AppointmentDate: req.AppointmentDate,
		StartTime:       req.StartTime,
		EndTime:         req.EndTime,
		Notes:           req.Notes,
		IsWalkin:        req.IsWalkin,
		Status:          "booked",
	}

	// Check if status is sent
	if s := r.URL.Query().Get("status"); s != "" {
		appt.Status = s
	}

	var services []domain.AppointmentService
	for _, s := range req.Services {
		services = append(services, *domain.NewAppointmentService(id, s.ServiceID, s.ServiceName, s.DurationMinutes, s.Price))
	}

	if err := h.uc.Update(r.Context(), appt, services); err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, apiResponse{Success: true, Data: appt})
}

// UpdateStatus handles PUT /appointments/{id}/status.
func (h *AppointmentHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, apiResponse{Success: false, Error: &apiError{Code: "INVALID_ID", Message: "Invalid appointment ID"}})
		return
	}

	var req struct {
		Status string `json:"status"`
		Note   string `json:"note"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, apiResponse{Success: false, Error: &apiError{Code: "INVALID_JSON", Message: "Invalid request body"}})
		return
	}

	if err := h.uc.UpdateStatus(r.Context(), id, req.Status, req.Note); err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, apiResponse{Success: true, Data: map[string]string{"status": req.Status}})
}

// Delete handles DELETE /appointments/{id}.
func (h *AppointmentHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, apiResponse{Success: false, Error: &apiError{Code: "INVALID_ID", Message: "Invalid appointment ID"}})
		return
	}

	if err := h.uc.Delete(r.Context(), id); err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, apiResponse{Success: true, Data: map[string]string{"message": "Appointment deleted"}})
}

// History handles GET /appointments/{id}/history.
func (h *AppointmentHandler) History(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, apiResponse{Success: false, Error: &apiError{Code: "INVALID_ID", Message: "Invalid appointment ID"}})
		return
	}

	history, err := h.uc.GetHistory(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, apiResponse{Success: true, Data: history})
}
