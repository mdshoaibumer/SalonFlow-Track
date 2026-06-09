package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/salonflow/salonflow-track/internal/core/domain"
)

// CloudBackupRepository manages cloud backup data.
type CloudBackupRepository interface {
	GetConfig(ctx context.Context) (*domain.CloudBackupConfig, error)
	SaveConfig(ctx context.Context, cfg *domain.CloudBackupConfig) error

	CreateHistory(ctx context.Context, h *domain.CloudBackupHistory) error
	UpdateHistoryStatus(ctx context.Context, id uuid.UUID, status, errorMsg string) error
	ListHistory(ctx context.Context, limit, offset int) ([]domain.CloudBackupHistory, int, error)
	GetStats(ctx context.Context) (*domain.CloudBackupStats, error)
}

// CloudBackupEngine handles the actual upload/download to cloud storage.
type CloudBackupEngine interface {
	Upload(cfg *domain.CloudBackupConfig, localPath, remotePath string, encrypt bool) error
	Download(cfg *domain.CloudBackupConfig, remotePath, localPath string) error
	TestConnection(cfg *domain.CloudBackupConfig) error
}
