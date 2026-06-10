package main

import (
	"context"

	"github.com/google/uuid"
	"github.com/salonflow/salonflow-track/internal/core/domain"
	"github.com/salonflow/salonflow-track/internal/core/usecase"
)

// SalaryService exposes salary operations to the Wails frontend.
type SalaryService struct {
	ctx context.Context
	uc  *usecase.SalaryUseCase
}

func NewSalaryService(uc *usecase.SalaryUseCase) *SalaryService {
	return &SalaryService{uc: uc}
}

func (s *SalaryService) SetContext(ctx context.Context) {
	s.ctx = ctx
}

func (s *SalaryService) GenerateSalary(input usecase.GenerateSalaryInput) (*usecase.GenerateSalaryOutput, error) {
	return s.uc.GenerateMonthlySalary(s.ctx, input)
}

func (s *SalaryService) GetSalary(id string) (*domain.SalaryRecord, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	return s.uc.GetSalaryByID(s.ctx, uid)
}

func (s *SalaryService) ListSalaries(input usecase.ListSalariesInput) ([]domain.SalaryRecord, error) {
	return s.uc.ListSalaries(s.ctx, input)
}

func (s *SalaryService) PaySalary(id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return s.uc.PaySalary(s.ctx, uid)
}

func (s *SalaryService) ListCycles(year int) ([]domain.SalaryCycle, error) {
	return s.uc.ListCycles(s.ctx, year)
}

func (s *SalaryService) CreateAdvance(input usecase.CreateAdvanceInput) (*domain.Advance, error) {
	return s.uc.CreateAdvance(s.ctx, input)
}

func (s *SalaryService) ApproveAdvance(id string) (*domain.Advance, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	return s.uc.ApproveAdvance(s.ctx, uid)
}

func (s *SalaryService) RejectAdvance(id string) (*domain.Advance, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	return s.uc.RejectAdvance(s.ctx, uid)
}

func (s *SalaryService) ListAdvances(input usecase.ListAdvancesInput) (*usecase.ListAdvancesOutput, error) {
	return s.uc.ListAdvances(s.ctx, input)
}

func (s *SalaryService) GetSalaryStats() (*usecase.SalaryStats, error) {
	return s.uc.GetStats(s.ctx)
}
