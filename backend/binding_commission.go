package main

import (
	"context"

	"github.com/google/uuid"
	"github.com/salonflow/salonflow-track/internal/core/domain"
	"github.com/salonflow/salonflow-track/internal/core/ports"
	"github.com/salonflow/salonflow-track/internal/core/usecase"
)

// CommissionService exposes commission operations to the Wails frontend.
type CommissionService struct {
	ctx context.Context
	uc  *usecase.CommissionUseCase
}

func NewCommissionService(uc *usecase.CommissionUseCase) *CommissionService {
	return &CommissionService{uc: uc}
}

func (s *CommissionService) SetContext(ctx context.Context) {
	s.ctx = ctx
}

func (s *CommissionService) CreateRule(input usecase.CreateRuleInput) (*domain.CommissionRule, error) {
	return s.uc.CreateRule(s.ctx, input)
}

func (s *CommissionService) GetRule(id string) (*domain.CommissionRule, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	return s.uc.GetRuleByID(s.ctx, uid)
}

func (s *CommissionService) ListRules(input usecase.ListRulesInput) (*usecase.ListRulesOutput, error) {
	return s.uc.ListRules(s.ctx, input)
}

func (s *CommissionService) UpdateRule(id string, input usecase.UpdateRuleInput) (*domain.CommissionRule, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	return s.uc.UpdateRule(s.ctx, uid, input)
}

func (s *CommissionService) DeleteRule(id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return s.uc.DeleteRule(s.ctx, uid)
}

func (s *CommissionService) GetStaffCommission(input usecase.GetStaffCommissionInput) (*usecase.StaffCommissionOutput, error) {
	return s.uc.GetStaffCommission(s.ctx, input)
}

func (s *CommissionService) GetMonthlyCommission(input usecase.MonthlyCommissionInput) ([]ports.CommissionStaffSummary, error) {
	return s.uc.GetMonthlyCommission(s.ctx, input)
}

func (s *CommissionService) GetCommissionStats() (*usecase.CommissionStats, error) {
	return s.uc.GetStats(s.ctx)
}
