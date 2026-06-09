package sqlite

import (
	"context"
	"database/sql"
	"log/slog"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/salonflow/salonflow-track/internal/core/domain"
)

func setupPrinterTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(`
		CREATE TABLE printer_settings (
			id TEXT PRIMARY KEY,
			default_printer TEXT NOT NULL DEFAULT '',
			paper_width TEXT NOT NULL DEFAULT '80mm',
			margin_top INTEGER NOT NULL DEFAULT 5,
			margin_bottom INTEGER NOT NULL DEFAULT 5,
			margin_left INTEGER NOT NULL DEFAULT 5,
			margin_right INTEGER NOT NULL DEFAULT 5,
			header_text TEXT NOT NULL DEFAULT '',
			footer_text TEXT NOT NULL DEFAULT 'Thank you for visiting!',
			show_logo INTEGER NOT NULL DEFAULT 0,
			show_qr INTEGER NOT NULL DEFAULT 0,
			upi_id TEXT NOT NULL DEFAULT '',
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL
		);
		CREATE TABLE print_jobs (
			id TEXT PRIMARY KEY,
			document_type TEXT NOT NULL,
			document_id TEXT NOT NULL DEFAULT '',
			printer_name TEXT NOT NULL DEFAULT '',
			paper_width TEXT NOT NULL DEFAULT '80mm',
			status TEXT NOT NULL DEFAULT 'queued',
			copies INTEGER NOT NULL DEFAULT 1,
			created_at TEXT NOT NULL
		);
	`)
	if err != nil {
		t.Fatal(err)
	}
	return db
}

func TestPrinterRepository_Settings(t *testing.T) {
	db := setupPrinterTestDB(t)
	defer db.Close()
	log := slog.New(slog.NewTextHandler(os.Stderr, nil))
	repo := NewPrinterRepository(db, log)
	ctx := context.Background()

	// Initially not found
	_, err := repo.GetSettings(ctx)
	if err == nil {
		t.Fatal("expected not found error")
	}

	// Save settings
	s := domain.NewPrinterSettings()
	s.DefaultPrinter = "ThermalPOS-80"
	s.PaperWidth = domain.PaperWidth80mm
	s.FooterText = "Visit again!"
	s.ShowQR = true

	err = repo.SaveSettings(ctx, s)
	if err != nil {
		t.Fatalf("SaveSettings: %v", err)
	}

	got, err := repo.GetSettings(ctx)
	if err != nil {
		t.Fatalf("GetSettings: %v", err)
	}
	if got.DefaultPrinter != "ThermalPOS-80" {
		t.Errorf("DefaultPrinter = %q", got.DefaultPrinter)
	}
	if got.PaperWidth != "80mm" {
		t.Errorf("PaperWidth = %q", got.PaperWidth)
	}
	if !got.ShowQR {
		t.Error("expected ShowQR = true")
	}

	// Update
	s.DefaultPrinter = "NewPrinter"
	repo.SaveSettings(ctx, s)
	got, _ = repo.GetSettings(ctx)
	if got.DefaultPrinter != "NewPrinter" {
		t.Errorf("updated DefaultPrinter = %q", got.DefaultPrinter)
	}
}

func TestPrinterRepository_PrintJobs(t *testing.T) {
	db := setupPrinterTestDB(t)
	defer db.Close()
	log := slog.New(slog.NewTextHandler(os.Stderr, nil))
	repo := NewPrinterRepository(db, log)
	ctx := context.Background()

	job := domain.NewPrintJob(domain.PrintDocInvoice, "INV-001", "ThermalPOS", "80mm", 1)
	err := repo.CreatePrintJob(ctx, job)
	if err != nil {
		t.Fatalf("CreatePrintJob: %v", err)
	}

	// Get
	got, err := repo.GetPrintJob(ctx, job.ID)
	if err != nil {
		t.Fatalf("GetPrintJob: %v", err)
	}
	if got.DocumentType != domain.PrintDocInvoice {
		t.Errorf("DocumentType = %q", got.DocumentType)
	}
	if got.Status != domain.PrintStatusQueued {
		t.Errorf("Status = %q", got.Status)
	}

	// Update status
	err = repo.UpdatePrintJobStatus(ctx, job.ID, domain.PrintStatusCompleted)
	if err != nil {
		t.Fatalf("UpdatePrintJobStatus: %v", err)
	}
	got, _ = repo.GetPrintJob(ctx, job.ID)
	if got.Status != domain.PrintStatusCompleted {
		t.Errorf("Status after update = %q", got.Status)
	}

	// Create more and list
	repo.CreatePrintJob(ctx, domain.NewPrintJob(domain.PrintDocReceipt, "INV-002", "ThermalPOS", "58mm", 2))
	jobs, total, err := repo.ListPrintJobs(ctx, 10, 0)
	if err != nil {
		t.Fatalf("ListPrintJobs: %v", err)
	}
	if total != 2 {
		t.Errorf("total = %d, want 2", total)
	}
	if len(jobs) != 2 {
		t.Errorf("len = %d, want 2", len(jobs))
	}
}
