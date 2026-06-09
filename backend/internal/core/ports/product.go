package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/salonflow/salonflow-track/internal/core/domain"
)

// ProductFilter holds filtering options for listing products.
type ProductFilter struct {
	Category string
	Status   string
	Search   string
	LowStock bool
	Limit    int
	Offset   int
}

// StockTransactionFilter holds filtering for stock history.
type StockTransactionFilter struct {
	ProductID       string
	TransactionType string
	DateFrom        string
	DateTo          string
	Limit           int
	Offset          int
}

// ProductRepository defines the contract for product/inventory persistence.
type ProductRepository interface {
	// Products
	CreateProduct(ctx context.Context, p *domain.Product) error
	GetProductByID(ctx context.Context, id uuid.UUID) (*domain.Product, error)
	UpdateProduct(ctx context.Context, p *domain.Product) error
	DeleteProduct(ctx context.Context, id uuid.UUID) error
	ListProducts(ctx context.Context, filter ProductFilter) ([]domain.Product, int, error)
	UpdateStock(ctx context.Context, productID uuid.UUID, delta float64) error

	// Stock Transactions
	CreateStockTransaction(ctx context.Context, st *domain.StockTransaction) error
	ListStockTransactions(ctx context.Context, filter StockTransactionFilter) ([]domain.StockTransaction, int, error)

	// Purchases
	CreatePurchaseEntry(ctx context.Context, pe *domain.PurchaseEntry) error
	GetPurchaseByID(ctx context.Context, id uuid.UUID) (*domain.PurchaseEntry, error)
	ListPurchases(ctx context.Context, dateFrom, dateTo string, limit, offset int) ([]domain.PurchaseEntry, int, error)

	// Reporting
	GetLowStockProducts(ctx context.Context) ([]domain.LowStockItem, error)
	GetInventoryStats(ctx context.Context) (*domain.InventoryStats, error)
	GetInventoryValue(ctx context.Context) (float64, error)

	// Code generation
	NextProductCode(ctx context.Context) (string, error)
	NextPurchaseNumber(ctx context.Context) (string, error)
}
