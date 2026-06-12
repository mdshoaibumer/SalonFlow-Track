package main

import (
	"context"

	"github.com/google/uuid"
	"github.com/salonflow/salonflow-track/internal/core/domain"
	"github.com/salonflow/salonflow-track/internal/core/usecase"
)

// StaffService exposes staff operations to the Wails frontend.
type StaffService struct {
	ctx   context.Context
	uc    *usecase.StaffUseCase
	guard *PermissionGuard
}

func NewStaffService(uc *usecase.StaffUseCase) *StaffService {
	return &StaffService{uc: uc}
}

func (s *StaffService) SetContext(ctx context.Context) {
	s.ctx = ctx
}

func (s *StaffService) ListStaff(input usecase.ListStaffInput) (*usecase.ListStaffOutput, error) {
	return s.uc.List(s.ctx, input)
}

func (s *StaffService) GetStaff(id string) (*domain.Staff, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	return s.uc.GetByID(s.ctx, uid)
}

func (s *StaffService) CreateStaff(input usecase.CreateStaffInput) (*domain.Staff, error) {
	return s.uc.Create(s.ctx, input)
}

func (s *StaffService) UpdateStaff(id string, input usecase.UpdateStaffInput) (*domain.Staff, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	return s.uc.Update(s.ctx, uid, input)
}

func (s *StaffService) DeleteStaff(id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return s.uc.Delete(s.ctx, uid)
}

func (s *StaffService) GetStaffStats() (*usecase.StaffStats, error) {
	return s.uc.Stats(s.ctx)
}
