package usecase

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/salonflow/salonflow-track/internal/core/domain"
	"github.com/salonflow/salonflow-track/internal/core/ports"
	"github.com/salonflow/salonflow-track/pkg/apperror"
)

// PerformanceUseCase handles staff performance business logic.
type PerformanceUseCase struct {
	perfRepo ports.PerformanceRepository
	log      *slog.Logger
}

// NewPerformanceUseCase creates a new PerformanceUseCase.
func NewPerformanceUseCase(perfRepo ports.PerformanceRepository, log *slog.Logger) *PerformanceUseCase {
	return &PerformanceUseCase{perfRepo: perfRepo, log: log}
}

// RecordInvoicePerformance updates staff performance when an invoice is created.
func (uc *PerformanceUseCase) RecordInvoicePerformance(ctx context.Context, staffID uuid.UUID, revenue float64, serviceCount int, commission float64) error {
	businessDate := time.Now().UTC().Format("2006-01-02")
	perf := domain.NewStaffPerformanceDaily(staffID, businessDate)
	perf.AddInvoice(revenue, serviceCount, commission)

	if err := uc.perfRepo.Upsert(ctx, perf); err != nil {
		uc.log.Error("failed to record performance", "staff_id", staffID, "error", err)
		return err
	}
	uc.log.Info("performance recorded", "staff_id", staffID, "revenue", revenue, "date", businessDate)
	return nil
}

// DailyPerformanceInput is the input for daily performance query.
type DailyPerformanceInput struct {
	StaffID string `json:"staff_id"`
	Date    string `json:"date"` // YYYY-MM-DD
}

// GetDailyPerformance returns daily performance for all staff.
func (uc *PerformanceUseCase) GetDailyPerformance(ctx context.Context, input DailyPerformanceInput) ([]domain.StaffPerformanceSummary, error) {
	filter := ports.PerformanceFilter{
		StaffID:  input.StaffID,
		DateFrom: input.Date,
		DateTo:   input.Date,
	}
	return uc.perfRepo.GetDaily(ctx, filter)
}

// PeriodPerformanceInput is the input for period-based performance queries.
type PeriodPerformanceInput struct {
	StaffID  string `json:"staff_id"`
	DateFrom string `json:"date_from"`
	DateTo   string `json:"date_to"`
}

// GetWeeklyPerformance returns weekly performance for all staff.
func (uc *PerformanceUseCase) GetWeeklyPerformance(ctx context.Context, input PeriodPerformanceInput) ([]domain.StaffPerformanceSummary, error) {
	filter := ports.PerformanceFilter{
		StaffID:  input.StaffID,
		DateFrom: input.DateFrom,
		DateTo:   input.DateTo,
	}
	return uc.perfRepo.GetWeekly(ctx, filter)
}

// GetMonthlyPerformance returns monthly performance for all staff.
func (uc *PerformanceUseCase) GetMonthlyPerformance(ctx context.Context, input PeriodPerformanceInput) ([]domain.StaffPerformanceSummary, error) {
	filter := ports.PerformanceFilter{
		StaffID:  input.StaffID,
		DateFrom: input.DateFrom,
		DateTo:   input.DateTo,
	}
	return uc.perfRepo.GetMonthly(ctx, filter)
}

// GetStaffPerformance returns raw daily records for a specific staff member.
func (uc *PerformanceUseCase) GetStaffPerformance(ctx context.Context, staffID uuid.UUID, dateFrom, dateTo string) ([]domain.StaffPerformanceDaily, error) {
	if dateFrom == "" {
		now := time.Now().UTC()
		dateFrom = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC).Format("2006-01-02")
		dateTo = now.Format("2006-01-02")
	}
	return uc.perfRepo.GetByStaff(ctx, staffID, dateFrom, dateTo)
}

// TopPerformersInput is the input for top performers query.
type TopPerformersInput struct {
	DateFrom string `json:"date_from"`
	DateTo   string `json:"date_to"`
	Limit    int    `json:"limit"`
}

// GetTopPerformers returns top performing staff.
func (uc *PerformanceUseCase) GetTopPerformers(ctx context.Context, input TopPerformersInput) ([]domain.StaffPerformanceSummary, error) {
	if input.DateFrom == "" {
		now := time.Now().UTC()
		input.DateFrom = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC).Format("2006-01-02")
		input.DateTo = now.Format("2006-01-02")
	}
	if input.Limit <= 0 {
		input.Limit = 10
	}
	return uc.perfRepo.GetTopPerformers(ctx, input.DateFrom, input.DateTo, input.Limit)
}

