package main

import (
	"context"

	"github.com/google/uuid"
	"github.com/salonflow/salonflow-track/internal/core/domain"
	"github.com/salonflow/salonflow-track/internal/core/usecase"
)

// ExpenseService exposes expense operations to the Wails frontend.
type ExpenseService struct {
	ctx context.Context
	uc  *usecase.ExpenseUseCase
}

func NewExpenseService(uc *usecase.ExpenseUseCase) *ExpenseService {
	return &ExpenseService{uc: uc}
}

func (s *ExpenseService) SetContext(ctx context.Context) {
	s.ctx = ctx
}

func (s *ExpenseService) ListCategories(activeOnly bool) ([]domain.ExpenseCategory, error) {
	return s.uc.ListCategories(s.ctx, activeOnly)
}

func (s *ExpenseService) CreateCategory(name, description string) (*domain.ExpenseCategory, error) {
	return s.uc.CreateCategory(s.ctx, name, description)
}

func (s *ExpenseService) UpdateCategory(id string, name, description string, isActive bool) (*domain.ExpenseCategory, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	return s.uc.UpdateCategory(s.ctx, uid, name, description, isActive)
}

func (s *ExpenseService) CreateExpense(input usecase.CreateExpenseInput) (*domain.Expense, error) {
	return s.uc.CreateExpense(s.ctx, input)
}

func (s *ExpenseService) GetExpense(id string) (*domain.Expense, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	return s.uc.GetExpense(s.ctx, uid)
}

func (s *ExpenseService) UpdateExpense(id string, input usecase.UpdateExpenseInput) (*domain.Expense, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	return s.uc.UpdateExpense(s.ctx, uid, input)
}

func (s *ExpenseService) DeleteExpense(id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return s.uc.DeleteExpense(s.ctx, uid)
}

func (s *ExpenseService) ListExpenses(input usecase.ListExpensesInput) (*usecase.ListExpensesOutput, error) {
	return s.uc.ListExpenses(s.ctx, input)
}

func (s *ExpenseService) GetExpenseStats() (*domain.ExpenseStats, error) {
	return s.uc.GetExpenseStats(s.ctx)
}

func (s *ExpenseService) GetExpenseReport(input usecase.ExpenseReportInput) (*usecase.ExpenseReport, error) {
	return s.uc.GetExpenseReport(s.ctx, input)
}

func (s *ExpenseService) GetProfitLoss(input usecase.ProfitLossInput) (*domain.ProfitLoss, error) {
	return s.uc.GetProfitLoss(s.ctx, input)
}

func (s *ExpenseService) GetMonthlyTrend(months int) ([]domain.MonthlyTrend, error) {
	return s.uc.GetMonthlyTrend(s.ctx, months)
}
