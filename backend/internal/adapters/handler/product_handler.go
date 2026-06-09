package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/salonflow/salonflow-track/internal/core/usecase"
	"github.com/salonflow/salonflow-track/pkg/apperror"
)

// ProductHandler handles product/inventory HTTP endpoints.
type ProductHandler struct {
	uc *usecase.ProductUseCase
}

// NewProductHandler creates a new ProductHandler.
func NewProductHandler(uc *usecase.ProductUseCase) *ProductHandler {
	return &ProductHandler{uc: uc}
}

// Routes registers product/inventory routes.
func (h *ProductHandler) Routes() chi.Router {
	r := chi.NewRouter()

	// Products
	r.Post("/", h.CreateProduct)
	r.Get("/", h.ListProducts)
	r.Get("/stats", h.GetStats)
	r.Get("/low-stock", h.GetLowStock)
	r.Get("/{id}", h.GetProduct)
	r.Put("/{id}", h.UpdateProduct)
	r.Delete("/{id}", h.DeleteProduct)

	// Stock
	r.Post("/stock/adjust", h.AdjustStock)
	r.Get("/stock/history", h.ListStockHistory)

	// Purchases
	r.Post("/purchases", h.CreatePurchase)
	r.Get("/purchases", h.ListPurchases)
	r.Get("/purchases/{id}", h.GetPurchase)

	return r
}

// --- Product Handlers ---

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var input usecase.CreateProductInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, &apperror.Error{Kind: apperror.KindValidation, Message: "invalid request body"})
		return
	}
	p, err := h.uc.CreateProduct(r.Context(), input)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, apiResponse{Success: true, Data: p})
}

func (h *ProductHandler) ListProducts(w http.ResponseWriter, r *http.Request) {
	input := usecase.ListProductsInput{
		Category: r.URL.Query().Get("category"),
		Status:   r.URL.Query().Get("status"),
		Search:   r.URL.Query().Get("search"),
		LowStock: r.URL.Query().Get("low_stock") == "true",
		Page:     queryInt(r, "page", 1),
		PerPage:  queryInt(r, "per_page", 20),
	}

	output, err := h.uc.ListProducts(r.Context(), input)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, apiResponse{
		Success: true,
		Data:    output.Products,
		Meta: &apiMeta{
			Page:       output.Page,
			PerPage:    output.PerPage,
			Total:      output.Total,
			TotalPages: output.TotalPages,
		},
	})
}

func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, &apperror.Error{Kind: apperror.KindValidation, Message: "invalid product ID"})
		return
	}
	p, err := h.uc.GetProduct(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, apiResponse{Success: true, Data: p})
}

func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, &apperror.Error{Kind: apperror.KindValidation, Message: "invalid product ID"})
		return
	}
	var input usecase.UpdateProductInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, &apperror.Error{Kind: apperror.KindValidation, Message: "invalid request body"})
		return
	}
	p, err := h.uc.UpdateProduct(r.Context(), id, input)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, apiResponse{Success: true, Data: p})
}

func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, &apperror.Error{Kind: apperror.KindValidation, Message: "invalid product ID"})
		return
	}
	if err := h.uc.DeleteProduct(r.Context(), id); err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, apiResponse{Success: true, Data: map[string]string{"message": "product deleted"}})
}

// --- Stock Handlers ---

func (h *ProductHandler) AdjustStock(w http.ResponseWriter, r *http.Request) {
	var input usecase.StockAdjustInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, &apperror.Error{Kind: apperror.KindValidation, Message: "invalid request body"})
		return
	}
	st, err := h.uc.AdjustStock(r.Context(), input)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, apiResponse{Success: true, Data: st})
}

func (h *ProductHandler) ListStockHistory(w http.ResponseWriter, r *http.Request) {
	input := usecase.ListStockHistoryInput{
		ProductID:       r.URL.Query().Get("product_id"),
		TransactionType: r.URL.Query().Get("transaction_type"),
		DateFrom:        r.URL.Query().Get("date_from"),
		DateTo:          r.URL.Query().Get("date_to"),
		Page:            queryInt(r, "page", 1),
		PerPage:         queryInt(r, "per_page", 20),
	}

	output, err := h.uc.ListStockHistory(r.Context(), input)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, apiResponse{
		Success: true,
		Data:    output.Transactions,
		Meta: &apiMeta{
			Page:       output.Page,
			PerPage:    output.PerPage,
			Total:      output.Total,
			TotalPages: output.TotalPages,
		},
	})
}

// --- Purchase Handlers ---

func (h *ProductHandler) CreatePurchase(w http.ResponseWriter, r *http.Request) {
	var input usecase.CreatePurchaseInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, &apperror.Error{Kind: apperror.KindValidation, Message: "invalid request body"})
		return
	}
	pe, err := h.uc.CreatePurchase(r.Context(), input)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, apiResponse{Success: true, Data: pe})
}

func (h *ProductHandler) ListPurchases(w http.ResponseWriter, r *http.Request) {
	input := usecase.ListPurchasesInput{
		DateFrom: r.URL.Query().Get("date_from"),
		DateTo:   r.URL.Query().Get("date_to"),
		Page:     queryInt(r, "page", 1),
		PerPage:  queryInt(r, "per_page", 20),
	}

	output, err := h.uc.ListPurchases(r.Context(), input)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, apiResponse{
		Success: true,
		Data:    output.Purchases,
		Meta: &apiMeta{
			Page:       output.Page,
			PerPage:    output.PerPage,
			Total:      output.Total,
			TotalPages: output.TotalPages,
		},
	})
}

func (h *ProductHandler) GetPurchase(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, &apperror.Error{Kind: apperror.KindValidation, Message: "invalid purchase ID"})
		return
	}
	pe, err := h.uc.GetPurchase(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, apiResponse{Success: true, Data: pe})
}

// --- Reporting Handlers ---

func (h *ProductHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.uc.GetInventoryStats(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, apiResponse{Success: true, Data: stats})
}

func (h *ProductHandler) GetLowStock(w http.ResponseWriter, r *http.Request) {
	items, err := h.uc.GetLowStockProducts(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, apiResponse{Success: true, Data: items})
}
