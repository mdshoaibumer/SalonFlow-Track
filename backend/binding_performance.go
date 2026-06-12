package main

import (
	"context"

	"github.com/google/uuid"
	"github.com/salonflow/salonflow-track/internal/core/domain"
	"github.com/salonflow/salonflow-track/internal/core/usecase"
)

// PerformanceService exposes performance tracking to the Wails frontend.
type PerformanceService struct {
	ctx   context.Context
	uc    *usecase.PerformanceUseCase
	guard *PermissionGuard
}

func NewPerformanceService(uc *usecase.PerformanceUseCase) *PerformanceService {
	return &PerformanceService{uc: uc}
}

func (s *PerformanceService) SetContext(ctx context.Context) {
	s.ctx = ctx
}

func (s *PerformanceService) GetDailyPerformance(input usecase.DailyPerformanceInput) ([]domain.StaffPerformanceSummary, error) {
	return s.uc.GetDailyPerformance(s.ctx, input)
}

func (s *PerformanceService) GetWeeklyPerformance(input usecase.PeriodPerformanceInput) ([]domain.StaffPerformanceSummary, error) {
	return s.uc.GetWeeklyPerformance(s.ctx, input)
}

func (s *PerformanceService) GetMonthlyPerformance(input usecase.PeriodPerformanceInput) ([]domain.StaffPerformanceSummary, error) {
	return s.uc.GetMonthlyPerformance(s.ctx, input)
}

func (s *PerformanceService) GetTopPerformers(input usecase.TopPerformersInput) ([]domain.StaffPerformanceSummary, error) {
	return s.uc.GetTopPerformers(s.ctx, input)
}

func (s *PerformanceService) GetRevenueTrend(dateFrom, dateTo string) ([]domain.RevenueTrendPoint, error) {
	return s.uc.GetRevenueTrend(s.ctx, dateFrom, dateTo)
}

func (s *PerformanceService) GetStaffRevenueTrend(staffID string, dateFrom, dateTo string) ([]domain.RevenueTrendPoint, error) {
	uid, err := uuid.Parse(staffID)
	if err != nil {
		return nil, err
	}
	return s.uc.GetStaffRevenueTrend(s.ctx, uid, dateFrom, dateTo)
}

func (s *PerformanceService) GetPerformanceStats() (*usecase.PerformanceStats, error) {
	return s.uc.GetStats(s.ctx)
}

func (s *PerformanceService) GetStaffPerformance(staffID string, dateFrom, dateTo string) ([]domain.StaffPerformanceDaily, error) {
	uid, err := uuid.Parse(staffID)
	if err != nil {
		return nil, err
	}
	return s.uc.GetStaffPerformance(s.ctx, uid, dateFrom, dateTo)
}
