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

func setupUpdateTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(`
		CREATE TABLE app_versions (
			id TEXT PRIMARY KEY,
			version TEXT NOT NULL UNIQUE,
			release_date TEXT NOT NULL,
			release_notes TEXT NOT NULL DEFAULT '',
			installed_at TEXT NOT NULL DEFAULT '',
			status TEXT NOT NULL DEFAULT 'available' CHECK (status IN ('available','downloading','downloaded','installing','installed','failed','rolled_back')),
			created_at TEXT NOT NULL
		);
		CREATE TABLE update_history (
			id TEXT PRIMARY KEY,
			from_version TEXT NOT NULL,
			to_version TEXT NOT NULL,
			update_date TEXT NOT NULL,
			status TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending','downloading','downloaded','installing','completed','failed','rolled_back')),
			error_message TEXT NOT NULL DEFAULT '',
			created_at TEXT NOT NULL
		);
	`)
	if err != nil {
		t.Fatal(err)
	}
	return db
}

func TestUpdateRepository_Versions(t *testing.T) {
	db := setupUpdateTestDB(t)
	defer db.Close()
	log := slog.New(slog.NewTextHandler(os.Stderr, nil))
	repo := NewUpdateRepository(db, log)
	ctx := context.Background()

	v := domain.NewAppVersion("1.1.0", "2025-06-01", "Bug fixes and improvements")
	err := repo.CreateVersion(ctx, v)
	if err != nil {
		t.Fatalf("CreateVersion: %v", err)
	}

	got, err := repo.GetVersionByName(ctx, "1.1.0")
	if err != nil {
		t.Fatalf("GetVersionByName: %v", err)
	}
	if got.ReleaseNotes != "Bug fixes and improvements" {
		t.Errorf("ReleaseNotes = %q", got.ReleaseNotes)
	}

	// Update to installed
	got.Status = domain.UpdateStatusInstalled
	got.InstalledAt = "2025-06-02T10:00:00Z"
	err = repo.UpdateVersion(ctx, got)
	if err != nil {
		t.Fatalf("UpdateVersion: %v", err)
	}

	installed, err := repo.GetInstalledVersion(ctx)
	if err != nil {
		t.Fatalf("GetInstalledVersion: %v", err)
	}
	if installed.Version != "1.1.0" {
		t.Errorf("installed version = %q, want 1.1.0", installed.Version)
	}
}

func TestUpdateRepository_ListVersions(t *testing.T) {
	db := setupUpdateTestDB(t)
	defer db.Close()
	log := slog.New(slog.NewTextHandler(os.Stderr, nil))
	repo := NewUpdateRepository(db, log)
	ctx := context.Background()

	repo.CreateVersion(ctx, domain.NewAppVersion("1.0.0", "2025-01-01", "Initial"))
	repo.CreateVersion(ctx, domain.NewAppVersion("1.1.0", "2025-02-01", "Features"))
	repo.CreateVersion(ctx, domain.NewAppVersion("1.2.0", "2025-03-01", "More features"))

	versions, total, err := repo.ListVersions(ctx, 10, 0)
	if err != nil {
		t.Fatalf("ListVersions: %v", err)
	}
	if total != 3 {
		t.Errorf("total = %d, want 3", total)
	}
	if len(versions) != 3 {
		t.Errorf("len = %d, want 3", len(versions))
	}
}

func TestUpdateRepository_UpdateHistory(t *testing.T) {
	db := setupUpdateTestDB(t)
	defer db.Close()
	log := slog.New(slog.NewTextHandler(os.Stderr, nil))
	repo := NewUpdateRepository(db, log)
	ctx := context.Background()

	rec := domain.NewUpdateRecord("1.0.0", "1.1.0")
	err := repo.CreateUpdateRecord(ctx, rec)
	if err != nil {
		t.Fatalf("CreateUpdateRecord: %v", err)
	}

	// Update status
	rec.Status = domain.UpdateHistoryCompleted
	err = repo.UpdateUpdateRecord(ctx, rec)
	if err != nil {
		t.Fatalf("UpdateUpdateRecord: %v", err)
	}

	records, total, err := repo.ListUpdateHistory(ctx, 10, 0)
	if err != nil {
		t.Fatalf("ListUpdateHistory: %v", err)
	}
	if total != 1 {
		t.Errorf("total = %d, want 1", total)
	}
	if records[0].Status != domain.UpdateHistoryCompleted {
		t.Errorf("status = %q, want completed", records[0].Status)
	}
}

func TestUpdateRepository_FailedUpdate(t *testing.T) {
	db := setupUpdateTestDB(t)
	defer db.Close()
	log := slog.New(slog.NewTextHandler(os.Stderr, nil))
	repo := NewUpdateRepository(db, log)
	ctx := context.Background()

	rec := domain.NewUpdateRecord("1.0.0", "1.1.0")
	repo.CreateUpdateRecord(ctx, rec)

	rec.Status = domain.UpdateHistoryFailed
	rec.ErrorMessage = "Download interrupted"
	err := repo.UpdateUpdateRecord(ctx, rec)
	if err != nil {
		t.Fatalf("UpdateUpdateRecord: %v", err)
	}

	records, _, _ := repo.ListUpdateHistory(ctx, 10, 0)
	if records[0].ErrorMessage != "Download interrupted" {
		t.Errorf("ErrorMessage = %q", records[0].ErrorMessage)
	}
}
