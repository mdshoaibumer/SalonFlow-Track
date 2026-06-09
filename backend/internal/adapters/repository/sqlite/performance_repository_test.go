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

func setupPerformanceTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.Exec(`
		CREATE TABLE staff (
			id TEXT PRIMARY KEY,
			staff_code TEXT NOT NULL,
			full_name TEXT NOT NULL,
			phone TEXT NOT NULL,
			email TEXT NOT NULL DEFAULT '',
			gender TEXT NOT NULL DEFAULT 'male',
			designation TEXT NOT NULL DEFAULT 'stylist',
			joining_date TEXT NOT NULL DEFAULT '',
			base_salary REAL NOT NULL DEFAULT 0,
			commission_percentage REAL NOT NULL DEFAULT 0,
			status TEXT NOT NULL DEFAULT 'active',
			created_at TEXT NOT NULL DEFAULT '',
			updated_at TEXT NOT NULL DEFAULT ''
		);
		CREATE TABLE staff_performance_daily (
			id TEXT PRIMARY KEY,
			staff_id TEXT NOT NULL REFERENCES staff(id),
			business_date TEXT NOT NULL,
			invoice_count INTEGER NOT NULL DEFAULT 0,
			customer_count INTEGER NOT NULL DEFAULT 0,
			service_count INTEGER NOT NULL DEFAULT 0,
			revenue REAL NOT NULL DEFAULT 0,
			commission_amount REAL NOT NULL DEFAULT 0,
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL
		);
		CREATE UNIQUE INDEX idx_staff_performance_daily_staff_date ON staff_performance_daily(staff_id, business_date);
	`)
	if err != nil {
		t.Fatal(err)
	}
	return db
}

func insertTestStaff(t *testing.T, db *sql.DB, id uuid.UUID, name string) {
	t.Helper()
	_, err := db.Exec("INSERT INTO staff (id, staff_code, full_name, phone, status) VALUES (?, ?, ?, ?, 'active')",
		id.String(), "STF-"+name, name, "9999"+name)
	if err != nil {
		t.Fatal(err)
	}
}

func TestPerformanceRepository_Upsert(t *testing.T) {
	db := setupPerformanceTestDB(t)
	defer db.Close()
	repo := NewPerformanceRepository(db, slog.New(slog.NewTextHandler(os.Stdout, nil)))

	staffID := uuid.New()
	insertTestStaff(t, db, staffID, "Nazim")

	// First insert
	perf := domain.NewStaffPerformanceDaily(staffID, "2026-06-09")
	perf.AddInvoice(2000, 3, 200)
	err := repo.Upsert(context.Background(), perf)
	if err != nil {
		t.Fatalf("first upsert failed: %v", err)
	}

	// Second upsert (should add to existing)
	perf2 := domain.NewStaffPerformanceDaily(staffID, "2026-06-09")
	perf2.AddInvoice(1500, 2, 150)
	err = repo.Upsert(context.Background(), perf2)
	if err != nil {
		t.Fatalf("second upsert failed: %v", err)
	}

	// Check data
	records, err := repo.GetByStaff(context.Background(), staffID, "2026-06-09", "2026-06-09")
	if err != nil {
		t.Fatalf("get by staff failed: %v", err)
	}
	if len(records) != 1 {
		t.Fatalf("expected 1 record, got %d", len(records))
	}
	if records[0].Revenue != 3500 {
		t.Errorf("expected revenue 3500, got %.2f", records[0].Revenue)
	}
	if records[0].InvoiceCount != 2 {
		t.Errorf("expected 2 invoices, got %d", records[0].InvoiceCount)
	}
	if records[0].ServiceCount != 5 {
		t.Errorf("expected 5 services, got %d", records[0].ServiceCount)
	}
}

func TestPerformanceRepository_GetDaily(t *testing.T) {
	db := setupPerformanceTestDB(t)
	defer db.Close()
	repo := NewPerformanceRepository(db, slog.New(slog.NewTextHandler(os.Stdout, nil)))

	staff1 := uuid.New()
	staff2 := uuid.New()
	insertTestStaff(t, db, staff1, "Nazim")
	insertTestStaff(t, db, staff2, "Ravi")

	p1 := domain.NewStaffPerformanceDaily(staff1, "2026-06-09")
	p1.AddInvoice(5000, 4, 500)
	_ = repo.Upsert(context.Background(), p1)

	p2 := domain.NewStaffPerformanceDaily(staff2, "2026-06-09")
	p2.AddInvoice(3000, 2, 300)
	_ = repo.Upsert(context.Background(), p2)

	summaries, err := repo.GetDaily(context.Background(), ports.PerformanceFilter{
		DateFrom: "2026-06-09",
		DateTo:   "2026-06-09",
	})
	if err != nil {
		t.Fatalf("get daily failed: %v", err)
	}
	if len(summaries) < 2 {
		t.Fatalf("expected at least 2 summaries, got %d", len(summaries))
	}
	// First should be highest revenue
	if summaries[0].Revenue != 5000 {
		t.Errorf("expected top revenue 5000, got %.2f", summaries[0].Revenue)
	}
	if summaries[0].Rank != 1 {
		t.Errorf("expected rank 1, got %d", summaries[0].Rank)
	}
}

func TestPerformanceRepository_GetTopPerformers(t *testing.T) {
	db := setupPerformanceTestDB(t)
	defer db.Close()
	repo := NewPerformanceRepository(db, slog.New(slog.NewTextHandler(os.Stdout, nil)))

	staff1 := uuid.New()
	insertTestStaff(t, db, staff1, "Nazim")

	p1 := domain.NewStaffPerformanceDaily(staff1, "2026-06-09")
	p1.AddInvoice(15000, 7, 1500)
	_ = repo.Upsert(context.Background(), p1)

	top, err := repo.GetTopPerformers(context.Background(), "2026-06-01", "2026-06-30", 5)
	if err != nil {
		t.Fatalf("get top performers failed: %v", err)
	}
	if len(top) != 1 {
		t.Fatalf("expected 1 top performer, got %d", len(top))
	}
	if top[0].Revenue != 15000 {
		t.Errorf("expected revenue 15000, got %.2f", top[0].Revenue)
	}
	// AvgBill = 15000/1 = 15000
	if top[0].AvgBill != 15000 {
		t.Errorf("expected avg bill 15000, got %.2f", top[0].AvgBill)
	}
}

func TestPerformanceRepository_GetRevenueTrend(t *testing.T) {
	db := setupPerformanceTestDB(t)
	defer db.Close()
	repo := NewPerformanceRepository(db, slog.New(slog.NewTextHandler(os.Stdout, nil)))

	staffID := uuid.New()
	insertTestStaff(t, db, staffID, "Nazim")

	for i, date := range []string{"2026-06-07", "2026-06-08", "2026-06-09"} {
		p := domain.NewStaffPerformanceDaily(staffID, date)
		p.AddInvoice(float64((i+1)*1000), i+1, float64((i+1)*100))
		_ = repo.Upsert(context.Background(), p)
	}

	trend, err := repo.GetRevenueTrend(context.Background(), "2026-06-07", "2026-06-09")
	if err != nil {
		t.Fatalf("get revenue trend failed: %v", err)
	}
	if len(trend) != 3 {
		t.Fatalf("expected 3 points, got %d", len(trend))
	}
	if trend[0].Revenue != 1000 {
		t.Errorf("expected first point revenue 1000, got %.2f", trend[0].Revenue)
	}
	if trend[2].Revenue != 3000 {
		t.Errorf("expected third point revenue 3000, got %.2f", trend[2].Revenue)
	}
}
