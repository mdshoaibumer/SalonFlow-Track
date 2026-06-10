package sqlite

import (
	"context"
	"database/sql"
	"log/slog"
	"os"
	"testing"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"github.com/salonflow/salonflow-track/internal/core/domain"
	"github.com/salonflow/salonflow-track/internal/core/ports"
)

func setupProductTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(`
		CREATE TABLE products (
			id TEXT PRIMARY KEY,
			product_code TEXT NOT NULL UNIQUE,
			name TEXT NOT NULL,
			category TEXT NOT NULL CHECK (category IN ('hair_care','facial','spa','coloring','treatment','retail','equipment','other')),
			brand TEXT NOT NULL DEFAULT '',
			unit TEXT NOT NULL DEFAULT 'pcs',
			sku TEXT NOT NULL DEFAULT '',
			purchase_price REAL NOT NULL DEFAULT 0,
			selling_price REAL NOT NULL DEFAULT 0,
			current_stock REAL NOT NULL DEFAULT 0,
			minimum_stock REAL NOT NULL DEFAULT 0,
			maximum_stock REAL NOT NULL DEFAULT 0,
			status TEXT NOT NULL DEFAULT 'active' CHECK (status IN ('active','inactive','discontinued')),
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL
		);
		CREATE TABLE stock_transactions (
			id TEXT PRIMARY KEY,
			product_id TEXT NOT NULL REFERENCES products(id),
			transaction_type TEXT NOT NULL CHECK (transaction_type IN ('purchase','consumption','sale','adjustment','return','damage')),
			quantity REAL NOT NULL,
			unit_cost REAL NOT NULL DEFAULT 0,
			reference_type TEXT NOT NULL DEFAULT '',
			reference_id TEXT NOT NULL DEFAULT '',
			notes TEXT NOT NULL DEFAULT '',
			transaction_date DATE NOT NULL,
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL
		);
		CREATE TABLE purchase_entries (
			id TEXT PRIMARY KEY,
			purchase_number TEXT NOT NULL UNIQUE,
			vendor_name TEXT NOT NULL,
			invoice_number TEXT NOT NULL DEFAULT '',
			purchase_date DATE NOT NULL,
			total_amount REAL NOT NULL DEFAULT 0,
			notes TEXT NOT NULL DEFAULT '',
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL
		);
		CREATE TABLE purchase_items (
			id TEXT PRIMARY KEY,
			purchase_entry_id TEXT NOT NULL REFERENCES purchase_entries(id),
			product_id TEXT NOT NULL REFERENCES products(id),
			quantity REAL NOT NULL,
			unit_price REAL NOT NULL,
			line_total REAL NOT NULL,
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL
		);
		CREATE TABLE product_code_seq (
			prefix TEXT PRIMARY KEY,
			seq INTEGER NOT NULL DEFAULT 0
		);
		INSERT INTO product_code_seq (prefix, seq) VALUES ('PRD', 0);
		INSERT INTO product_code_seq (prefix, seq) VALUES ('PUR', 0);
	`)
	if err != nil {
		t.Fatal(err)
	}
	return db
}

func TestProductRepository_CreateAndGet(t *testing.T) {
	db := setupProductTestDB(t)
	defer db.Close()
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	repo := NewProductRepository(db, log)
	ctx := context.Background()

	p := domain.NewProduct("Shampoo Pro", "hair_care", "L'Oreal", "ml", 350, 500, 10, 50)
	code, err := repo.NextProductCode(ctx)
	if err != nil {
		t.Fatalf("NextProductCode: %v", err)
	}
	p.ProductCode = code

	err = repo.CreateProduct(ctx, p)
	if err != nil {
		t.Fatalf("CreateProduct: %v", err)
	}

	got, err := repo.GetProductByID(ctx, p.ID)
	if err != nil {
		t.Fatalf("GetProductByID: %v", err)
	}
	if got.Name != "Shampoo Pro" {
		t.Errorf("expected name 'Shampoo Pro', got '%s'", got.Name)
	}
	if got.ProductCode != "PRD-000001" {
		t.Errorf("expected code PRD-000001, got %s", got.ProductCode)
	}
	if got.Category != "hair_care" {
		t.Errorf("expected category hair_care, got %s", got.Category)
	}
	if got.PurchasePrice != 350 {
		t.Errorf("expected purchase price 350, got %f", got.PurchasePrice)
	}
	if got.SellingPrice != 500 {
		t.Errorf("expected selling price 500, got %f", got.SellingPrice)
	}
}

func TestProductRepository_UpdateProduct(t *testing.T) {
	db := setupProductTestDB(t)
	defer db.Close()
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	repo := NewProductRepository(db, log)
	ctx := context.Background()

	p := domain.NewProduct("Conditioner", "hair_care", "Matrix", "ml", 250, 400, 5, 30)
	p.ProductCode = "PRD-000001"
	repo.CreateProduct(ctx, p)

	p.Name = "Conditioner Pro"
	p.SellingPrice = 450
	err := repo.UpdateProduct(ctx, p)
	if err != nil {
		t.Fatalf("UpdateProduct: %v", err)
	}

	got, _ := repo.GetProductByID(ctx, p.ID)
	if got.Name != "Conditioner Pro" {
		t.Errorf("expected name 'Conditioner Pro', got '%s'", got.Name)
	}
	if got.SellingPrice != 450 {
		t.Errorf("expected selling price 450, got %f", got.SellingPrice)
	}
}

