package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/salonflow/salonflow-track/internal/core/usecase"
	"github.com/salonflow/salonflow-track/pkg/apperror"
)

// ExpenseHandler handles expense HTTP endpoints.
type ExpenseHandler struct {
	uc *usecase.ExpenseUseCase
}

// NewExpenseHandler creates a new ExpenseHandler.
func NewExpenseHandler(uc *usecase.ExpenseUseCase) *ExpenseHandler {
	return &ExpenseHandler{uc: uc}
}

// Routes registers expense routes.
func (h *ExpenseHandler) Routes() chi.Router {
	r := chi.NewRouter()

	// Categories
	r.Get("/categories", h.ListCategories)
	r.Post("/categories", h.CreateCategory)
	r.Put("/categories/{id}", h.UpdateCategory)

	// Expenses CRUD
	r.Post("/", h.CreateExpense)
	r.Get("/", h.ListExpenses)
	r.Get("/stats", h.GetStats)
	r.Get("/report", h.GetReport)
	r.Get("/trend", h.GetMonthlyTrend)
	r.Get("/{id}", h.GetExpense)
	r.Put("/{id}", h.UpdateExpense)
	r.Delete("/{id}", h.DeleteExpense)

	// Profit & Loss
	r.Get("/profit-loss", h.GetProfitLoss)

	return r
}

// --- Category Handlers ---

func (h *ExpenseHandler) ListCategories(w http.ResponseWriter, r *http.Request) {
	activeOnly := r.URL.Query().Get("active_only") == "true"
	cats, err := h.uc.ListCategories(r.Context(), activeOnly)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, apiResponse{Success: true, Data: cats})
}

func (h *ExpenseHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, &apperror.Error{Kind: apperror.KindValidation, Message: "invalid request body"})
		return
	}
	cat, err := h.uc.CreateCategory(r.Context(), input.Name, input.Description)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, apiResponse{Success: true, Data: cat})
}

func (h *ExpenseHandler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, &apperror.Error{Kind: apperror.KindValidation, Message: "invalid category ID"})
		return
	}
	var input struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		IsActive    bool   `json:"is_active"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, &apperror.Error{Kind: apperror.KindValidation, Message: "invalid request body"})
		return
	}
	cat, err := h.uc.UpdateCategory(r.Context(), id, input.Name, input.Description, input.IsActive)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, apiResponse{Success: true, Data: cat})
}

// --- Expense Handlers ---

func (h *ExpenseHandler) CreateExpense(w http.ResponseWriter, r *http.Request) {
	var input usecase.CreateExpenseInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, &apperror.Error{Kind: apperror.KindValidation, Message: "invalid request body"})
		return
	}
	exp, err := h.uc.CreateExpense(r.Context(), input)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, apiResponse{Success: true, Data: exp})
}

func (h *ExpenseHandler) ListExpenses(w http.ResponseWriter, r *http.Request) {
	input := usecase.ListExpensesInput{
		CategoryID:    r.URL.Query().Get("category_id"),
		Status:        r.URL.Query().Get("status"),
		PaymentMethod: r.URL.Query().Get("payment_method"),
		DateFrom:      r.URL.Query().Get("date_from"),
		DateTo:        r.URL.Query().Get("date_to"),
		Search:        r.URL.Query().Get("search"),
		Page:          queryInt(r, "page", 1),
		PerPage:       queryInt(r, "per_page", 20),
	}

	output, err := h.uc.ListExpenses(r.Context(), input)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, apiResponse{
		Success: true,
		Data:    output.Expenses,
		Meta: &apiMeta{
			Page:       output.Page,
			PerPage:    output.PerPage,
			Total:      output.Total,
			TotalPages: output.TotalPages,
		},
	})
}

func (h *ExpenseHandler) GetExpense(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, &apperror.Error{Kind: apperror.KindValidation, Message: "invalid expense ID"})
		return
	}
	exp, err := h.uc.GetExpense(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, apiResponse{Success: true, Data: exp})
}

func (h *ExpenseHandler) UpdateExpense(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, &apperror.Error{Kind: apperror.KindValidation, Message: "invalid expense ID"})
		return
	}
	var input usecase.UpdateExpenseInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, &apperror.Error{Kind: apperror.KindValidation, Message: "invalid request body"})
		return
	}
	exp, err := h.uc.UpdateExpense(r.Context(), id, input)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, apiResponse{Success: true, Data: exp})
}

func (h *ExpenseHandler) DeleteExpense(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, &apperror.Error{Kind: apperror.KindValidation, Message: "invalid expense ID"})
		return
	}
	if err := h.uc.DeleteExpense(r.Context(), id); err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, apiResponse{Success: true, Data: map[string]string{"message": "expense deleted"}})
}

// --- Reporting Handlers ---

func (h *ExpenseHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.uc.GetExpenseStats(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, apiResponse{Success: true, Data: stats})
}

func (h *ExpenseHandler) GetReport(w http.ResponseWriter, r *http.Request) {
	input := usecase.ExpenseReportInput{
		DateFrom: r.URL.Query().Get("date_from"),
		DateTo:   r.URL.Query().Get("date_to"),
	}
	report, err := h.uc.GetExpenseReport(r.Context(), input)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, apiResponse{Success: true, Data: report})
}

func (h *ExpenseHandler) GetProfitLoss(w http.ResponseWriter, r *http.Request) {
	input := usecase.ProfitLossInput{
		DateFrom: r.URL.Query().Get("date_from"),
		DateTo:   r.URL.Query().Get("date_to"),
	}
	pl, err := h.uc.GetProfitLoss(r.Context(), input)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, apiResponse{Success: true, Data: pl})
}

func (h *ExpenseHandler) GetMonthlyTrend(w http.ResponseWriter, r *http.Request) {
	months := queryInt(r, "months", 6)
	trends, err := h.uc.GetMonthlyTrend(r.Context(), months)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, apiResponse{Success: true, Data: trends})
}
