package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/salonflow/salonflow-track/internal/core/domain"
)

// BackupRepository manages backup and restore history in the database.
type BackupRepository interface {
	CreateBackupRecord(ctx context.Context, record *domain.BackupRecord) error
	UpdateBackupRecord(ctx context.Context, record *domain.BackupRecord) error
	GetBackupByID(ctx context.Context, id uuid.UUID) (*domain.BackupRecord, error)
	ListBackups(ctx context.Context, limit, offset int) ([]domain.BackupRecord, int, error)
	DeleteBackupRecord(ctx context.Context, id uuid.UUID) error
	GetBackupStats(ctx context.Context) (*domain.BackupStats, error)

	CreateRestoreRecord(ctx context.Context, record *domain.RestoreRecord) error
	UpdateRestoreRecord(ctx context.Context, record *domain.RestoreRecord) error
	ListRestores(ctx context.Context, limit, offset int) ([]domain.RestoreRecord, int, error)
}

// BackupEngine handles the actual file-system backup and restore operations.
type BackupEngine interface {
	CreateBackup(dbPath, destPath string) (int64, string, error)
	VerifyBackup(backupPath, expectedChecksum string) *domain.BackupVerification
	RestoreBackup(backupPath, dbPath string) error
	BackupDir() string
}
