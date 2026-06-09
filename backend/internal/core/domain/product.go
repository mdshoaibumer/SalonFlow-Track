package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/salonflow/salonflow-track/pkg/uid"
)

// Product categories.
var ValidProductCategories = map[string]bool{
	"hair_care": true,
	"facial":    true,
	"spa":       true,
	"coloring":  true,
	"treatment": true,
	"retail":    true,
	"equipment": true,
	"other":     true,
}

// Product statuses.
var ValidProductStatuses = map[string]bool{
	"active":       true,
	"inactive":     true,
	"discontinued": true,
}

// Stock transaction types.
var ValidTransactionTypes = map[string]bool{
	"purchase":    true,
	"consumption": true,
	"sale":        true,
	"adjustment":  true,
	"return":      true,
	"damage":      true,
}

// Product represents an inventory product.
type Product struct {
	ID            uuid.UUID `json:"id"`
	ProductCode   string    `json:"product_code"`
	Name          string    `json:"name"`
	Category      string    `json:"category"`
	Brand         string    `json:"brand"`
	Unit          string    `json:"unit"`
	SKU           string    `json:"sku"`
	PurchasePrice float64   `json:"purchase_price"`
	SellingPrice  float64   `json:"selling_price"`
	CurrentStock  float64   `json:"current_stock"`
	MinimumStock  float64   `json:"minimum_stock"`
	MaximumStock  float64   `json:"maximum_stock"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// NewProduct creates a new Product with defaults.
func NewProduct(name, category, brand, unit string, purchasePrice, sellingPrice, minStock, maxStock float64) *Product {
	now := time.Now().UTC()
	return &Product{
		ID:            uid.New(),
		Name:          name,
		Category:      category,
		Brand:         brand,
		Unit:          unit,
		PurchasePrice: purchasePrice,
		SellingPrice:  sellingPrice,
		CurrentStock:  0,
		MinimumStock:  minStock,
		MaximumStock:  maxStock,
		Status:        "active",
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}

// Validate validates a Product.
func (p *Product) Validate() error {
	if p.Name == "" {
		return ErrProductNameRequired
	}
	if !ValidProductCategories[p.Category] {
		return ErrProductInvalidCategory
	}
	if p.PurchasePrice < 0 {
		return ErrProductInvalidPrice
	}
	if p.SellingPrice < 0 {
		return ErrProductInvalidPrice
	}
	return nil
}

// IsLowStock returns true if current stock is below minimum.
func (p *Product) IsLowStock() bool {
	return p.CurrentStock < p.MinimumStock
}

// InventoryValue returns the total value of current stock.
func (p *Product) InventoryValue() float64 {
	return p.CurrentStock * p.PurchasePrice
}

// StockTransaction represents a stock movement.
type StockTransaction struct {
	ID              uuid.UUID `json:"id"`
	ProductID       uuid.UUID `json:"product_id"`
	ProductName     string    `json:"product_name,omitempty"`
	TransactionType string    `json:"transaction_type"`
	Quantity        float64   `json:"quantity"`
	UnitCost        float64   `json:"unit_cost"`
	ReferenceType   string    `json:"reference_type"`
	ReferenceID     string    `json:"reference_id"`
	Notes           string    `json:"notes"`
	TransactionDate string    `json:"transaction_date"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// NewStockTransaction creates a new stock transaction.
func NewStockTransaction(productID uuid.UUID, txType string, qty, unitCost float64, refType, refID, notes, txDate string) *StockTransaction {
	now := time.Now().UTC()
	return &StockTransaction{
		ID:              uid.New(),
		ProductID:       productID,
		TransactionType: txType,
		Quantity:        qty,
		UnitCost:        unitCost,
		ReferenceType:   refType,
		ReferenceID:     refID,
		Notes:           notes,
		TransactionDate: txDate,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
}

// Validate validates a StockTransaction.
func (st *StockTransaction) Validate() error {
	if st.ProductID == uuid.Nil {
		return ErrProductRequired
	}
	if !ValidTransactionTypes[st.TransactionType] {
		return ErrInvalidTransactionType
	}
	if st.Quantity == 0 {
		return ErrInvalidQuantity
	}
	if st.TransactionDate == "" {
		return ErrTransactionDateRequired
	}
	return nil
}

// PurchaseEntry represents a purchase order/invoice.
type PurchaseEntry struct {
	ID             uuid.UUID      `json:"id"`
	PurchaseNumber string         `json:"purchase_number"`
	VendorName     string         `json:"vendor_name"`
	InvoiceNumber  string         `json:"invoice_number"`
	PurchaseDate   string         `json:"purchase_date"`
	TotalAmount    float64        `json:"total_amount"`
	Notes          string         `json:"notes"`
	Items          []PurchaseItem `json:"items,omitempty"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
}

// NewPurchaseEntry creates a new PurchaseEntry.
func NewPurchaseEntry(vendorName, invoiceNumber, purchaseDate, notes string) *PurchaseEntry {
	now := time.Now().UTC()
	return &PurchaseEntry{
		ID:            uid.New(),
		VendorName:    vendorName,
		InvoiceNumber: invoiceNumber,
		PurchaseDate:  purchaseDate,
		Notes:         notes,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}

// Validate validates a PurchaseEntry.
func (pe *PurchaseEntry) Validate() error {
	if pe.VendorName == "" {
		return ErrPurchaseVendorRequired
	}
	if pe.PurchaseDate == "" {
		return ErrPurchaseDateRequired
	}
	if len(pe.Items) == 0 {
		return ErrPurchaseItemsRequired
	}
	return nil
}

// CalculateTotal recalculates the total from items.
func (pe *PurchaseEntry) CalculateTotal() {
	var total float64
	for i := range pe.Items {
		pe.Items[i].LineTotal = pe.Items[i].Quantity * pe.Items[i].UnitPrice
		total += pe.Items[i].LineTotal
	}
	pe.TotalAmount = total
}

// PurchaseItem represents a line item in a purchase entry.
type PurchaseItem struct {
	ID              uuid.UUID `json:"id"`
	PurchaseEntryID uuid.UUID `json:"purchase_entry_id"`
	ProductID       uuid.UUID `json:"product_id"`
	ProductName     string    `json:"product_name,omitempty"`
	Quantity        float64   `json:"quantity"`
	UnitPrice       float64   `json:"unit_price"`
	LineTotal       float64   `json:"line_total"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// InventoryStats holds dashboard stats.
type InventoryStats struct {
	TotalProducts  int     `json:"total_products"`
	ActiveProducts int     `json:"active_products"`
	LowStockCount  int     `json:"low_stock_count"`
	TotalValue     float64 `json:"total_value"`
	TotalPurchases float64 `json:"total_purchases_this_month"`
}

// LowStockItem represents a product with low stock.
type LowStockItem struct {
	ProductID    string  `json:"product_id"`
	ProductCode  string  `json:"product_code"`
	ProductName  string  `json:"product_name"`
	Category     string  `json:"category"`
	CurrentStock float64 `json:"current_stock"`
	MinimumStock float64 `json:"minimum_stock"`
	Deficit      float64 `json:"deficit"`
}
