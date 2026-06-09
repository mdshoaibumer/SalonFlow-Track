package ports

import (
	"context"

	"github.com/salonflow/salonflow-track/internal/core/domain"
)

// UpdateRepository manages update history in the database.
type UpdateRepository interface {
	CreateVersion(ctx context.Context, version *domain.AppVersion) error
	UpdateVersion(ctx context.Context, version *domain.AppVersion) error
	GetVersionByName(ctx context.Context, version string) (*domain.AppVersion, error)
	GetInstalledVersion(ctx context.Context) (*domain.AppVersion, error)
	ListVersions(ctx context.Context, limit, offset int) ([]domain.AppVersion, int, error)

	CreateUpdateRecord(ctx context.Context, record *domain.UpdateRecord) error
	UpdateUpdateRecord(ctx context.Context, record *domain.UpdateRecord) error
	ListUpdateHistory(ctx context.Context, limit, offset int) ([]domain.UpdateRecord, int, error)
}

// UpdateEngine handles checking for updates, downloading, verifying, and installing.
type UpdateEngine interface {
	CheckForUpdate(currentVersion string) (*domain.VersionInfo, error)
	DownloadUpdate(info *domain.VersionInfo, destPath string) error
	VerifyDownload(filePath, expectedChecksum string) (bool, error)
	UpdateDir() string
}
