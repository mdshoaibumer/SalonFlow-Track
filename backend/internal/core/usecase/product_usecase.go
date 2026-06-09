package usecase

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/salonflow/salonflow-track/internal/core/domain"
	"github.com/salonflow/salonflow-track/internal/core/ports"
	"github.com/salonflow/salonflow-track/pkg/apperror"
	"github.com/salonflow/salonflow-track/pkg/uid"
)

// ProductUseCase handles product/inventory business logic.
type ProductUseCase struct {
	productRepo ports.ProductRepository
	log         *slog.Logger
}

// NewProductUseCase creates a new ProductUseCase.
func NewProductUseCase(productRepo ports.ProductRepository, log *slog.Logger) *ProductUseCase {
	return &ProductUseCase{productRepo: productRepo, log: log}
}

// --- Input/Output DTOs ---

type CreateProductInput struct {
	Name          string  `json:"name"`
	Category      string  `json:"category"`
	Brand         string  `json:"brand"`
	Unit          string  `json:"unit"`
	SKU           string  `json:"sku"`
	PurchasePrice float64 `json:"purchase_price"`
	SellingPrice  float64 `json:"selling_price"`
	MinimumStock  float64 `json:"minimum_stock"`
	MaximumStock  float64 `json:"maximum_stock"`
}

type UpdateProductInput struct {
	Name          string  `json:"name"`
	Category      string  `json:"category"`
	Brand         string  `json:"brand"`
	Unit          string  `json:"unit"`
	SKU           string  `json:"sku"`
	PurchasePrice float64 `json:"purchase_price"`
	SellingPrice  float64 `json:"selling_price"`
	MinimumStock  float64 `json:"minimum_stock"`
	MaximumStock  float64 `json:"maximum_stock"`
	Status        string  `json:"status"`
}

type ListProductsInput struct {
	Category string
	Status   string
	Search   string
	LowStock bool
	Page     int
	PerPage  int
}

type ListProductsOutput struct {
	Products   []domain.Product `json:"products"`
	Total      int              `json:"total"`
	Page       int              `json:"page"`
	PerPage    int              `json:"per_page"`
	TotalPages int              `json:"total_pages"`
}

type CreatePurchaseInput struct {
	VendorName    string              `json:"vendor_name"`
	InvoiceNumber string              `json:"invoice_number"`
	PurchaseDate  string              `json:"purchase_date"`
	Notes         string              `json:"notes"`
	Items         []PurchaseItemInput `json:"items"`
}

type PurchaseItemInput struct {
	ProductID string  `json:"product_id"`
	Quantity  float64 `json:"quantity"`
	UnitPrice float64 `json:"unit_price"`
}

type StockAdjustInput struct {
	ProductID       string  `json:"product_id"`
	TransactionType string  `json:"transaction_type"`
	Quantity        float64 `json:"quantity"`
	Notes           string  `json:"notes"`
}

type ListStockHistoryInput struct {
	ProductID       string
	TransactionType string
	DateFrom        string
	DateTo          string
	Page            int
	PerPage         int
}

type ListStockHistoryOutput struct {
	Transactions []domain.StockTransaction `json:"transactions"`
	Total        int                       `json:"total"`
	Page         int                       `json:"page"`
	PerPage      int                       `json:"per_page"`
	TotalPages   int                       `json:"total_pages"`
}

type ListPurchasesInput struct {
	DateFrom string
	DateTo   string
	Page     int
	PerPage  int
}

type ListPurchasesOutput struct {
	Purchases  []domain.PurchaseEntry `json:"purchases"`
	Total      int                    `json:"total"`
	Page       int                    `json:"page"`
	PerPage    int                    `json:"per_page"`
	TotalPages int                    `json:"total_pages"`
}

// --- Product CRUD ---

func (uc *ProductUseCase) CreateProduct(ctx context.Context, input CreateProductInput) (*domain.Product, error) {
	p := domain.NewProduct(input.Name, input.Category, input.Brand, input.Unit,
		input.PurchasePrice, input.SellingPrice, input.MinimumStock, input.MaximumStock)
	p.SKU = input.SKU

	if err := p.Validate(); err != nil {
		return nil, &apperror.Error{Kind: apperror.KindValidation, Message: err.Error()}
	}

	code, err := uc.productRepo.NextProductCode(ctx)
	if err != nil {
		return nil, err
	}
	p.ProductCode = code

	if err := uc.productRepo.CreateProduct(ctx, p); err != nil {
		return nil, err
	}

	uc.log.Info("product created", "id", p.ID, "code", p.ProductCode, "name", p.Name)
	return p, nil
}

func (uc *ProductUseCase) GetProduct(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	return uc.productRepo.GetProductByID(ctx, id)
}

