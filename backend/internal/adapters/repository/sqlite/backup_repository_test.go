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

func setupBackupTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(`
		CREATE TABLE backup_history (
			id TEXT PRIMARY KEY,
			backup_name TEXT NOT NULL,
			backup_type TEXT NOT NULL CHECK (backup_type IN ('manual','daily','before_update','before_restore')),
			backup_path TEXT NOT NULL,
			file_size INTEGER NOT NULL DEFAULT 0,
			checksum TEXT NOT NULL DEFAULT '',
			status TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending','completed','failed','corrupted','verified')),
			error_message TEXT NOT NULL DEFAULT '',
			created_at TEXT NOT NULL
		);
		CREATE TABLE restore_history (
			id TEXT PRIMARY KEY,
			backup_id TEXT NOT NULL REFERENCES backup_history(id),
			backup_name TEXT NOT NULL,
			restore_date TEXT NOT NULL,
			status TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending','completed','failed')),
			notes TEXT NOT NULL DEFAULT '',
			error_message TEXT NOT NULL DEFAULT '',
			created_at TEXT NOT NULL
		);
	`)
	if err != nil {
		t.Fatal(err)
	}
	return db
}

func TestBackupRepository_CreateAndGet(t *testing.T) {
	db := setupBackupTestDB(t)
	defer db.Close()
	log := slog.New(slog.NewTextHandler(os.Stderr, nil))
	repo := NewBackupRepository(db, log)
	ctx := context.Background()

	rec := &domain.BackupRecord{
		ID:         uuid.New(),
		BackupName: "manual_2025-01-15_10-30-00",
		BackupType: domain.BackupTypeManual,
		BackupPath: "/tmp/backups/2025/01/manual_2025-01-15_10-30-00.db",
		FileSize:   1024000,
		Checksum:   "abc123def456",
		Status:     domain.BackupStatusCompleted,
		CreatedAt:  time.Now().UTC(),
	}

	// Create
	err := repo.CreateBackupRecord(ctx, rec)
	if err != nil {
		t.Fatalf("CreateBackupRecord: %v", err)
	}

	// Get by ID
	got, err := repo.GetBackupByID(ctx, rec.ID)
	if err != nil {
		t.Fatalf("GetBackupByID: %v", err)
	}
	if got.BackupName != rec.BackupName {
		t.Errorf("BackupName = %q, want %q", got.BackupName, rec.BackupName)
	}
	if got.FileSize != rec.FileSize {
		t.Errorf("FileSize = %d, want %d", got.FileSize, rec.FileSize)
	}
	if got.Status != rec.Status {
		t.Errorf("Status = %q, want %q", got.Status, rec.Status)
	}
}

func TestBackupRepository_UpdateRecord(t *testing.T) {
	db := setupBackupTestDB(t)
	defer db.Close()
	log := slog.New(slog.NewTextHandler(os.Stderr, nil))
	repo := NewBackupRepository(db, log)
	ctx := context.Background()

	rec := &domain.BackupRecord{
		ID:         uuid.New(),
		BackupName: "test_backup",
		BackupType: domain.BackupTypeManual,
		BackupPath: "/tmp/test.db",
		Status:     domain.BackupStatusPending,
		CreatedAt:  time.Now().UTC(),
	}
	repo.CreateBackupRecord(ctx, rec)

	// Update
	rec.Status = domain.BackupStatusCompleted
	rec.FileSize = 2048
	rec.Checksum = "sha256hash"
	err := repo.UpdateBackupRecord(ctx, rec)
	if err != nil {
		t.Fatalf("UpdateBackupRecord: %v", err)
	}

	got, _ := repo.GetBackupByID(ctx, rec.ID)
	if got.Status != domain.BackupStatusCompleted {
		t.Errorf("Status = %q, want completed", got.Status)
	}
	if got.FileSize != 2048 {
		t.Errorf("FileSize = %d, want 2048", got.FileSize)
	}
}

func TestBackupRepository_ListBackups(t *testing.T) {
	db := setupBackupTestDB(t)
	defer db.Close()
	log := slog.New(slog.NewTextHandler(os.Stderr, nil))
	repo := NewBackupRepository(db, log)
	ctx := context.Background()

	// Insert 3 records
	for i := 0; i < 3; i++ {
		rec := &domain.BackupRecord{
			ID:         uuid.New(),
			BackupName: "backup_" + string(rune('a'+i)),
			BackupType: domain.BackupTypeManual,
			BackupPath: "/tmp/" + string(rune('a'+i)) + ".db",
			Status:     domain.BackupStatusCompleted,
			CreatedAt:  time.Now().UTC().Add(time.Duration(i) * time.Hour),
		}
		repo.CreateBackupRecord(ctx, rec)
	}

	records, total, err := repo.ListBackups(ctx, 10, 0)
	if err != nil {
		t.Fatalf("ListBackups: %v", err)
	}
	if total != 3 {
		t.Errorf("total = %d, want 3", total)
	}
	if len(records) != 3 {
		t.Errorf("len = %d, want 3", len(records))
	}

	// Test pagination
	records, _, err = repo.ListBackups(ctx, 2, 0)
	if err != nil {
		t.Fatalf("ListBackups paginated: %v", err)
	}
	if len(records) != 2 {
		t.Errorf("paginated len = %d, want 2", len(records))
	}
}

