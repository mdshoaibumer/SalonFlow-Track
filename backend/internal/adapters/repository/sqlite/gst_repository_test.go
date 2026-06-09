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
)

func setupGSTTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(`
		CREATE TABLE gst_settings (
			id TEXT PRIMARY KEY,
			business_name TEXT NOT NULL DEFAULT '',
			gstin TEXT NOT NULL DEFAULT '',
			state TEXT NOT NULL DEFAULT '',
			address TEXT NOT NULL DEFAULT '',
			hsn_code TEXT NOT NULL DEFAULT '',
			cgst_rate REAL NOT NULL DEFAULT 9.0,
			sgst_rate REAL NOT NULL DEFAULT 9.0,
			igst_rate REAL NOT NULL DEFAULT 18.0,
			is_gst_enabled INTEGER NOT NULL DEFAULT 0,
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL
		);
		CREATE TABLE tax_rates (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			hsn_code TEXT NOT NULL DEFAULT '',
			cgst_rate REAL NOT NULL DEFAULT 9.0,
			sgst_rate REAL NOT NULL DEFAULT 9.0,
			igst_rate REAL NOT NULL DEFAULT 18.0,
			category TEXT NOT NULL DEFAULT '',
			is_active INTEGER NOT NULL DEFAULT 1,
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL
		);
		CREATE TABLE invoice_tax_lines (
			id TEXT PRIMARY KEY,
			invoice_id TEXT NOT NULL,
			item_id TEXT NOT NULL DEFAULT '',
			taxable_amount REAL NOT NULL DEFAULT 0,
			cgst_rate REAL NOT NULL DEFAULT 0,
			cgst_amount REAL NOT NULL DEFAULT 0,
			sgst_rate REAL NOT NULL DEFAULT 0,
			sgst_amount REAL NOT NULL DEFAULT 0,
			igst_rate REAL NOT NULL DEFAULT 0,
			igst_amount REAL NOT NULL DEFAULT 0,
			total_tax REAL NOT NULL DEFAULT 0,
			is_interstate INTEGER NOT NULL DEFAULT 0,
			hsn_code TEXT NOT NULL DEFAULT '',
			created_at TEXT NOT NULL
		);
	`)
	if err != nil {
		t.Fatal(err)
	}
	return db
}

func TestGSTRepository_Settings(t *testing.T) {
	db := setupGSTTestDB(t)
	defer db.Close()
	log := slog.New(slog.NewTextHandler(os.Stderr, nil))
	repo := NewGSTRepository(db, log)
	ctx := context.Background()

	// Initially not found
	_, err := repo.GetSettings(ctx)
	if err == nil {
		t.Fatal("expected not found error")
	}

	// Save settings
	s := domain.NewGSTSettings()
	s.BusinessName = "Glamour Salon"
	s.GSTIN = "27AABCU9603R1ZM"
	s.State = "Maharashtra"
	s.IsGSTEnabled = true

	err = repo.SaveSettings(ctx, s)
	if err != nil {
		t.Fatalf("SaveSettings: %v", err)
	}

	// Get settings
	got, err := repo.GetSettings(ctx)
	if err != nil {
		t.Fatalf("GetSettings: %v", err)
	}
	if got.BusinessName != "Glamour Salon" {
		t.Errorf("BusinessName = %q", got.BusinessName)
	}
	if got.GSTIN != "27AABCU9603R1ZM" {
		t.Errorf("GSTIN = %q", got.GSTIN)
	}
	if !got.IsGSTEnabled {
		t.Error("expected IsGSTEnabled = true")
	}

	// Update
	s.BusinessName = "Super Salon"
	err = repo.SaveSettings(ctx, s)
	if err != nil {
		t.Fatalf("SaveSettings update: %v", err)
	}
	got, _ = repo.GetSettings(ctx)
	if got.BusinessName != "Super Salon" {
		t.Errorf("updated BusinessName = %q", got.BusinessName)
	}
}

