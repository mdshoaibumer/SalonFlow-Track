package ports

import (
	"context"

	"github.com/salonflow/salonflow-track/internal/core/domain"
)

// AnalyticsRepository provides read-only analytics queries across all domain tables.
type AnalyticsRepository interface {
	// Dashboard
	GetDashboardStats(ctx context.Context) (*domain.DashboardStats, error)
	GetKPIMetrics(ctx context.Context, dateFrom, dateTo string) (*domain.KPIMetrics, error)

	// Revenue
	GetRevenueReport(ctx context.Context, dateFrom, dateTo string, groupBy string) (*domain.RevenueReport, error)

	// Customer
	GetCustomerReport(ctx context.Context, dateFrom, dateTo string) (*domain.CustomerReport, error)

	// Staff
	GetStaffReport(ctx context.Context, dateFrom, dateTo string) (*domain.StaffReport, error)

	// Service
	GetServiceReport(ctx context.Context, dateFrom, dateTo string) (*domain.ServiceReport, error)

	// Expense
	GetExpenseReport(ctx context.Context, dateFrom, dateTo string) (*domain.ExpenseReport, error)

	// Inventory
	GetInventoryReport(ctx context.Context) (*domain.InventoryReport, error)

	// P&L
	GetProfitLossReport(ctx context.Context, dateFrom, dateTo string, groupBy string) (*domain.ProfitLossReport, error)
}
