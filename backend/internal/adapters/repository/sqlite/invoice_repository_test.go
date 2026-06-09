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
	"github.com/salonflow/salonflow-track/pkg/uid"
)

func setupInvoiceTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.Exec(`
		CREATE TABLE invoices (
			id TEXT PRIMARY KEY,
			invoice_number TEXT NOT NULL UNIQUE,
			customer_id TEXT NOT NULL,
			staff_id TEXT NOT NULL,
			subtotal REAL NOT NULL DEFAULT 0,
			discount REAL NOT NULL DEFAULT 0,
			tax REAL NOT NULL DEFAULT 0,
			grand_total REAL NOT NULL DEFAULT 0,
			payment_status TEXT NOT NULL DEFAULT 'pending',
			payment_method TEXT NOT NULL DEFAULT '',
			notes TEXT NOT NULL DEFAULT '',
			invoice_date TEXT NOT NULL,
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL
		);
		CREATE TABLE invoice_items (
			id TEXT PRIMARY KEY,
			invoice_id TEXT NOT NULL,
			service_id TEXT NOT NULL,
			service_name_snapshot TEXT NOT NULL,
			quantity INTEGER NOT NULL DEFAULT 1,
			unit_price REAL NOT NULL,
			discount REAL NOT NULL DEFAULT 0,
			line_total REAL NOT NULL,
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL
		);
		CREATE TABLE payments (
			id TEXT PRIMARY KEY,
			invoice_id TEXT NOT NULL,
			amount REAL NOT NULL,
			payment_method TEXT NOT NULL,
			reference_number TEXT NOT NULL DEFAULT '',
			payment_date TEXT NOT NULL,
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL
		);
	`)
	if err != nil {
		t.Fatal(err)
	}
	return db
}

func TestInvoiceRepository_Create(t *testing.T) {
	db := setupInvoiceTestDB(t)
	defer db.Close()
	repo := NewInvoiceRepository(db, slog.New(slog.NewTextHandler(os.Stdout, nil)))

	customerID := uid.New()
	staffID := uid.New()
	inv := domain.NewInvoice(customerID, staffID, "INV-2026-000001")
	item := domain.NewInvoiceItem(inv.ID, uid.New(), "Hair Cut", 1, 300, 0)
	inv.AddItem(item)

	err := repo.Create(context.Background(), inv)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestInvoiceRepository_GetByID(t *testing.T) {
	db := setupInvoiceTestDB(t)
	defer db.Close()
	repo := NewInvoiceRepository(db, slog.New(slog.NewTextHandler(os.Stdout, nil)))

	customerID := uid.New()
	staffID := uid.New()
	inv := domain.NewInvoice(customerID, staffID, "INV-2026-000001")
	item := domain.NewInvoiceItem(inv.ID, uid.New(), "Hair Cut", 1, 300, 0)
	inv.AddItem(item)
	_ = repo.Create(context.Background(), inv)

	found, err := repo.GetByID(context.Background(), inv.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if found.InvoiceNumber != "INV-2026-000001" {
		t.Errorf("expected INV-2026-000001, got %q", found.InvoiceNumber)
	}
	if len(found.Items) != 1 {
		t.Errorf("expected 1 item, got %d", len(found.Items))
	}
}

func TestInvoiceRepository_List(t *testing.T) {
	db := setupInvoiceTestDB(t)
	defer db.Close()
	repo := NewInvoiceRepository(db, slog.New(slog.NewTextHandler(os.Stdout, nil)))

	for i := 1; i <= 3; i++ {
		inv := domain.NewInvoice(uid.New(), uid.New(), domain.GenerateInvoiceNumber(2026, i))
		item := domain.NewInvoiceItem(inv.ID, uid.New(), "Service", 1, 500, 0)
		inv.AddItem(item)
		_ = repo.Create(context.Background(), inv)
	}

	invoices, total, err := repo.List(context.Background(), ports.InvoiceFilter{Limit: 10})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if total != 3 {
		t.Errorf("expected 3, got %d", total)
	}
	if len(invoices) != 3 {
		t.Errorf("expected 3 invoices, got %d", len(invoices))
	}
}

func TestInvoiceRepository_GetNextSequence(t *testing.T) {
	db := setupInvoiceTestDB(t)
	defer db.Close()
	repo := NewInvoiceRepository(db, slog.New(slog.NewTextHandler(os.Stdout, nil)))

	seq, err := repo.GetNextSequence(context.Background(), 2026)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if seq != 1 {
		t.Errorf("expected 1, got %d", seq)
	}

	inv := domain.NewInvoice(uid.New(), uid.New(), "INV-2026-000001")
	item := domain.NewInvoiceItem(inv.ID, uid.New(), "Test", 1, 100, 0)
	inv.AddItem(item)
	_ = repo.Create(context.Background(), inv)

	seq, err = repo.GetNextSequence(context.Background(), 2026)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if seq != 2 {
		t.Errorf("expected 2, got %d", seq)
	}
}

func TestInvoiceRepository_RecordPayment(t *testing.T) {
	db := setupInvoiceTestDB(t)
	defer db.Close()
	repo := NewInvoiceRepository(db, slog.New(slog.NewTextHandler(os.Stdout, nil)))

	invoiceID := uid.New()
	inv := domain.NewInvoice(uid.New(), uid.New(), "INV-2026-000001")
	inv.ID = invoiceID
	item := domain.NewInvoiceItem(inv.ID, uid.New(), "Test", 1, 500, 0)
	inv.AddItem(item)
	_ = repo.Create(context.Background(), inv)

	payment := domain.NewPayment(invoiceID, 500, domain.PaymentMethodCash, "")
	err := repo.RecordPayment(context.Background(), payment)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	payments, err := repo.GetPayments(context.Background(), invoiceID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(payments) != 1 {
		t.Errorf("expected 1 payment, got %d", len(payments))
	}
}

func TestInvoiceRepository_UpdateStatus(t *testing.T) {
	db := setupInvoiceTestDB(t)
	defer db.Close()
	repo := NewInvoiceRepository(db, slog.New(slog.NewTextHandler(os.Stdout, nil)))

	inv := domain.NewInvoice(uid.New(), uid.New(), "INV-2026-000001")
	item := domain.NewInvoiceItem(inv.ID, uid.New(), "Test", 1, 500, 0)
	inv.AddItem(item)
	_ = repo.Create(context.Background(), inv)

	err := repo.UpdateStatus(context.Background(), inv.ID, domain.PaymentStatusPaid, domain.PaymentMethodCash)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	found, _ := repo.GetByID(context.Background(), inv.ID)
	if found.PaymentStatus != domain.PaymentStatusPaid {
		t.Errorf("expected paid, got %q", found.PaymentStatus)
	}
}

func TestInvoiceRepository_NotFound(t *testing.T) {
	db := setupInvoiceTestDB(t)
	defer db.Close()
	repo := NewInvoiceRepository(db, slog.New(slog.NewTextHandler(os.Stdout, nil)))

	_, err := repo.GetByID(context.Background(), uuid.New())
	if err == nil {
		t.Error("expected not found error")
	}
}