func TestProductRepository_DeleteProduct(t *testing.T) {
	db := setupProductTestDB(t)
	defer db.Close()
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	repo := NewProductRepository(db, log)
	ctx := context.Background()

	p := domain.NewProduct("Gel", "hair_care", "Gatsby", "pcs", 100, 180, 10, 100)
	p.ProductCode = "PRD-000001"
	repo.CreateProduct(ctx, p)

	err := repo.DeleteProduct(ctx, p.ID)
	if err != nil {
		t.Fatalf("DeleteProduct: %v", err)
	}

	_, err = repo.GetProductByID(ctx, p.ID)
	if err == nil {
		t.Fatal("expected error after delete, got nil")
	}
}

func TestProductRepository_ListProducts(t *testing.T) {
	db := setupProductTestDB(t)
	defer db.Close()
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	repo := NewProductRepository(db, log)
	ctx := context.Background()

	p1 := domain.NewProduct("Shampoo A", "hair_care", "Brand A", "ml", 200, 350, 10, 50)
	p1.ProductCode = "PRD-000001"
	p2 := domain.NewProduct("Facial Cream", "facial", "Brand B", "pcs", 500, 750, 5, 20)
	p2.ProductCode = "PRD-000002"
	p3 := domain.NewProduct("Spa Oil", "spa", "Brand C", "ml", 300, 600, 3, 15)
	p3.ProductCode = "PRD-000003"
	repo.CreateProduct(ctx, p1)
	repo.CreateProduct(ctx, p2)
	repo.CreateProduct(ctx, p3)

	// All products
	products, total, err := repo.ListProducts(ctx, ports.ProductFilter{Limit: 10})
	if err != nil {
		t.Fatalf("ListProducts: %v", err)
	}
	if total != 3 {
		t.Errorf("expected total 3, got %d", total)
	}
	if len(products) != 3 {
		t.Errorf("expected 3 products, got %d", len(products))
	}

	// Filter by category
	products, total, err = repo.ListProducts(ctx, ports.ProductFilter{Category: "facial", Limit: 10})
	if err != nil {
		t.Fatalf("ListProducts by category: %v", err)
	}
	if total != 1 {
		t.Errorf("expected 1 facial product, got %d", total)
	}
	if len(products) > 0 && products[0].Name != "Facial Cream" {
		t.Errorf("expected 'Facial Cream', got '%s'", products[0].Name)
	}

	// Search
	_, total, err = repo.ListProducts(ctx, ports.ProductFilter{Search: "Shampoo", Limit: 10})
	if err != nil {
		t.Fatalf("ListProducts search: %v", err)
	}
	if total != 1 {
		t.Errorf("expected 1 shampoo search result, got %d", total)
	}
}

func TestProductRepository_StockAndTransactions(t *testing.T) {
	db := setupProductTestDB(t)
	defer db.Close()
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	repo := NewProductRepository(db, log)
	ctx := context.Background()

	p := domain.NewProduct("Test Product", "retail", "TestBrand", "pcs", 100, 200, 5, 50)
	p.ProductCode = "PRD-000001"
	repo.CreateProduct(ctx, p)

	// Update stock +10
	err := repo.UpdateStock(ctx, p.ID, 10)
	if err != nil {
		t.Fatalf("UpdateStock +10: %v", err)
	}
	got, _ := repo.GetProductByID(ctx, p.ID)
	if got.CurrentStock != 10 {
		t.Errorf("expected stock 10 after +10, got %f", got.CurrentStock)
	}

	// Update stock -3
	err = repo.UpdateStock(ctx, p.ID, -3)
	if err != nil {
		t.Fatalf("UpdateStock -3: %v", err)
	}
	got, _ = repo.GetProductByID(ctx, p.ID)
	if got.CurrentStock != 7 {
		t.Errorf("expected stock 7 after -3, got %f", got.CurrentStock)
	}

	// Create stock transaction
	tx := domain.NewStockTransaction(p.ID, "purchase", 10, 100, "", "", "initial stock", "2026-06-01")
	err = repo.CreateStockTransaction(ctx, tx)
	if err != nil {
		t.Fatalf("CreateStockTransaction: %v", err)
	}

	// List stock transactions
	txns, total, err := repo.ListStockTransactions(ctx, ports.StockTransactionFilter{ProductID: p.ID.String(), Limit: 10})
	if err != nil {
		t.Fatalf("ListStockTransactions: %v", err)
	}
	if total != 1 {
		t.Errorf("expected 1 transaction, got %d", total)
	}
	if len(txns) > 0 && txns[0].TransactionType != "purchase" {
		t.Errorf("expected transaction type 'purchase', got '%s'", txns[0].TransactionType)
	}
}

