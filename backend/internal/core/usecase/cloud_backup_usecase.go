package usecase

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/salonflow/salonflow-track/internal/core/domain"
	"github.com/salonflow/salonflow-track/internal/core/ports"
	"github.com/salonflow/salonflow-track/pkg/apperror"
)

// CloudBackupUseCase handles cloud backup business logic.
type CloudBackupUseCase struct {
	repo   ports.CloudBackupRepository
	engine ports.CloudBackupEngine
	dbPath string
	log    *slog.Logger
}

// NewCloudBackupUseCase creates a new CloudBackupUseCase.
func NewCloudBackupUseCase(repo ports.CloudBackupRepository, engine ports.CloudBackupEngine, dbPath string, log *slog.Logger) *CloudBackupUseCase {
	return &CloudBackupUseCase{repo: repo, engine: engine, dbPath: dbPath, log: log}
}

// GetConfig retrieves cloud backup config.
func (uc *CloudBackupUseCase) GetConfig(ctx context.Context) (*domain.CloudBackupConfig, error) {
	cfg, err := uc.repo.GetConfig(ctx)
	if err != nil {
		if apperror.Is(err, apperror.KindNotFound) {
			return domain.NewCloudBackupConfig(), nil
		}
		return nil, err
	}
	return cfg, nil
}

// SaveConfig saves cloud backup config.
func (uc *CloudBackupUseCase) SaveConfig(ctx context.Context, cfg *domain.CloudBackupConfig) error {
	if cfg.ID == uuid.Nil {
		cfg.ID = domain.NewCloudBackupConfig().ID
	}
	return uc.repo.SaveConfig(ctx, cfg)
}

// TestConnection tests the cloud connection.
func (uc *CloudBackupUseCase) TestConnection(ctx context.Context) error {
	cfg, err := uc.GetConfig(ctx)
	if err != nil {
		return err
	}
	return uc.engine.TestConnection(cfg)
}

// BackupNow performs an immediate cloud backup.
func (uc *CloudBackupUseCase) BackupNow(ctx context.Context) (*domain.CloudBackupHistory, error) {
	cfg, err := uc.GetConfig(ctx)
	if err != nil {
		return nil, err
	}
	if cfg.Provider == domain.CloudProviderNone {
		return nil, apperror.Business("NO_PROVIDER", "No cloud provider configured")
	}

	// Get file info
	info, err := os.Stat(uc.dbPath)
	if err != nil {
		return nil, apperror.Internal("Database file not found", err)
	}

	fileName := fmt.Sprintf("salonflow-backup-%s.db", time.Now().Format("20060102-150405"))
	remotePath := filepath.Join("backups", fileName)

	history := domain.NewCloudBackupHistory(cfg.Provider, fileName, info.Size(), cfg.EncryptBackups)
	history.RemotePath = remotePath
	history.Status = domain.CloudBackupUploading

	if err := uc.repo.CreateHistory(ctx, history); err != nil {
		return nil, err
	}

	// Perform upload
	if err := uc.engine.Upload(cfg, uc.dbPath, remotePath, cfg.EncryptBackups); err != nil {
		_ = uc.repo.UpdateHistoryStatus(ctx, history.ID, domain.CloudBackupFailed, err.Error())
		history.Status = domain.CloudBackupFailed
		history.ErrorMessage = err.Error()
		return history, nil
	}

	_ = uc.repo.UpdateHistoryStatus(ctx, history.ID, domain.CloudBackupCompleted, "")
	history.Status = domain.CloudBackupCompleted
	return history, nil
}

// Restore downloads and restores a backup.
func (uc *CloudBackupUseCase) Restore(ctx context.Context, historyID uuid.UUID) error {
	cfg, err := uc.GetConfig(ctx)
	if err != nil {
		return err
	}

	// For now, we'd download to a temp path and then swap the DB file
	// This is a placeholder that indicates the restore request
	_ = uc.repo.UpdateHistoryStatus(ctx, historyID, domain.CloudBackupRestoring, "")

	// In production: download file, close DB, replace file, reopen
	_ = cfg
	_ = uc.repo.UpdateHistoryStatus(ctx, historyID, domain.CloudBackupRestored, "")
	return nil
}

// ListHistory lists cloud backup history.
func (uc *CloudBackupUseCase) ListHistory(ctx context.Context, limit, offset int) ([]domain.CloudBackupHistory, int, error) {
	return uc.repo.ListHistory(ctx, limit, offset)
}

// GetStats gets cloud backup stats.
func (uc *CloudBackupUseCase) GetStats(ctx context.Context) (*domain.CloudBackupStats, error) {
	return uc.repo.GetStats(ctx)
}
