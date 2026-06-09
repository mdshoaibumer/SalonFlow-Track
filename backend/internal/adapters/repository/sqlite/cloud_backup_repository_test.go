package sqlite

import (
	"context"
	"database/sql"
	"log/slog"
	"os"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/salonflow/salonflow-track/internal/core/domain"
	"github.com/salonflow/salonflow-track/pkg/uid"
)

func setupCloudBackupTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(`
		CREATE TABLE cloud_backup_config (
			id TEXT PRIMARY KEY,
			provider TEXT NOT NULL DEFAULT 'none',
			bucket_name TEXT NOT NULL DEFAULT '',
			region TEXT NOT NULL DEFAULT '',
			access_key TEXT NOT NULL DEFAULT '',
			endpoint TEXT NOT NULL DEFAULT '',
			encrypt_backups INTEGER NOT NULL DEFAULT 1,
			auto_backup INTEGER NOT NULL DEFAULT 0,
			auto_backup_interval_hours INTEGER NOT NULL DEFAULT 24,
			max_versions INTEGER NOT NULL DEFAULT 10,
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL
		);
		CREATE TABLE cloud_backup_history (
			id TEXT PRIMARY KEY,
			provider TEXT NOT NULL,
			file_name TEXT NOT NULL,
			file_size INTEGER NOT NULL DEFAULT 0,
			remote_path TEXT NOT NULL DEFAULT '',
			status TEXT NOT NULL DEFAULT 'pending',
			is_encrypted INTEGER NOT NULL DEFAULT 0,
			error_message TEXT NOT NULL DEFAULT '',
			started_at TEXT NOT NULL DEFAULT '',
			completed_at TEXT NOT NULL DEFAULT '',
			created_at TEXT NOT NULL
		);
	`)
	if err != nil {
		t.Fatal(err)
	}
	return db
}

func TestCloudBackupRepository_Config(t *testing.T) {
	db := setupCloudBackupTestDB(t)
	defer db.Close()
	log := slog.New(slog.NewTextHandler(os.Stderr, nil))
	repo := NewCloudBackupRepository(db, log)
	ctx := context.Background()

	// Initially no config
	_, err := repo.GetConfig(ctx)
	if err == nil {
		t.Error("expected not found error")
	}

	// Save config
	cfg := domain.NewCloudBackupConfig()
	cfg.Provider = domain.CloudProviderS3
	cfg.BucketName = "salon-backups"
	cfg.Region = "ap-south-1"
	cfg.AccessKey = "AKIAIOSFODNN7EXAMPLE"

	err = repo.SaveConfig(ctx, cfg)
	if err != nil {
		t.Fatalf("SaveConfig failed: %v", err)
	}

	// Get config
	got, err := repo.GetConfig(ctx)
	if err != nil {
		t.Fatalf("GetConfig failed: %v", err)
	}
	if got.Provider != domain.CloudProviderS3 {
		t.Errorf("got provider %q, want %q", got.Provider, domain.CloudProviderS3)
	}
	if got.BucketName != "salon-backups" {
		t.Errorf("got bucket %q, want %q", got.BucketName, "salon-backups")
	}

	// Update config (upsert)
	cfg.Provider = domain.CloudProviderGDrive
	err = repo.SaveConfig(ctx, cfg)
	if err != nil {
		t.Fatalf("SaveConfig update failed: %v", err)
	}
	got, _ = repo.GetConfig(ctx)
	if got.Provider != domain.CloudProviderGDrive {
		t.Errorf("got provider %q, want %q", got.Provider, domain.CloudProviderGDrive)
	}
}

func TestCloudBackupRepository_History(t *testing.T) {
	db := setupCloudBackupTestDB(t)
	defer db.Close()
	log := slog.New(slog.NewTextHandler(os.Stderr, nil))
	repo := NewCloudBackupRepository(db, log)
	ctx := context.Background()

	h := domain.NewCloudBackupHistory(domain.CloudProviderS3, "backup-2024.db", 1024*1024, true)
	h.RemotePath = "backups/backup-2024.db"

	err := repo.CreateHistory(ctx, h)
	if err != nil {
		t.Fatalf("CreateHistory failed: %v", err)
	}

	list, total, err := repo.ListHistory(ctx, 10, 0)
	if err != nil {
		t.Fatalf("ListHistory failed: %v", err)
	}
	if total != 1 {
		t.Errorf("got total %d, want 1", total)
	}
	if len(list) != 1 {
		t.Errorf("got %d entries, want 1", len(list))
	}
	if list[0].FileName != "backup-2024.db" {
		t.Errorf("got filename %q, want %q", list[0].FileName, "backup-2024.db")
	}
}

func TestCloudBackupRepository_UpdateHistoryStatus(t *testing.T) {
	db := setupCloudBackupTestDB(t)
	defer db.Close()
	log := slog.New(slog.NewTextHandler(os.Stderr, nil))
	repo := NewCloudBackupRepository(db, log)
	ctx := context.Background()

	h := domain.NewCloudBackupHistory(domain.CloudProviderS3, "test.db", 512, false)
	_ = repo.CreateHistory(ctx, h)

	err := repo.UpdateHistoryStatus(ctx, h.ID, domain.CloudBackupCompleted, "")
	if err != nil {
		t.Fatalf("UpdateHistoryStatus failed: %v", err)
	}

	list, _, _ := repo.ListHistory(ctx, 10, 0)
	if list[0].Status != domain.CloudBackupCompleted {
		t.Errorf("got status %q, want %q", list[0].Status, domain.CloudBackupCompleted)
	}
}

func TestCloudBackupRepository_Stats(t *testing.T) {
	db := setupCloudBackupTestDB(t)
	defer db.Close()
	log := slog.New(slog.NewTextHandler(os.Stderr, nil))
	repo := NewCloudBackupRepository(db, log)
	ctx := context.Background()

	// Save config first
	cfg := domain.NewCloudBackupConfig()
	cfg.Provider = domain.CloudProviderS3
	cfg.AutoBackup = true
	_ = repo.SaveConfig(ctx, cfg)

	// Add a completed backup
	h := &domain.CloudBackupHistory{
		ID:          uid.New(),
		Provider:    domain.CloudProviderS3,
		FileName:    "backup.db",
		FileSize:    2048,
		RemotePath:  "backups/backup.db",
		Status:      domain.CloudBackupCompleted,
		IsEncrypted: true,
		CompletedAt: time.Now().UTC().Format(time.RFC3339),
		CreatedAt:   time.Now().UTC(),
	}
	_ = repo.CreateHistory(ctx, h)

	stats, err := repo.GetStats(ctx)
	if err != nil {
		t.Fatalf("GetStats failed: %v", err)
	}
	if stats.Provider != domain.CloudProviderS3 {
		t.Errorf("got provider %q, want %q", stats.Provider, domain.CloudProviderS3)
	}
	if stats.TotalBackups != 1 {
		t.Errorf("got total %d, want 1", stats.TotalBackups)
	}
	if stats.AutoEnabled != true {
		t.Error("expected auto_enabled true")
	}
}