func TestGSTRepository_TaxRates(t *testing.T) {
	db := setupGSTTestDB(t)
	defer db.Close()
	log := slog.New(slog.NewTextHandler(os.Stderr, nil))
	repo := NewGSTRepository(db, log)
	ctx := context.Background()

	rate := domain.NewTaxRate("Hair Services", "9902", "hair", 9, 9, 18)
	err := repo.CreateTaxRate(ctx, rate)
	if err != nil {
		t.Fatalf("CreateTaxRate: %v", err)
	}

	rate2 := domain.NewTaxRate("Spa Services", "9902", "spa", 9, 9, 18)
	repo.CreateTaxRate(ctx, rate2)

	// List all
	rates, err := repo.ListTaxRates(ctx, "")
	if err != nil {
		t.Fatalf("ListTaxRates: %v", err)
	}
	if len(rates) != 2 {
		t.Errorf("len = %d, want 2", len(rates))
	}

	// List by category
	rates, _ = repo.ListTaxRates(ctx, "hair")
	if len(rates) != 1 {
		t.Errorf("hair rates = %d, want 1", len(rates))
	}

	// Get by ID
	got, err := repo.GetTaxRate(ctx, rate.ID)
	if err != nil {
		t.Fatalf("GetTaxRate: %v", err)
	}
	if got.Name != "Hair Services" {
		t.Errorf("Name = %q", got.Name)
	}

	// Get by category
	got, err = repo.GetTaxRateByCategory(ctx, "hair")
	if err != nil {
		t.Fatalf("GetTaxRateByCategory: %v", err)
	}
	if got.Category != "hair" {
		t.Errorf("Category = %q", got.Category)
	}

	// Update
	rate.Name = "Premium Hair"
	err = repo.UpdateTaxRate(ctx, rate)
	if err != nil {
		t.Fatalf("UpdateTaxRate: %v", err)
	}
	got, _ = repo.GetTaxRate(ctx, rate.ID)
	if got.Name != "Premium Hair" {
		t.Errorf("updated Name = %q", got.Name)
	}

	// Delete
	err = repo.DeleteTaxRate(ctx, rate.ID)
	if err != nil {
		t.Fatalf("DeleteTaxRate: %v", err)
	}
	_, err = repo.GetTaxRate(ctx, rate.ID)
	if err == nil {
		t.Error("expected not found after delete")
	}
}

func TestGSTRepository_TaxLines(t *testing.T) {
	db := setupGSTTestDB(t)
	defer db.Close()
	log := slog.New(slog.NewTextHandler(os.Stderr, nil))
	repo := NewGSTRepository(db, log)
	ctx := context.Background()

	invoiceID := uuid.New()
	lines := []domain.InvoiceTaxLine{
		*domain.NewInvoiceTaxLine(invoiceID, "item1", 1000, false, 9, 9, 18, "9902"),
		*domain.NewInvoiceTaxLine(invoiceID, "item2", 500, false, 9, 9, 18, "9902"),
	}

	err := repo.CreateTaxLines(ctx, lines)
	if err != nil {
		t.Fatalf("CreateTaxLines: %v", err)
	}

	got, err := repo.GetTaxLinesByInvoice(ctx, invoiceID)
	if err != nil {
		t.Fatalf("GetTaxLinesByInvoice: %v", err)
	}
	if len(got) != 2 {
		t.Errorf("len = %d, want 2", len(got))
	}

	// Verify calculation
	if got[0].CGSTAmount != 90 {
		t.Errorf("CGST = %.2f, want 90", got[0].CGSTAmount)
	}
	if got[0].SGSTAmount != 90 {
		t.Errorf("SGST = %.2f, want 90", got[0].SGSTAmount)
	}
	if got[0].TotalTax != 180 {
		t.Errorf("TotalTax = %.2f, want 180", got[0].TotalTax)
	}
}

func TestGSTRepository_Report(t *testing.T) {
	db := setupGSTTestDB(t)
	defer db.Close()
	log := slog.New(slog.NewTextHandler(os.Stderr, nil))
	repo := NewGSTRepository(db, log)
	ctx := context.Background()

	invoiceID := uuid.New()
	lines := []domain.InvoiceTaxLine{
		*domain.NewInvoiceTaxLine(invoiceID, "item1", 1000, false, 9, 9, 18, "9902"),
	}
	repo.CreateTaxLines(ctx, lines)

	report, err := repo.GetGSTReport(ctx, domain.GSTReportFilter{
		Period:    "daily",
		StartDate: "2020-01-01",
		EndDate:   "2030-12-31",
	})
	if err != nil {
		t.Fatalf("GetGSTReport: %v", err)
	}
	if report.TotalInvoices != 1 {
		t.Errorf("TotalInvoices = %d, want 1", report.TotalInvoices)
	}
	if report.TaxableAmount != 1000 {
		t.Errorf("TaxableAmount = %.2f, want 1000", report.TaxableAmount)
	}
	if report.TotalCGST != 90 {
		t.Errorf("TotalCGST = %.2f, want 90", report.TotalCGST)
	}
}