func (uc *ProductUseCase) UpdateProduct(ctx context.Context, id uuid.UUID, input UpdateProductInput) (*domain.Product, error) {
	p, err := uc.productRepo.GetProductByID(ctx, id)
	if err != nil {
		return nil, err
	}

	p.Name = input.Name
	p.Category = input.Category
	p.Brand = input.Brand
	p.Unit = input.Unit
	p.SKU = input.SKU
	p.PurchasePrice = input.PurchasePrice
	p.SellingPrice = input.SellingPrice
	p.MinimumStock = input.MinimumStock
	p.MaximumStock = input.MaximumStock
	if input.Status != "" && domain.ValidProductStatuses[input.Status] {
		p.Status = input.Status
	}

	if err := p.Validate(); err != nil {
		return nil, &apperror.Error{Kind: apperror.KindValidation, Message: err.Error()}
	}

	if err := uc.productRepo.UpdateProduct(ctx, p); err != nil {
		return nil, err
	}

	uc.log.Info("product updated", "id", p.ID)
	return p, nil
}

func (uc *ProductUseCase) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	_, err := uc.productRepo.GetProductByID(ctx, id)
	if err != nil {
		return err
	}
	if err := uc.productRepo.DeleteProduct(ctx, id); err != nil {
		return err
	}
	uc.log.Info("product deleted", "id", id)
	return nil
}

func (uc *ProductUseCase) ListProducts(ctx context.Context, input ListProductsInput) (*ListProductsOutput, error) {
	if input.Page <= 0 {
		input.Page = 1
	}
	if input.PerPage <= 0 {
		input.PerPage = 20
	}
	if input.PerPage > 100 {
		input.PerPage = 100
	}

	filter := ports.ProductFilter{
		Category: input.Category,
		Status:   input.Status,
		Search:   input.Search,
		LowStock: input.LowStock,
		Limit:    input.PerPage,
		Offset:   (input.Page - 1) * input.PerPage,
	}

	products, total, err := uc.productRepo.ListProducts(ctx, filter)
	if err != nil {
		return nil, err
	}

	totalPages := total / input.PerPage
	if total%input.PerPage > 0 {
		totalPages++
	}

	return &ListProductsOutput{
		Products:   products,
		Total:      total,
		Page:       input.Page,
		PerPage:    input.PerPage,
		TotalPages: totalPages,
	}, nil
}

// --- Purchase ---

func (uc *ProductUseCase) CreatePurchase(ctx context.Context, input CreatePurchaseInput) (*domain.PurchaseEntry, error) {
	pe := domain.NewPurchaseEntry(input.VendorName, input.InvoiceNumber, input.PurchaseDate, input.Notes)

	now := time.Now().UTC()
	for _, itemInput := range input.Items {
		productID, err := uuid.Parse(itemInput.ProductID)
		if err != nil {
			return nil, &apperror.Error{Kind: apperror.KindValidation, Message: "invalid product ID in items"}
		}
		// Verify product exists
		_, err = uc.productRepo.GetProductByID(ctx, productID)
		if err != nil {
			return nil, err
		}
		item := domain.PurchaseItem{
			ID:              uid.New(),
			PurchaseEntryID: pe.ID,
			ProductID:       productID,
			Quantity:        itemInput.Quantity,
			UnitPrice:       itemInput.UnitPrice,
			LineTotal:       itemInput.Quantity * itemInput.UnitPrice,
			CreatedAt:       now,
			UpdatedAt:       now,
		}
		pe.Items = append(pe.Items, item)
	}

	if err := pe.Validate(); err != nil {
		return nil, &apperror.Error{Kind: apperror.KindValidation, Message: err.Error()}
	}

	pe.CalculateTotal()

	// Generate purchase number
	purNum, err := uc.productRepo.NextPurchaseNumber(ctx)
	if err != nil {
		return nil, err
	}
	pe.PurchaseNumber = purNum

	// Persist
	if err := uc.productRepo.CreatePurchaseEntry(ctx, pe); err != nil {
		return nil, err
	}

	// Update stock and create stock transactions
	txDate := input.PurchaseDate
	if txDate == "" {
		txDate = time.Now().UTC().Format("2006-01-02")
	}
	for _, item := range pe.Items {
		if err := uc.productRepo.UpdateStock(ctx, item.ProductID, item.Quantity); err != nil {
			uc.log.Error("failed to update stock for purchase", "product_id", item.ProductID, "error", err)
			continue
		}
		st := domain.NewStockTransaction(item.ProductID, "purchase", item.Quantity, item.UnitPrice,
			"purchase", pe.ID.String(), "Purchase: "+pe.PurchaseNumber, txDate)
		if err := uc.productRepo.CreateStockTransaction(ctx, st); err != nil {
			uc.log.Error("failed to create stock transaction", "product_id", item.ProductID, "error", err)
		}
	}

	uc.log.Info("purchase created", "id", pe.ID, "number", pe.PurchaseNumber, "total", pe.TotalAmount)
	return pe, nil
}

