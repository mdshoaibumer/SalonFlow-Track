package main

import (
	"context"

	"github.com/google/uuid"
	"github.com/salonflow/salonflow-track/internal/core/domain"
	"github.com/salonflow/salonflow-track/internal/core/usecase"
)

// BackupService exposes backup operations to the Wails frontend.
type BackupService struct {
	ctx context.Context
	uc  *usecase.BackupUseCase
}

func NewBackupService(uc *usecase.BackupUseCase) *BackupService {
	return &BackupService{uc: uc}
}

func (s *BackupService) SetContext(ctx context.Context) {
	s.ctx = ctx
}

func (s *BackupService) CreateBackup(backupType string) (*domain.BackupRecord, error) {
	return s.uc.CreateBackup(s.ctx, backupType)
}

func (s *BackupService) VerifyBackup(id string) (*domain.BackupVerification, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	return s.uc.VerifyBackup(s.ctx, uid)
}

func (s *BackupService) RestoreBackup(id string, notes string) (*domain.RestoreRecord, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	return s.uc.RestoreBackup(s.ctx, uid, notes)
}

func (s *BackupService) ListBackups(page, perPage int) ([]domain.BackupRecord, int, error) {
	return s.uc.ListBackups(s.ctx, page, perPage)
}

func (s *BackupService) ListRestores(page, perPage int) ([]domain.RestoreRecord, int, error) {
	return s.uc.ListRestores(s.ctx, page, perPage)
}

func (s *BackupService) GetBackupStats() (*domain.BackupStats, error) {
	return s.uc.GetStats(s.ctx)
}

func (s *BackupService) DeleteBackup(id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return s.uc.DeleteBackup(s.ctx, uid)
}
