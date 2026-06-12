package main

import (
	"context"

	"github.com/google/uuid"
	"github.com/salonflow/salonflow-track/internal/core/domain"
	"github.com/salonflow/salonflow-track/internal/core/usecase"
)

// CustomerService exposes customer operations to the Wails frontend.
type CustomerService struct {
	ctx      context.Context
	uc       *usecase.CustomerUseCase
	guard    *PermissionGuard
	licGuard *LicenseGuard
}

func NewCustomerService(uc *usecase.CustomerUseCase) *CustomerService {
	return &CustomerService{uc: uc}
}

func (s *CustomerService) SetContext(ctx context.Context) {
	s.ctx = ctx
}

func (s *CustomerService) ListCustomers(input usecase.ListCustomerInput) (*usecase.ListCustomerOutput, error) {
	if err := s.guard.Require("customers.read"); err != nil {
		return nil, err
	}
	return s.uc.List(s.ctx, input)
}

func (s *CustomerService) GetCustomer(id string) (*domain.Customer, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	return s.uc.GetByID(s.ctx, uid)
}

func (s *CustomerService) CreateCustomer(input usecase.CreateCustomerInput) (*domain.Customer, error) {
	if err := s.guard.Require("customers.create"); err != nil {
		return nil, err
	}
	if err := s.licGuard.RequireActive(domain.OpCustomerCreate); err != nil {
		return nil, err
	}
	return s.uc.Create(s.ctx, input)
}

func (s *CustomerService) UpdateCustomer(id string, input usecase.UpdateCustomerInput) (*domain.Customer, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	return s.uc.Update(s.ctx, uid, input)
}

func (s *CustomerService) DeleteCustomer(id string) error {
	if err := s.guard.Require("customers.delete"); err != nil {
		return err
	}
	uid, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return s.uc.Delete(s.ctx, uid)
}

func (s *CustomerService) GetCustomerStats() (*usecase.CustomerStats, error) {
	return s.uc.Stats(s.ctx)
}
