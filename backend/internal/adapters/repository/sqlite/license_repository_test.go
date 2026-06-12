package sqlite

import (
	"context"
	"database/sql"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"github.com/salonflow/salonflow-track/internal/core/domain"
)

func setupLicenseTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(`
		CREATE TABLE licenses (
			id TEXT PRIMARY KEY,
			license_key TEXT NOT NULL UNIQUE,
			customer_name TEXT NOT NULL,
			salon_name TEXT NOT NULL,
			device_id TEXT NOT NULL DEFAULT '',
			issued_date TEXT NOT NULL,
			expiry_date TEXT NOT NULL,
			grace_until TEXT NOT NULL DEFAULT '',
			status TEXT NOT NULL DEFAULT 'active' CHECK (status IN ('active','grace_period','expired','suspended')),
			signature TEXT NOT NULL DEFAULT '',
			last_validation TEXT NOT NULL DEFAULT '',
			last_verified_at TEXT NOT NULL DEFAULT '',
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL
		);
		CREATE TABLE license_events (
			id TEXT PRIMARY KEY,
			license_id TEXT NOT NULL REFERENCES licenses(id),
			event_type TEXT NOT NULL CHECK (event_type IN ('activated','renewed','expired','validated','suspended','grace_started','restricted','imported')),
			event_date TEXT NOT NULL,
			notes TEXT NOT NULL DEFAULT '',
			created_at TEXT NOT NULL
		);
	`)
	if err != nil {
		t.Fatal(err)
	}
	return db
}

func TestLicenseRepository_CreateAndGet(t *testing.T) {
	db := setupLicenseTestDB(t)
	defer db.Close()
	log := slog.New(slog.NewTextHandler(os.Stderr, nil))
	repo := NewLicenseRepository(db, log)
	ctx := context.Background()

	lic := &domain.License{
		ID:             uuid.New(),
		LicenseKey:     "SALONFLOW-ABCD-EFGH-IJKL",
		CustomerName:   "Test Customer",
		SalonName:      "Test Salon",
		DeviceID:       "device123",
		IssuedDate:     "2025-01-01",
		ExpiryDate:     "2025-02-01",
		Status:         domain.LicenseStatusActive,
		Signature:      "sig123",
		LastValidation: time.Now().Format(time.RFC3339),
		CreatedAt:      time.Now().UTC(),
		UpdatedAt:      time.Now().UTC(),
	}

	err := repo.CreateLicense(ctx, lic)
	if err != nil {
		t.Fatalf("CreateLicense: %v", err)
	}

	// Get by key
	got, err := repo.GetLicenseByKey(ctx, lic.LicenseKey)
	if err != nil {
		t.Fatalf("GetLicenseByKey: %v", err)
	}
	if got.SalonName != "Test Salon" {
		t.Errorf("SalonName = %q, want Test Salon", got.SalonName)
	}

	// Get active
	active, err := repo.GetActiveLicense(ctx)
	if err != nil {
		t.Fatalf("GetActiveLicense: %v", err)
	}
	if active.LicenseKey != lic.LicenseKey {
		t.Errorf("LicenseKey = %q, want %q", active.LicenseKey, lic.LicenseKey)
	}
}

func TestLicenseRepository_Update(t *testing.T) {
	db := setupLicenseTestDB(t)
	defer db.Close()
	log := slog.New(slog.NewTextHandler(os.Stderr, nil))
	repo := NewLicenseRepository(db, log)
	ctx := context.Background()

	lic := &domain.License{
		ID:             uuid.New(),
		LicenseKey:     "SALONFLOW-1111-2222-3333",
		CustomerName:   "Update Test",
		SalonName:      "Salon X",
		DeviceID:       "dev456",
		IssuedDate:     "2025-01-01",
		ExpiryDate:     "2025-02-01",
		Status:         domain.LicenseStatusActive,
		Signature:      "sig_orig",
		LastValidation: time.Now().Format(time.RFC3339),
		CreatedAt:      time.Now().UTC(),
		UpdatedAt:      time.Now().UTC(),
	}
	repo.CreateLicense(ctx, lic)

	// Update
	lic.Status = domain.LicenseStatusGracePeriod
	lic.ExpiryDate = "2025-03-01"
	lic.Signature = "sig_new"
	err := repo.UpdateLicense(ctx, lic)
	if err != nil {
		t.Fatalf("UpdateLicense: %v", err)
	}

	got, _ := repo.GetActiveLicense(ctx)
	if got.Status != domain.LicenseStatusGracePeriod {
		t.Errorf("Status = %q, want grace_period", got.Status)
	}
	if got.ExpiryDate != "2025-03-01" {
		t.Errorf("ExpiryDate = %q, want 2025-03-01", got.ExpiryDate)
	}
}

func TestLicenseRepository_Events(t *testing.T) {
	db := setupLicenseTestDB(t)
	defer db.Close()
	log := slog.New(slog.NewTextHandler(os.Stderr, nil))
	repo := NewLicenseRepository(db, log)
	ctx := context.Background()

	licID := uuid.New()
	lic := &domain.License{
		ID:             licID,
		LicenseKey:     "SALONFLOW-AAAA-BBBB-CCCC",
		CustomerName:   "Events Test",
		SalonName:      "Salon Y",
		DeviceID:       "dev789",
		IssuedDate:     "2025-01-01",
		ExpiryDate:     "2025-02-01",
		Status:         domain.LicenseStatusActive,
		Signature:      "sig",
		LastValidation: time.Now().Format(time.RFC3339),
		CreatedAt:      time.Now().UTC(),
		UpdatedAt:      time.Now().UTC(),
	}
	repo.CreateLicense(ctx, lic)

	// Create events
	ev1 := domain.NewLicenseEvent(licID, domain.LicenseEventActivated, "First activation")
	err := repo.CreateEvent(ctx, ev1)
	if err != nil {
		t.Fatalf("CreateEvent: %v", err)
	}

	ev2 := domain.NewLicenseEvent(licID, domain.LicenseEventValidated, "")
	repo.CreateEvent(ctx, ev2)

	// List events
	events, total, err := repo.ListEvents(ctx, licID, 10, 0)
	if err != nil {
		t.Fatalf("ListEvents: %v", err)
	}
	if total != 2 {
		t.Errorf("total = %d, want 2", total)
	}
	if len(events) != 2 {
		t.Errorf("len = %d, want 2", len(events))
	}
}

func TestLicenseRepository_NotFound(t *testing.T) {
	db := setupLicenseTestDB(t)
	defer db.Close()
	log := slog.New(slog.NewTextHandler(os.Stderr, nil))
	repo := NewLicenseRepository(db, log)
	ctx := context.Background()

	_, err := repo.GetLicenseByKey(ctx, "SALONFLOW-NOPE-NOPE-NOPE")
	if err == nil {
		t.Error("expected error for non-existent key")
	}

	_, err = repo.GetActiveLicense(ctx)
	if err == nil {
		t.Error("expected error for no license")
	}
}
