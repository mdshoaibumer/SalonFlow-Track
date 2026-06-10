package usecase

import (
	"context"
	"log/slog"
	"os"

	"github.com/google/uuid"
	"github.com/salonflow/salonflow-track/internal/adapters/backup"
	"github.com/salonflow/salonflow-track/internal/core/domain"
	"github.com/salonflow/salonflow-track/internal/core/ports"
	"github.com/salonflow/salonflow-track/pkg/apperror"
)

// BackupUseCase manages backup and restore operations.
type BackupUseCase struct {
	repo   ports.BackupRepository
	engine ports.BackupEngine
	dbPath string
	log    *slog.Logger
}

// NewBackupUseCase creates a new BackupUseCase.
func NewBackupUseCase(repo ports.BackupRepository, engine ports.BackupEngine, dbPath string, log *slog.Logger) *BackupUseCase {
	return &BackupUseCase{repo: repo, engine: engine, dbPath: dbPath, log: log}
}

// CreateBackup creates a new backup.
func (uc *BackupUseCase) CreateBackup(ctx context.Context, backupType string) (*domain.BackupRecord, error) {
	if backupType == "" {
		backupType = domain.BackupTypeManual
	}

	name := backup.GenerateBackupName(backupType)
	path := backup.GenerateBackupPath(uc.engine.BackupDir())

	record := domain.NewBackupRecord(name, backupType, path)
	if err := uc.repo.CreateBackupRecord(ctx, record); err != nil {
		return nil, err
	}

	uc.log.Info("creating backup", "name", name, "type", backupType, "path", path)

	size, checksum, err := uc.engine.CreateBackup(uc.dbPath, path)
	if err != nil {
		record.Status = domain.BackupStatusFailed
		record.ErrorMessage = err.Error()
		_ = uc.repo.UpdateBackupRecord(ctx, record)
		uc.log.Error("backup failed", "error", err)
		return record, apperror.Internal("backup failed", err)
	}

	record.FileSize = size
	record.Checksum = checksum
	record.Status = domain.BackupStatusCompleted
	_ = uc.repo.UpdateBackupRecord(ctx, record)

	uc.log.Info("backup completed", "name", name, "size", size)
	return record, nil
}

// VerifyBackup verifies a specific backup.
func (uc *BackupUseCase) VerifyBackup(ctx context.Context, id uuid.UUID) (*domain.BackupVerification, error) {
	record, err := uc.repo.GetBackupByID(ctx, id)
	if err != nil {
		return nil, err
	}

	v := uc.engine.VerifyBackup(record.BackupPath, record.Checksum)
	v.BackupID = id.String()

	// Update record status based on verification
	if v.Status == domain.BackupStatusVerified {
		record.Status = domain.BackupStatusVerified
	} else {
		record.Status = domain.BackupStatusCorrupted
		record.ErrorMessage = v.ErrorMessage
	}
	_ = uc.repo.UpdateBackupRecord(ctx, record)

	return v, nil
}

// RestoreBackup restores from a backup.
func (uc *BackupUseCase) RestoreBackup(ctx context.Context, id uuid.UUID, notes string) (*domain.RestoreRecord, error) {
	backupRec, err := uc.repo.GetBackupByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Verify backup first
	v := uc.engine.VerifyBackup(backupRec.BackupPath, backupRec.Checksum)
	if v.Status == domain.BackupStatusCorrupted {
		return nil, apperror.Business("corrupted_backup", "backup is corrupted and cannot be restored: "+v.ErrorMessage)
	}

	// Create a safety backup before restore
	uc.log.Info("creating safety backup before restore")
	_, safetyErr := uc.CreateBackup(ctx, domain.BackupTypeBeforeRestore)
	if safetyErr != nil {
		uc.log.Warn("safety backup failed, proceeding with restore", "error", safetyErr)
	}

	// Create restore record
	restoreRec := domain.NewRestoreRecord(backupRec.ID, backupRec.BackupName, notes)
	if err := uc.repo.CreateRestoreRecord(ctx, restoreRec); err != nil {
		return nil, err
	}

	uc.log.Info("restoring backup", "backup_name", backupRec.BackupName, "path", backupRec.BackupPath)

	// Perform restore
	if err := uc.engine.RestoreBackup(backupRec.BackupPath, uc.dbPath); err != nil {
		restoreRec.Status = domain.RestoreStatusFailed
		restoreRec.ErrorMessage = err.Error()
		_ = uc.repo.UpdateRestoreRecord(ctx, restoreRec)
		uc.log.Error("restore failed", "error", err)
		return restoreRec, apperror.Internal("restore failed", err)
	}

	restoreRec.Status = domain.RestoreStatusCompleted
	_ = uc.repo.UpdateRestoreRecord(ctx, restoreRec)
	uc.log.Info("restore completed", "backup_name", backupRec.BackupName)

	return restoreRec, nil
}

// ListBackups returns backup history.
func (uc *BackupUseCase) ListBackups(ctx context.Context, page, perPage int) ([]domain.BackupRecord, int, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}
	return uc.repo.ListBackups(ctx, perPage, (page-1)*perPage)
}

// ListRestores returns restore history.
func (uc *BackupUseCase) ListRestores(ctx context.Context, page, perPage int) ([]domain.RestoreRecord, int, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}
	return uc.repo.ListRestores(ctx, perPage, (page-1)*perPage)
}

// GetStats returns backup dashboard stats.
func (uc *BackupUseCase) GetStats(ctx context.Context) (*domain.BackupStats, error) {
	return uc.repo.GetBackupStats(ctx)
}

// DeleteBackup deletes a backup record and its file.
func (uc *BackupUseCase) DeleteBackup(ctx context.Context, id uuid.UUID) error {
	record, err := uc.repo.GetBackupByID(ctx, id)
	if err != nil {
		return err
	}

	// Remove file if exists
	if record.BackupPath != "" {
		os.Remove(record.BackupPath)
	}

	return uc.repo.DeleteBackupRecord(ctx, id)
}
