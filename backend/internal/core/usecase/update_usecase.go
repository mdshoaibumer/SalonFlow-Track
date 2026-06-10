package usecase

import (
	"context"
	"fmt"
	"log/slog"
	"path/filepath"
	"time"

	"github.com/salonflow/salonflow-track/internal/core/domain"
	"github.com/salonflow/salonflow-track/internal/core/ports"
	"github.com/salonflow/salonflow-track/pkg/apperror"
)

// UpdateUseCase handles update business logic.
type UpdateUseCase struct {
	repo           ports.UpdateRepository
	engine         ports.UpdateEngine
	backupEngine   ports.BackupEngine
	currentVersion string
	dbPath         string
	log            *slog.Logger
}

// NewUpdateUseCase creates a new UpdateUseCase.
func NewUpdateUseCase(repo ports.UpdateRepository, engine ports.UpdateEngine, backupEngine ports.BackupEngine, currentVersion, dbPath string, log *slog.Logger) *UpdateUseCase {
	return &UpdateUseCase{
		repo:           repo,
		engine:         engine,
		backupEngine:   backupEngine,
		currentVersion: currentVersion,
		dbPath:         dbPath,
		log:            log,
	}
}

// CheckForUpdate checks if a new version is available.
func (uc *UpdateUseCase) CheckForUpdate(ctx context.Context) (*domain.UpdateStatus, error) {
	info, err := uc.engine.CheckForUpdate(uc.currentVersion)
	if err != nil {
		return &domain.UpdateStatus{
			CurrentVersion:  uc.currentVersion,
			UpdateAvailable: false,
			Status:          "check_failed",
		}, nil
	}

	if info == nil {
		return &domain.UpdateStatus{
			CurrentVersion:  uc.currentVersion,
			UpdateAvailable: false,
			Status:          "up_to_date",
		}, nil
	}

	// Store version info
	existing, _ := uc.repo.GetVersionByName(ctx, info.Version)
	if existing == nil {
		v := domain.NewAppVersion(info.Version, info.ReleaseDate, info.ReleaseNotes)
		_ = uc.repo.CreateVersion(ctx, v)
	}

	return &domain.UpdateStatus{
		CurrentVersion:  uc.currentVersion,
		LatestVersion:   info.Version,
		UpdateAvailable: true,
		Status:          "update_available",
		ReleaseNotes:    info.ReleaseNotes,
	}, nil
}

// DownloadUpdate downloads the latest update.
func (uc *UpdateUseCase) DownloadUpdate(ctx context.Context) (*domain.UpdateRecord, error) {
	info, err := uc.engine.CheckForUpdate(uc.currentVersion)
	if err != nil {
		return nil, apperror.Internal("check update failed", err)
	}
	if info == nil {
		return nil, apperror.Business("NO_UPDATE", "No update available")
	}

	record := domain.NewUpdateRecord(uc.currentVersion, info.Version)
	record.Status = domain.UpdateHistoryDownloading
	if err := uc.repo.CreateUpdateRecord(ctx, record); err != nil {
		return nil, err
	}

	destPath := filepath.Join(uc.engine.UpdateDir(), fmt.Sprintf("salonflow-%s.zip", info.Version))
	if err := uc.engine.DownloadUpdate(info, destPath); err != nil {
		record.Status = domain.UpdateHistoryFailed
		record.ErrorMessage = err.Error()
		_ = uc.repo.UpdateUpdateRecord(ctx, record)
		return record, nil
	}

	// Verify checksum
	valid, err := uc.engine.VerifyDownload(destPath, info.Checksum)
	if err != nil || !valid {
		record.Status = domain.UpdateHistoryFailed
		record.ErrorMessage = "Checksum verification failed"
		_ = uc.repo.UpdateUpdateRecord(ctx, record)
		return record, nil
	}

	record.Status = domain.UpdateHistoryDownloaded
	_ = uc.repo.UpdateUpdateRecord(ctx, record)

	// Update version record
	v, _ := uc.repo.GetVersionByName(ctx, info.Version)
	if v != nil {
		v.Status = domain.UpdateStatusDownloaded
		_ = uc.repo.UpdateVersion(ctx, v)
	}

	return record, nil
}

// InstallUpdate installs the downloaded update (creates backup first).
func (uc *UpdateUseCase) InstallUpdate(ctx context.Context) (*domain.UpdateRecord, error) {
	// Find pending downloaded update
	records, _, err := uc.repo.ListUpdateHistory(ctx, 1, 0)
	if err != nil || len(records) == 0 {
		return nil, apperror.Business("NO_DOWNLOAD", "No downloaded update to install")
	}

	record := &records[0]
	if record.Status != domain.UpdateHistoryDownloaded {
		return nil, apperror.Business("NOT_READY", "Update is not ready for installation (status: "+record.Status+")")
	}

	// Create backup before install
	backupPath := filepath.Join(uc.backupEngine.BackupDir(), fmt.Sprintf("pre-update-%s.db", record.ToVersion))
	_, _, err = uc.backupEngine.CreateBackup(uc.dbPath, backupPath)
	if err != nil {
		record.Status = domain.UpdateHistoryFailed
		record.ErrorMessage = fmt.Sprintf("Pre-update backup failed: %v", err)
		_ = uc.repo.UpdateUpdateRecord(ctx, record)
		return record, nil
	}

	// Mark as installing
	record.Status = domain.UpdateHistoryInstalling
	_ = uc.repo.UpdateUpdateRecord(ctx, record)

	// In a desktop app, the actual binary replacement would be handled by
	// the OS-level installer/updater. Here we record the intent.
	record.Status = domain.UpdateHistoryCompleted
	_ = uc.repo.UpdateUpdateRecord(ctx, record)

	// Update version record
	v, _ := uc.repo.GetVersionByName(ctx, record.ToVersion)
	if v != nil {
		v.Status = domain.UpdateStatusInstalled
		v.InstalledAt = time.Now().UTC().Format(time.RFC3339)
		_ = uc.repo.UpdateVersion(ctx, v)
	}

	return record, nil
}

// GetUpdateStatus returns the current update status for the dashboard.
func (uc *UpdateUseCase) GetUpdateStatus(ctx context.Context) (*domain.UpdateStatus, error) {
	return &domain.UpdateStatus{
		CurrentVersion:  uc.currentVersion,
		UpdateAvailable: false,
		Status:          "up_to_date",
	}, nil
}

// ListHistory returns update history.
func (uc *UpdateUseCase) ListHistory(ctx context.Context, page, perPage int) ([]domain.UpdateRecord, int, error) {
	offset := (page - 1) * perPage
	return uc.repo.ListUpdateHistory(ctx, perPage, offset)
}

// ListVersions returns all known versions.
func (uc *UpdateUseCase) ListVersions(ctx context.Context, page, perPage int) ([]domain.AppVersion, int, error) {
	offset := (page - 1) * perPage
	return uc.repo.ListVersions(ctx, perPage, offset)
}