func TestProductRepository_PurchaseEntry(t *testing.T) {
	db := setupProductTestDB(t)
	defer db.Close()
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	repo := NewProductRepository(db, log)
	ctx := context.Background()

	// Create a product first
	p := domain.NewProduct("Hair Dye", "coloring", "Garnier", "pcs", 150, 250, 10, 100)
	p.ProductCode = "PRD-000001"
	repo.CreateProduct(ctx, p)

	// Create purchase entry
	pe := domain.NewPurchaseEntry("Beauty Supplies Co", "INV-001", "2026-06-15", "Monthly restock")
	purNum, err := repo.NextPurchaseNumber(ctx)
	if err != nil {
		t.Fatalf("NextPurchaseNumber: %v", err)
	}
	pe.PurchaseNumber = purNum
	pe.Items = []domain.PurchaseItem{
		{ID: uuid.New(), ProductID: p.ID, Quantity: 20, UnitPrice: 150, LineTotal: 3000},
	}
	pe.TotalAmount = 3000

	err = repo.CreatePurchaseEntry(ctx, pe)
	if err != nil {
		t.Fatalf("CreatePurchaseEntry: %v", err)
	}

	got, err := repo.GetPurchaseByID(ctx, pe.ID)
	if err != nil {
		t.Fatalf("GetPurchaseByID: %v", err)
	}
	if got.PurchaseNumber != "PUR-000001" {
		t.Errorf("expected PUR-000001, got %s", got.PurchaseNumber)
	}
	if got.VendorName != "Beauty Supplies Co" {
		t.Errorf("expected vendor 'Beauty Supplies Co', got '%s'", got.VendorName)
	}
	if got.TotalAmount != 3000 {
		t.Errorf("expected total 3000, got %f", got.TotalAmount)
	}
	if len(got.Items) != 1 {
		t.Errorf("expected 1 item, got %d", len(got.Items))
	}

	// List purchases
	purchases, total, err := repo.ListPurchases(ctx, "", "", 10, 0)
	if err != nil {
		t.Fatalf("ListPurchases: %v", err)
	}
	if total != 1 {
		t.Errorf("expected 1 purchase, got %d", total)
	}
	if len(purchases) != 1 {
		t.Errorf("expected 1 result, got %d", len(purchases))
	}
}

func TestProductRepository_LowStockAndStats(t *testing.T) {
	db := setupProductTestDB(t)
	defer db.Close()
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	repo := NewProductRepository(db, log)
	ctx := context.Background()

	// Product with low stock
	p1 := domain.NewProduct("Low Stock Item", "retail", "Brand", "pcs", 100, 200, 10, 50)
	p1.ProductCode = "PRD-000001"
	p1.CurrentStock = 3 // below minimum of 10
	repo.CreateProduct(ctx, p1)

	// Product with healthy stock
	p2 := domain.NewProduct("Healthy Stock", "retail", "Brand", "pcs", 80, 150, 5, 30)
	p2.ProductCode = "PRD-000002"
	p2.CurrentStock = 20
	repo.CreateProduct(ctx, p2)

	lowItems, err := repo.GetLowStockProducts(ctx)
	if err != nil {
		t.Fatalf("GetLowStockProducts: %v", err)
	}
	if len(lowItems) != 1 {
		t.Errorf("expected 1 low stock item, got %d", len(lowItems))
	}
	if len(lowItems) > 0 && lowItems[0].ProductName != "Low Stock Item" {
		t.Errorf("expected 'Low Stock Item', got '%s'", lowItems[0].ProductName)
	}

	stats, err := repo.GetInventoryStats(ctx)
	if err != nil {
		t.Fatalf("GetInventoryStats: %v", err)
	}
	if stats.TotalProducts != 2 {
		t.Errorf("expected 2 total products, got %d", stats.TotalProducts)
	}
	if stats.ActiveProducts != 2 {
		t.Errorf("expected 2 active, got %d", stats.ActiveProducts)
	}
	if stats.LowStockCount != 1 {
		t.Errorf("expected 1 low stock, got %d", stats.LowStockCount)
	}
	// Value = 3*100 + 20*80 = 300+1600 = 1900
	if stats.TotalValue != 1900 {
		t.Errorf("expected total value 1900, got %f", stats.TotalValue)
	}
}

func TestProductRepository_NextCodes(t *testing.T) {
	db := setupProductTestDB(t)
	defer db.Close()
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	repo := NewProductRepository(db, log)
	ctx := context.Background()

	c1, _ := repo.NextProductCode(ctx)
	c2, _ := repo.NextProductCode(ctx)
	if c1 != "PRD-000001" {
		t.Errorf("expected PRD-000001, got %s", c1)
	}
	if c2 != "PRD-000002" {
		t.Errorf("expected PRD-000002, got %s", c2)
	}

	p1, _ := repo.NextPurchaseNumber(ctx)
	p2, _ := repo.NextPurchaseNumber(ctx)
	if p1 != "PUR-000001" {
		t.Errorf("expected PUR-000001, got %s", p1)
	}
	if p2 != "PUR-000002" {
		t.Errorf("expected PUR-000002, got %s", p2)
	}
}
