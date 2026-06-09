package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/salonflow/salonflow-track/internal/core/domain"
)

// PerformanceRepository defines persistence operations for staff performance.
type PerformanceRepository interface {
	Upsert(ctx context.Context, perf *domain.StaffPerformanceDaily) error
	GetDaily(ctx context.Context, filter PerformanceFilter) ([]domain.StaffPerformanceSummary, error)
	GetWeekly(ctx context.Context, filter PerformanceFilter) ([]domain.StaffPerformanceSummary, error)
	GetMonthly(ctx context.Context, filter PerformanceFilter) ([]domain.StaffPerformanceSummary, error)
	GetByStaff(ctx context.Context, staffID uuid.UUID, dateFrom, dateTo string) ([]domain.StaffPerformanceDaily, error)
	GetTopPerformers(ctx context.Context, dateFrom, dateTo string, limit int) ([]domain.StaffPerformanceSummary, error)
	GetRevenueTrend(ctx context.Context, dateFrom, dateTo string) ([]domain.RevenueTrendPoint, error)
	GetStaffRevenueTrend(ctx context.Context, staffID uuid.UUID, dateFrom, dateTo string) ([]domain.RevenueTrendPoint, error)
}

// PerformanceFilter holds query parameters for performance queries.
type PerformanceFilter struct {
	StaffID  string
	DateFrom string
	DateTo   string
}
