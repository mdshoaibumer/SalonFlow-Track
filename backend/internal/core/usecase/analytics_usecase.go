package usecase

import (
	"context"
	"log/slog"

	"github.com/salonflow/salonflow-track/internal/core/domain"
	"github.com/salonflow/salonflow-track/internal/core/ports"
)

// AnalyticsUseCase provides business analytics operations.
type AnalyticsUseCase struct {
	repo ports.AnalyticsRepository
	log  *slog.Logger
}

// NewAnalyticsUseCase creates a new AnalyticsUseCase.
func NewAnalyticsUseCase(repo ports.AnalyticsRepository, log *slog.Logger) *AnalyticsUseCase {
	return &AnalyticsUseCase{repo: repo, log: log}
}

// GetDashboard returns executive dashboard stats.
func (uc *AnalyticsUseCase) GetDashboard(ctx context.Context) (*domain.DashboardStats, error) {
	return uc.repo.GetDashboardStats(ctx)
}

// GetKPIs returns KPI metrics for a date range.
func (uc *AnalyticsUseCase) GetKPIs(ctx context.Context, dateFrom, dateTo string) (*domain.KPIMetrics, error) {
	return uc.repo.GetKPIMetrics(ctx, dateFrom, dateTo)
}

// GetRevenueReport returns revenue analytics.
func (uc *AnalyticsUseCase) GetRevenueReport(ctx context.Context, dateFrom, dateTo, groupBy string) (*domain.RevenueReport, error) {
	return uc.repo.GetRevenueReport(ctx, dateFrom, dateTo, groupBy)
}

// GetCustomerReport returns customer analytics.
func (uc *AnalyticsUseCase) GetCustomerReport(ctx context.Context, dateFrom, dateTo string) (*domain.CustomerReport, error) {
	return uc.repo.GetCustomerReport(ctx, dateFrom, dateTo)
}

// GetStaffReport returns staff analytics.
func (uc *AnalyticsUseCase) GetStaffReport(ctx context.Context, dateFrom, dateTo string) (*domain.StaffReport, error) {
	return uc.repo.GetStaffReport(ctx, dateFrom, dateTo)
}

// GetServiceReport returns service analytics.
func (uc *AnalyticsUseCase) GetServiceReport(ctx context.Context, dateFrom, dateTo string) (*domain.ServiceReport, error) {
	return uc.repo.GetServiceReport(ctx, dateFrom, dateTo)
}

// GetExpenseReport returns expense analytics.
func (uc *AnalyticsUseCase) GetExpenseReport(ctx context.Context, dateFrom, dateTo string) (*domain.ExpenseReport, error) {
	return uc.repo.GetExpenseReport(ctx, dateFrom, dateTo)
}

// GetInventoryReport returns inventory analytics.
func (uc *AnalyticsUseCase) GetInventoryReport(ctx context.Context) (*domain.InventoryReport, error) {
	return uc.repo.GetInventoryReport(ctx)
}

// GetProfitLossReport returns P&L data.
func (uc *AnalyticsUseCase) GetProfitLossReport(ctx context.Context, dateFrom, dateTo, groupBy string) (*domain.ProfitLossReport, error) {
	return uc.repo.GetProfitLossReport(ctx, dateFrom, dateTo, groupBy)
}
