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

// WhatsAppHandler handles WhatsApp API endpoints.
type WhatsAppHandler struct {
	uc *usecase.WhatsAppUseCase
}

// NewWhatsAppHandler creates a new WhatsAppHandler.
func NewWhatsAppHandler(uc *usecase.WhatsAppUseCase) *WhatsAppHandler {
	return &WhatsAppHandler{uc: uc}
}

// Routes returns WhatsApp routes.
func (h *WhatsAppHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Post("/templates", h.CreateTemplate)
	r.Get("/templates", h.ListTemplates)
	r.Put("/templates/{id}", h.UpdateTemplate)
	r.Delete("/templates/{id}", h.DeleteTemplate)
	r.Post("/send", h.SendMessage)
	r.Get("/messages", h.ListMessages)
	r.Get("/stats", h.GetStats)
	r.Post("/rules", h.CreateRule)
	r.Get("/rules", h.ListRules)
	r.Put("/rules/{id}", h.UpdateRule)
	r.Delete("/rules/{id}", h.DeleteRule)
	return r
}

// CreateTemplate handles POST /whatsapp/templates.
func (h *WhatsAppHandler) CreateTemplate(w http.ResponseWriter, r *http.Request) {
	var tmpl domain.WhatsAppTemplate
	if err := json.NewDecoder(r.Body).Decode(&tmpl); err != nil {
		writeJSON(w, http.StatusBadRequest, apiResponse{Success: false, Error: &apiError{Code: "INVALID_JSON", Message: "Invalid request body"}})
		return
	}
	newTmpl := domain.NewWhatsAppTemplate(tmpl.Name, tmpl.Category, tmpl.Body)
	newTmpl.Variables = tmpl.Variables
	if err := h.uc.CreateTemplate(r.Context(), newTmpl); err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, apiResponse{Success: true, Data: newTmpl})
}

// ListTemplates handles GET /whatsapp/templates.
func (h *WhatsAppHandler) ListTemplates(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")
	templates, err := h.uc.ListTemplates(r.Context(), category)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, apiResponse{Success: true, Data: templates})
}

// UpdateTemplate handles PUT /whatsapp/templates/{id}.
func (h *WhatsAppHandler) UpdateTemplate(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, apiResponse{Success: false, Error: &apiError{Code: "INVALID_ID", Message: "Invalid ID"}})
		return
	}
	var tmpl domain.WhatsAppTemplate
	if err := json.NewDecoder(r.Body).Decode(&tmpl); err != nil {
		writeJSON(w, http.StatusBadRequest, apiResponse{Success: false, Error: &apiError{Code: "INVALID_JSON", Message: "Invalid request body"}})
		return
	}
	tmpl.ID = id
	if err := h.uc.UpdateTemplate(r.Context(), &tmpl); err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, apiResponse{Success: true, Data: tmpl})
}

// DeleteTemplate handles DELETE /whatsapp/templates/{id}.
func (h *WhatsAppHandler) DeleteTemplate(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, apiResponse{Success: false, Error: &apiError{Code: "INVALID_ID", Message: "Invalid ID"}})
		return
	}
	if err := h.uc.DeleteTemplate(r.Context(), id); err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, apiResponse{Success: true, Data: map[string]string{"message": "Template deleted"}})
}

// SendMessage handles POST /whatsapp/send.
func (h *WhatsAppHandler) SendMessage(w http.ResponseWriter, r *http.Request) {
	var req struct {
		TemplateID string            `json:"template_id"`
		Phone      string            `json:"phone"`
		Name       string            `json:"name"`
		Variables  map[string]string `json:"variables"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, apiResponse{Success: false, Error: &apiError{Code: "INVALID_JSON", Message: "Invalid request body"}})
		return
	}
	msg, err := h.uc.SendMessage(r.Context(), req.TemplateID, req.Phone, req.Name, req.Variables)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, apiResponse{Success: true, Data: msg})
}

// ListMessages handles GET /whatsapp/messages.
func (h *WhatsAppHandler) ListMessages(w http.ResponseWriter, r *http.Request) {
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
	status := r.URL.Query().Get("status")

	messages, total, err := h.uc.ListMessages(r.Context(), limit, offset, status)
	if err != nil {
		writeError(w, err)
		return
	}
	page := offset/limit + 1
	totalPages := (total + limit - 1) / limit
	writeJSON(w, http.StatusOK, apiResponse{
		Success: true,
		Data:    messages,
		Meta:    &apiMeta{Page: page, PerPage: limit, Total: total, TotalPages: totalPages},
	})
}

// GetStats handles GET /whatsapp/stats.
func (h *WhatsAppHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.uc.GetStats(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, apiResponse{Success: true, Data: stats})
}

// CreateRule handles POST /whatsapp/rules.
func (h *WhatsAppHandler) CreateRule(w http.ResponseWriter, r *http.Request) {
	var rule domain.AutomationRule
	if err := json.NewDecoder(r.Body).Decode(&rule); err != nil {
		writeJSON(w, http.StatusBadRequest, apiResponse{Success: false, Error: &apiError{Code: "INVALID_JSON", Message: "Invalid request body"}})
		return
	}
	newRule := domain.NewAutomationRule(rule.Name, rule.TriggerType, rule.TemplateID, rule.DelayMinutes)
	if err := h.uc.CreateRule(r.Context(), newRule); err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, apiResponse{Success: true, Data: newRule})
}

// ListRules handles GET /whatsapp/rules.
func (h *WhatsAppHandler) ListRules(w http.ResponseWriter, r *http.Request) {
	rules, err := h.uc.ListRules(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, apiResponse{Success: true, Data: rules})
}

// UpdateRule handles PUT /whatsapp/rules/{id}.
func (h *WhatsAppHandler) UpdateRule(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, apiResponse{Success: false, Error: &apiError{Code: "INVALID_ID", Message: "Invalid ID"}})
		return
	}
	var rule domain.AutomationRule
	if err := json.NewDecoder(r.Body).Decode(&rule); err != nil {
		writeJSON(w, http.StatusBadRequest, apiResponse{Success: false, Error: &apiError{Code: "INVALID_JSON", Message: "Invalid request body"}})
		return
	}
	rule.ID = id
	if err := h.uc.UpdateRule(r.Context(), &rule); err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, apiResponse{Success: true, Data: rule})
}

// DeleteRule handles DELETE /whatsapp/rules/{id}.
func (h *WhatsAppHandler) DeleteRule(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, apiResponse{Success: false, Error: &apiError{Code: "INVALID_ID", Message: "Invalid ID"}})
		return
	}
	if err := h.uc.DeleteRule(r.Context(), id); err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, apiResponse{Success: true, Data: map[string]string{"message": "Rule deleted"}})
}
