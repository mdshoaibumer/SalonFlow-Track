package main

import (
	"context"

	"github.com/google/uuid"
	"github.com/salonflow/salonflow-track/internal/core/domain"
	"github.com/salonflow/salonflow-track/internal/core/usecase"
)

// CustomerService exposes customer operations to the Wails frontend.
type CustomerService struct {
	ctx context.Context
	uc  *usecase.CustomerUseCase
}

func NewCustomerService(uc *usecase.CustomerUseCase) *CustomerService {
	return &CustomerService{uc: uc}
}

func (s *CustomerService) SetContext(ctx context.Context) {
	s.ctx = ctx
}

func (s *CustomerService) ListCustomers(input usecase.ListCustomerInput) (*usecase.ListCustomerOutput, error) {
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
	uid, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return s.uc.Delete(s.ctx, uid)
}

func (s *CustomerService) GetCustomerStats() (*usecase.CustomerStats, error) {
	return s.uc.Stats(s.ctx)
}
