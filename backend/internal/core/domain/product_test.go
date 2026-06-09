package domain

import (
	"testing"

	"github.com/google/uuid"
)

func TestProduct_Validate(t *testing.T) {
	tests := []struct {
		name    string
		product *Product
		wantErr error
	}{
		{
			name:    "valid product",
			product: NewProduct("Shampoo", "hair_care", "L'Oreal", "ml", 200, 350, 5, 30),
			wantErr: nil,
		},
		{
			name:    "empty name",
			product: NewProduct("", "hair_care", "Brand", "ml", 200, 350, 5, 30),
			wantErr: ErrProductNameRequired,
		},
		{
			name:    "invalid category",
			product: NewProduct("Shampoo", "invalid_cat", "Brand", "ml", 200, 350, 5, 30),
			wantErr: ErrProductInvalidCategory,
		},
		{
			name:    "negative purchase price",
			product: NewProduct("Shampoo", "hair_care", "Brand", "ml", -10, 350, 5, 30),
			wantErr: ErrProductInvalidPrice,
		},
		{
			name:    "negative selling price",
			product: NewProduct("Shampoo", "hair_care", "Brand", "ml", 200, -5, 5, 30),
			wantErr: ErrProductInvalidPrice,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.product.Validate()
			if err != tt.wantErr {
				t.Errorf("Validate() = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestProduct_IsLowStock(t *testing.T) {
	p := NewProduct("Test", "retail", "", "pcs", 100, 200, 10, 50)
	p.CurrentStock = 5
	if !p.IsLowStock() {
		t.Error("expected IsLowStock=true when stock 5 < min 10")
	}

	p.CurrentStock = 15
	if p.IsLowStock() {
		t.Error("expected IsLowStock=false when stock 15 >= min 10")
	}
}

func TestProduct_InventoryValue(t *testing.T) {
	p := NewProduct("Test", "retail", "", "pcs", 100, 200, 10, 50)
	p.CurrentStock = 25
	val := p.InventoryValue()
	if val != 2500 {
		t.Errorf("expected InventoryValue 2500, got %f", val)
	}
}

func TestStockTransaction_Validate(t *testing.T) {
	tests := []struct {
		name    string
		tx      *StockTransaction
		wantErr error
	}{
		{
			name:    "valid transaction",
			tx:      NewStockTransaction(uuid.New(), "purchase", 10, 100, "", "", "test", "2026-06-01"),
			wantErr: nil,
		},
		{
			name:    "nil product ID",
			tx:      NewStockTransaction(uuid.Nil, "purchase", 10, 100, "", "", "test", "2026-06-01"),
			wantErr: ErrProductRequired,
		},
		{
			name:    "invalid type",
			tx:      NewStockTransaction(uuid.New(), "invalid", 10, 100, "", "", "test", "2026-06-01"),
			wantErr: ErrInvalidTransactionType,
		},
		{
			name:    "zero quantity",
			tx:      NewStockTransaction(uuid.New(), "purchase", 0, 100, "", "", "test", "2026-06-01"),
			wantErr: ErrInvalidQuantity,
		},
		{
			name:    "empty date",
			tx:      NewStockTransaction(uuid.New(), "purchase", 10, 100, "", "", "test", ""),
			wantErr: ErrTransactionDateRequired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.tx.Validate()
			if err != tt.wantErr {
				t.Errorf("Validate() = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestPurchaseEntry_Validate(t *testing.T) {
	tests := []struct {
		name    string
		pe      *PurchaseEntry
		wantErr error
	}{
		{
			name: "valid purchase",
			pe: func() *PurchaseEntry {
				pe := NewPurchaseEntry("Vendor", "INV-1", "2026-06-01", "notes")
				pe.Items = []PurchaseItem{{ID: uuid.New(), ProductID: uuid.New(), Quantity: 5, UnitPrice: 100}}
				return pe
			}(),
			wantErr: nil,
		},
		{
			name:    "empty vendor",
			pe:      NewPurchaseEntry("", "INV-1", "2026-06-01", ""),
			wantErr: ErrPurchaseVendorRequired,
		},
		{
			name:    "empty date",
			pe:      NewPurchaseEntry("Vendor", "INV-1", "", ""),
			wantErr: ErrPurchaseDateRequired,
		},
		{
			name:    "no items",
			pe:      NewPurchaseEntry("Vendor", "INV-1", "2026-06-01", ""),
			wantErr: ErrPurchaseItemsRequired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.pe.Validate()
			if err != tt.wantErr {
				t.Errorf("Validate() = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestPurchaseEntry_CalculateTotal(t *testing.T) {
	pe := NewPurchaseEntry("Vendor", "INV-1", "2026-06-01", "")
	pe.Items = []PurchaseItem{
		{ID: uuid.New(), ProductID: uuid.New(), Quantity: 5, UnitPrice: 100},
		{ID: uuid.New(), ProductID: uuid.New(), Quantity: 3, UnitPrice: 250},
	}
	pe.CalculateTotal()

	if pe.TotalAmount != 1250 {
		t.Errorf("expected total 1250, got %f", pe.TotalAmount)
	}
	if pe.Items[0].LineTotal != 500 {
		t.Errorf("expected item[0] lineTotal 500, got %f", pe.Items[0].LineTotal)
	}
	if pe.Items[1].LineTotal != 750 {
		t.Errorf("expected item[1] lineTotal 750, got %f", pe.Items[1].LineTotal)
	}
}