// GetRevenueTrend returns overall revenue trend.
func (uc *PerformanceUseCase) GetRevenueTrend(ctx context.Context, dateFrom, dateTo string) ([]domain.RevenueTrendPoint, error) {
	if dateFrom == "" {
		now := time.Now().UTC()
		dateFrom = now.AddDate(0, 0, -30).Format("2006-01-02")
		dateTo = now.Format("2006-01-02")
	}
	return uc.perfRepo.GetRevenueTrend(ctx, dateFrom, dateTo)
}

// GetStaffRevenueTrend returns revenue trend for a specific staff member.
func (uc *PerformanceUseCase) GetStaffRevenueTrend(ctx context.Context, staffID uuid.UUID, dateFrom, dateTo string) ([]domain.RevenueTrendPoint, error) {
	if dateFrom == "" {
		now := time.Now().UTC()
		dateFrom = now.AddDate(0, 0, -30).Format("2006-01-02")
		dateTo = now.Format("2006-01-02")
	}
	return uc.perfRepo.GetStaffRevenueTrend(ctx, staffID, dateFrom, dateTo)
}

// PerformanceStats holds dashboard performance statistics.
type PerformanceStats struct {
	TopPerformerToday   *domain.StaffPerformanceSummary `json:"top_performer_today"`
	TopPerformerMonth   *domain.StaffPerformanceSummary `json:"top_performer_month"`
	TotalRevenueToday   float64                         `json:"total_revenue_today"`
	TotalCustomersToday int                             `json:"total_customers_today"`
	AvgBillToday        float64                         `json:"avg_bill_today"`
}

// GetStats returns performance dashboard stats.
func (uc *PerformanceUseCase) GetStats(ctx context.Context) (*PerformanceStats, error) {
	today := time.Now().UTC().Format("2006-01-02")
	now := time.Now().UTC()
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC).Format("2006-01-02")

	stats := &PerformanceStats{}

	// Top performer today
	topToday, err := uc.perfRepo.GetTopPerformers(ctx, today, today, 1)
	if err == nil && len(topToday) > 0 {
		stats.TopPerformerToday = &topToday[0]
		stats.TotalRevenueToday = topToday[0].Revenue
		stats.TotalCustomersToday = topToday[0].CustomerCount
	}

	// Get all today's data for totals
	dailyFilter := ports.PerformanceFilter{DateFrom: today, DateTo: today}
	dailyData, err := uc.perfRepo.GetDaily(ctx, dailyFilter)
	if err == nil {
		var totalRev float64
		var totalCust int
		var totalInv int
		for _, d := range dailyData {
			totalRev += d.Revenue
			totalCust += d.CustomerCount
			totalInv += d.InvoiceCount
		}
		stats.TotalRevenueToday = totalRev
		stats.TotalCustomersToday = totalCust
		if totalInv > 0 {
			stats.AvgBillToday = totalRev / float64(totalInv)
		}
	}

	// Top performer this month
	topMonth, err := uc.perfRepo.GetTopPerformers(ctx, monthStart, today, 1)
	if err == nil && len(topMonth) > 0 {
		stats.TopPerformerMonth = &topMonth[0]
	}

	return stats, nil
}

// Ensure PerformanceUseCase implements interface at compile time.
var _ performanceUseCaseInterface = (*PerformanceUseCase)(nil)

type performanceUseCaseInterface interface {
	RecordInvoicePerformance(ctx context.Context, staffID uuid.UUID, revenue float64, serviceCount int, commission float64) error
}

// CalculateServiceCommission calculates commission based on service commission settings.
func CalculateServiceCommission(commissionType string, commissionValue, revenue float64) float64 {
	switch commissionType {
	case "percentage":
		return revenue * commissionValue / 100
	case "fixed":
		return commissionValue
	default:
		return 0
	}
}

// Validate UUID helper (unexported, used internally).
func parseUUID(id string) (uuid.UUID, error) {
	parsed, err := uuid.Parse(id)
	if err != nil {
		return uuid.Nil, &apperror.Error{Kind: apperror.KindValidation, Message: "invalid UUID: " + id}
	}
	return parsed, nil
}
