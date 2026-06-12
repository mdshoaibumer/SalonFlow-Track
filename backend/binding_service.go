package main

import (
	"context"

	"github.com/google/uuid"
	"github.com/salonflow/salonflow-track/internal/core/domain"
	"github.com/salonflow/salonflow-track/internal/core/usecase"
)

// ServiceService exposes salon service operations to the Wails frontend.
type ServiceService struct {
	ctx   context.Context
	uc    *usecase.ServiceUseCase
	guard *PermissionGuard
}

func NewServiceService(uc *usecase.ServiceUseCase) *ServiceService {
	return &ServiceService{uc: uc}
}

func (s *ServiceService) SetContext(ctx context.Context) {
	s.ctx = ctx
}

func (s *ServiceService) ListServices(input usecase.ListServiceInput) (*usecase.ListServiceOutput, error) {
	return s.uc.List(s.ctx, input)
}

func (s *ServiceService) GetService(id string) (*domain.Service, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	return s.uc.GetByID(s.ctx, uid)
}

func (s *ServiceService) CreateService(input usecase.CreateServiceInput) (*domain.Service, error) {
	return s.uc.Create(s.ctx, input)
}

func (s *ServiceService) UpdateService(id string, input usecase.UpdateServiceInput) (*domain.Service, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	return s.uc.Update(s.ctx, uid, input)
}

func (s *ServiceService) DeleteService(id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return s.uc.Delete(s.ctx, uid)
}

func (s *ServiceService) GetServiceStats() (*usecase.ServiceStats, error) {
	return s.uc.Stats(s.ctx)
}
