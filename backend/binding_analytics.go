package main

import (
	"context"

	"github.com/salonflow/salonflow-track/internal/core/domain"
	"github.com/salonflow/salonflow-track/internal/core/usecase"
)

// AnalyticsService exposes analytics/reports to the Wails frontend.
type AnalyticsService struct {
	ctx context.Context
	uc  *usecase.AnalyticsUseCase
}

func NewAnalyticsService(uc *usecase.AnalyticsUseCase) *AnalyticsService {
	return &AnalyticsService{uc: uc}
}

func (s *AnalyticsService) SetContext(ctx context.Context) {
	s.ctx = ctx
}

func (s *AnalyticsService) GetDashboard() (*domain.DashboardStats, error) {
	return s.uc.GetDashboard(s.ctx)
}

func (s *AnalyticsService) GetKPIs(dateFrom, dateTo string) (*domain.KPIMetrics, error) {
	return s.uc.GetKPIs(s.ctx, dateFrom, dateTo)
}

func (s *AnalyticsService) GetRevenueReport(dateFrom, dateTo, groupBy string) (*domain.RevenueReport, error) {
	return s.uc.GetRevenueReport(s.ctx, dateFrom, dateTo, groupBy)
}

func (s *AnalyticsService) GetCustomerReport(dateFrom, dateTo string) (*domain.CustomerReport, error) {
	return s.uc.GetCustomerReport(s.ctx, dateFrom, dateTo)
}

func (s *AnalyticsService) GetStaffReport(dateFrom, dateTo string) (*domain.StaffReport, error) {
	return s.uc.GetStaffReport(s.ctx, dateFrom, dateTo)
}

func (s *AnalyticsService) GetServiceReport(dateFrom, dateTo string) (*domain.ServiceReport, error) {
	return s.uc.GetServiceReport(s.ctx, dateFrom, dateTo)
}

func (s *AnalyticsService) GetExpenseAnalytics(dateFrom, dateTo string) (*domain.ExpenseReport, error) {
	return s.uc.GetExpenseReport(s.ctx, dateFrom, dateTo)
}

func (s *AnalyticsService) GetInventoryReport() (*domain.InventoryReport, error) {
	return s.uc.GetInventoryReport(s.ctx)
}

func (s *AnalyticsService) GetProfitLossReport(dateFrom, dateTo, groupBy string) (*domain.ProfitLossReport, error) {
	return s.uc.GetProfitLossReport(s.ctx, dateFrom, dateTo, groupBy)
}
