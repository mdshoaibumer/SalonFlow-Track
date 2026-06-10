package main

import (
	"context"

	"github.com/google/uuid"
	"github.com/salonflow/salonflow-track/internal/core/domain"
	"github.com/salonflow/salonflow-track/internal/core/usecase"
)

// LicenseService exposes license operations to the Wails frontend.
type LicenseService struct {
	ctx context.Context
	uc  *usecase.LicenseUseCase
}

func NewLicenseService(uc *usecase.LicenseUseCase) *LicenseService {
	return &LicenseService{uc: uc}
}

func (s *LicenseService) SetContext(ctx context.Context) {
	s.ctx = ctx
}

func (s *LicenseService) GetStatus() (*domain.LicenseStatus, error) {
	return s.uc.GetStatus(s.ctx)
}

func (s *LicenseService) Validate() (*domain.LicenseValidation, error) {
	return s.uc.Validate(s.ctx)
}

func (s *LicenseService) Activate(key, customerName, salonName string) (*domain.License, error) {
	return s.uc.Activate(s.ctx, key, customerName, salonName)
}

func (s *LicenseService) Renew(key string) (*domain.License, error) {
	return s.uc.Renew(s.ctx, key)
}

func (s *LicenseService) ListEvents(page, perPage int) ([]domain.LicenseEvent, int, error) {
	return s.uc.ListEvents(s.ctx, page, perPage)
}

// UpdateService exposes update operations to the Wails frontend.
type UpdateService struct {
	ctx context.Context
	uc  *usecase.UpdateUseCase
}

func NewUpdateService(uc *usecase.UpdateUseCase) *UpdateService {
	return &UpdateService{uc: uc}
}

func (s *UpdateService) SetContext(ctx context.Context) {
	s.ctx = ctx
}

func (s *UpdateService) CheckForUpdate() (*domain.UpdateStatus, error) {
	return s.uc.CheckForUpdate(s.ctx)
}

func (s *UpdateService) DownloadUpdate() (*domain.UpdateRecord, error) {
	return s.uc.DownloadUpdate(s.ctx)
}

func (s *UpdateService) InstallUpdate() (*domain.UpdateRecord, error) {
	return s.uc.InstallUpdate(s.ctx)
}

func (s *UpdateService) GetUpdateStatus() (*domain.UpdateStatus, error) {
	return s.uc.GetUpdateStatus(s.ctx)
}

func (s *UpdateService) ListUpdateHistory(page, perPage int) ([]domain.UpdateRecord, int, error) {
	return s.uc.ListHistory(s.ctx, page, perPage)
}

// ImportService exposes import operations to the Wails frontend.
type ImportService struct {
	ctx context.Context
	uc  *usecase.ImportUseCase
}

func NewImportService(uc *usecase.ImportUseCase) *ImportService {
	return &ImportService{uc: uc}
}

func (s *ImportService) SetContext(ctx context.Context) {
	s.ctx = ctx
}

func (s *ImportService) Upload(fileName, filePath, targetEntity string) (*domain.ImportJob, []string, []domain.ColumnMapping, error) {
	return s.uc.Upload(s.ctx, fileName, filePath, targetEntity)
}

func (s *ImportService) Validate(jobID string, mappings []domain.ColumnMapping) (*domain.ImportPreview, error) {
	uid, err := uuid.Parse(jobID)
	if err != nil {
		return nil, err
	}
	return s.uc.Validate(s.ctx, uid, mappings)
}

func (s *ImportService) Process(jobID string) (*domain.ImportJob, error) {
	uid, err := uuid.Parse(jobID)
	if err != nil {
		return nil, err
	}
	return s.uc.Process(s.ctx, uid)
}

func (s *ImportService) ListJobs(page, perPage int) ([]domain.ImportJob, int, error) {
	return s.uc.ListJobs(s.ctx, page, perPage)
}

func (s *ImportService) GetJob(id string) (*domain.ImportJob, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	return s.uc.GetJob(s.ctx, uid)
}

func (s *ImportService) ListLogs(jobID, status string, page, perPage int) ([]domain.ImportLog, int, error) {
	uid, err := uuid.Parse(jobID)
	if err != nil {
		return nil, 0, err
	}
	return s.uc.ListLogs(s.ctx, uid, status, page, perPage)
}