func TestBackupRepository_DeleteRecord(t *testing.T) {
	db := setupBackupTestDB(t)
	defer db.Close()
	log := slog.New(slog.NewTextHandler(os.Stderr, nil))
	repo := NewBackupRepository(db, log)
	ctx := context.Background()

	rec := &domain.BackupRecord{
		ID:         uuid.New(),
		BackupName: "to_delete",
		BackupType: domain.BackupTypeManual,
		BackupPath: "/tmp/del.db",
		Status:     domain.BackupStatusCompleted,
		CreatedAt:  time.Now().UTC(),
	}
	repo.CreateBackupRecord(ctx, rec)

	err := repo.DeleteBackupRecord(ctx, rec.ID)
	if err != nil {
		t.Fatalf("DeleteBackupRecord: %v", err)
	}

	_, err = repo.GetBackupByID(ctx, rec.ID)
	if err == nil {
		t.Error("expected not-found error after delete")
	}
}

func TestBackupRepository_Stats(t *testing.T) {
	db := setupBackupTestDB(t)
	defer db.Close()
	log := slog.New(slog.NewTextHandler(os.Stderr, nil))
	repo := NewBackupRepository(db, log)
	ctx := context.Background()

	// Empty stats
	stats, err := repo.GetBackupStats(ctx)
	if err != nil {
		t.Fatalf("GetBackupStats: %v", err)
	}
	if stats.TotalBackups != 0 {
		t.Errorf("TotalBackups = %d, want 0", stats.TotalBackups)
	}

	// Add a backup
	rec := &domain.BackupRecord{
		ID:         uuid.New(),
		BackupName: "stats_test",
		BackupType: domain.BackupTypeDaily,
		BackupPath: "/tmp/stats.db",
		FileSize:   5000,
		Status:     domain.BackupStatusCompleted,
		CreatedAt:  time.Now().UTC(),
	}
	repo.CreateBackupRecord(ctx, rec)

	stats, _ = repo.GetBackupStats(ctx)
	if stats.TotalBackups != 1 {
		t.Errorf("TotalBackups = %d, want 1", stats.TotalBackups)
	}
	if stats.LastBackupName != "stats_test" {
		t.Errorf("LastBackupName = %q, want stats_test", stats.LastBackupName)
	}
	if stats.LastBackupSize != 5000 {
		t.Errorf("LastBackupSize = %d, want 5000", stats.LastBackupSize)
	}
}

func TestBackupRepository_RestoreRecords(t *testing.T) {
	db := setupBackupTestDB(t)
	defer db.Close()
	log := slog.New(slog.NewTextHandler(os.Stderr, nil))
	repo := NewBackupRepository(db, log)
	ctx := context.Background()

	// Create a backup first (FK reference)
	backupRec := &domain.BackupRecord{
		ID:         uuid.New(),
		BackupName: "for_restore",
		BackupType: domain.BackupTypeManual,
		BackupPath: "/tmp/restore.db",
		Status:     domain.BackupStatusCompleted,
		CreatedAt:  time.Now().UTC(),
	}
	repo.CreateBackupRecord(ctx, backupRec)

	// Create restore record
	restoreRec := &domain.RestoreRecord{
		ID:          uuid.New(),
		BackupID:    backupRec.ID,
		BackupName:  backupRec.BackupName,
		RestoreDate: time.Now().UTC().Format(time.RFC3339),
		Status:      domain.RestoreStatusCompleted,
		Notes:       "test restore",
		CreatedAt:   time.Now().UTC(),
	}
	err := repo.CreateRestoreRecord(ctx, restoreRec)
	if err != nil {
		t.Fatalf("CreateRestoreRecord: %v", err)
	}

	// List restores
	restores, total, err := repo.ListRestores(ctx, 10, 0)
	if err != nil {
		t.Fatalf("ListRestores: %v", err)
	}
	if total != 1 {
		t.Errorf("total = %d, want 1", total)
	}
	if len(restores) != 1 {
		t.Errorf("len = %d, want 1", len(restores))
	}
	if restores[0].BackupName != "for_restore" {
		t.Errorf("BackupName = %q, want for_restore", restores[0].BackupName)
	}

	// Update restore record
	restoreRec.Status = domain.RestoreStatusFailed
	restoreRec.ErrorMessage = "disk full"
	err = repo.UpdateRestoreRecord(ctx, restoreRec)
	if err != nil {
		t.Fatalf("UpdateRestoreRecord: %v", err)
	}

	// Verify stats count
	stats, _ := repo.GetBackupStats(ctx)
	if stats.TotalRestores != 1 {
		t.Errorf("TotalRestores = %d, want 1", stats.TotalRestores)
	}
}