func (uc *ProductUseCase) GetPurchase(ctx context.Context, id uuid.UUID) (*domain.PurchaseEntry, error) {
	return uc.productRepo.GetPurchaseByID(ctx, id)
}

func (uc *ProductUseCase) ListPurchases(ctx context.Context, input ListPurchasesInput) (*ListPurchasesOutput, error) {
	if input.Page <= 0 {
		input.Page = 1
	}
	if input.PerPage <= 0 {
		input.PerPage = 20
	}
	if input.PerPage > 100 {
		input.PerPage = 100
	}

	purchases, total, err := uc.productRepo.ListPurchases(ctx, input.DateFrom, input.DateTo, input.PerPage, (input.Page-1)*input.PerPage)
	if err != nil {
		return nil, err
	}

	totalPages := total / input.PerPage
	if total%input.PerPage > 0 {
		totalPages++
	}

	return &ListPurchasesOutput{
		Purchases:  purchases,
		Total:      total,
		Page:       input.Page,
		PerPage:    input.PerPage,
		TotalPages: totalPages,
	}, nil
}

// --- Stock Operations ---

func (uc *ProductUseCase) AdjustStock(ctx context.Context, input StockAdjustInput) (*domain.StockTransaction, error) {
	productID, err := uuid.Parse(input.ProductID)
	if err != nil {
		return nil, &apperror.Error{Kind: apperror.KindValidation, Message: "invalid product ID"}
	}

	product, err := uc.productRepo.GetProductByID(ctx, productID)
	if err != nil {
		return nil, err
	}

	// For consumption/sale/damage, quantity should be negative (stock decreases)
	delta := input.Quantity
	switch input.TransactionType {
	case "consumption", "sale", "damage":
		if delta > 0 {
			delta = -delta
		}
		if product.CurrentStock+delta < 0 {
			return nil, &apperror.Error{Kind: apperror.KindBusiness, Message: domain.ErrInsufficientStock.Error()}
		}
	case "purchase", "return":
		if delta < 0 {
			delta = -delta
		}
	case "adjustment":
		// Can be positive or negative
	default:
		return nil, &apperror.Error{Kind: apperror.KindValidation, Message: domain.ErrInvalidTransactionType.Error()}
	}

	txDate := time.Now().UTC().Format("2006-01-02")
	st := domain.NewStockTransaction(productID, input.TransactionType, delta, product.PurchasePrice,
		"manual", "", input.Notes, txDate)

	if err := st.Validate(); err != nil {
		return nil, &apperror.Error{Kind: apperror.KindValidation, Message: err.Error()}
	}

	if err := uc.productRepo.UpdateStock(ctx, productID, delta); err != nil {
		return nil, err
	}

	if err := uc.productRepo.CreateStockTransaction(ctx, st); err != nil {
		return nil, err
	}

	uc.log.Info("stock adjusted", "product_id", productID, "type", input.TransactionType, "qty", delta)
	return st, nil
}

func (uc *ProductUseCase) ListStockHistory(ctx context.Context, input ListStockHistoryInput) (*ListStockHistoryOutput, error) {
	if input.Page <= 0 {
		input.Page = 1
	}
	if input.PerPage <= 0 {
		input.PerPage = 20
	}
	if input.PerPage > 100 {
		input.PerPage = 100
	}

	filter := ports.StockTransactionFilter{
		ProductID:       input.ProductID,
		TransactionType: input.TransactionType,
		DateFrom:        input.DateFrom,
		DateTo:          input.DateTo,
		Limit:           input.PerPage,
		Offset:          (input.Page - 1) * input.PerPage,
	}

	txns, total, err := uc.productRepo.ListStockTransactions(ctx, filter)
	if err != nil {
		return nil, err
	}

	totalPages := total / input.PerPage
	if total%input.PerPage > 0 {
		totalPages++
	}

	return &ListStockHistoryOutput{
		Transactions: txns,
		Total:        total,
		Page:         input.Page,
		PerPage:      input.PerPage,
		TotalPages:   totalPages,
	}, nil
}

// --- Reporting ---

func (uc *ProductUseCase) GetLowStockProducts(ctx context.Context) ([]domain.LowStockItem, error) {
	return uc.productRepo.GetLowStockProducts(ctx)
}

func (uc *ProductUseCase) GetInventoryStats(ctx context.Context) (*domain.InventoryStats, error) {
	return uc.productRepo.GetInventoryStats(ctx)
}
