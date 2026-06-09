package sqlite

import (
	"context"
	"database/sql"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/salonflow/salonflow-track/internal/core/domain"
	"github.com/salonflow/salonflow-track/internal/core/ports"
	"github.com/salonflow/salonflow-track/pkg/apperror"
)

// PerformanceRepository is the SQLite implementation of ports.PerformanceRepository.
type PerformanceRepository struct {
	db  *sql.DB
	log *slog.Logger
}

// NewPerformanceRepository creates a new PerformanceRepository.
func NewPerformanceRepository(db *sql.DB, log *slog.Logger) *PerformanceRepository {
	return &PerformanceRepository{db: db, log: log}
}

// Upsert inserts or updates a daily performance record.
func (r *PerformanceRepository) Upsert(ctx context.Context, perf *domain.StaffPerformanceDaily) error {
	query := `
		INSERT INTO staff_performance_daily (id, staff_id, business_date, invoice_count, customer_count, service_count, revenue, commission_amount, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(staff_id, business_date) DO UPDATE SET
			invoice_count = invoice_count + excluded.invoice_count,
			customer_count = customer_count + excluded.customer_count,
			service_count = service_count + excluded.service_count,
			revenue = revenue + excluded.revenue,
			commission_amount = commission_amount + excluded.commission_amount,
			updated_at = excluded.updated_at`

	_, err := r.db.ExecContext(ctx, query,
		perf.ID.String(),
		perf.StaffID.String(),
		perf.BusinessDate,
		perf.InvoiceCount,
		perf.CustomerCount,
		perf.ServiceCount,
		perf.Revenue,
		perf.CommissionAmount,
		perf.CreatedAt.Format(time.RFC3339),
		perf.UpdatedAt.Format(time.RFC3339),
	)
	if err != nil {
		return apperror.Database("upsert performance", err)
	}
	return nil
}

// GetDaily returns daily performance summaries.
func (r *PerformanceRepository) GetDaily(ctx context.Context, filter ports.PerformanceFilter) ([]domain.StaffPerformanceSummary, error) {
	dateFrom := filter.DateFrom
	dateTo := filter.DateTo
	if dateFrom == "" {
		dateFrom = time.Now().UTC().Format("2006-01-02")
	}
	if dateTo == "" {
		dateTo = dateFrom
	}

	query := `
		SELECT s.id, s.full_name,
			COALESCE(SUM(p.revenue), 0),
			COALESCE(SUM(p.customer_count), 0),
			COALESCE(SUM(p.invoice_count), 0),
			COALESCE(SUM(p.service_count), 0),
			COALESCE(SUM(p.commission_amount), 0)
		FROM staff s
		LEFT JOIN staff_performance_daily p ON s.id = p.staff_id AND p.business_date BETWEEN ? AND ?
		WHERE s.status = 'active'`

	args := []interface{}{dateFrom, dateTo}

	if filter.StaffID != "" {
		query += " AND s.id = ?"
		args = append(args, filter.StaffID)
	}

	query += " GROUP BY s.id, s.full_name ORDER BY COALESCE(SUM(p.revenue), 0) DESC"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, apperror.Database("get daily performance", err)
	}
	defer rows.Close()

	return r.scanSummaries(rows)
}

// GetWeekly returns weekly performance summaries (last 7 days).
func (r *PerformanceRepository) GetWeekly(ctx context.Context, filter ports.PerformanceFilter) ([]domain.StaffPerformanceSummary, error) {
	if filter.DateFrom == "" {
		now := time.Now().UTC()
		filter.DateFrom = now.AddDate(0, 0, -6).Format("2006-01-02")
		filter.DateTo = now.Format("2006-01-02")
	}
	return r.GetDaily(ctx, filter)
}

// GetMonthly returns monthly performance summaries.
func (r *PerformanceRepository) GetMonthly(ctx context.Context, filter ports.PerformanceFilter) ([]domain.StaffPerformanceSummary, error) {
	if filter.DateFrom == "" {
		now := time.Now().UTC()
		filter.DateFrom = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC).Format("2006-01-02")
		filter.DateTo = now.Format("2006-01-02")
	}
	return r.GetDaily(ctx, filter)
}

// GetByStaff returns raw daily performance records for a specific staff member.
func (r *PerformanceRepository) GetByStaff(ctx context.Context, staffID uuid.UUID, dateFrom, dateTo string) ([]domain.StaffPerformanceDaily, error) {
	query := `
		SELECT id, staff_id, business_date, invoice_count, customer_count, service_count, revenue, commission_amount, created_at, updated_at
		FROM staff_performance_daily
		WHERE staff_id = ? AND business_date BETWEEN ? AND ?
		ORDER BY business_date DESC`

	rows, err := r.db.QueryContext(ctx, query, staffID.String(), dateFrom, dateTo)
	if err != nil {
		return nil, apperror.Database("get staff performance", err)
	}
	defer rows.Close()

	var results []domain.StaffPerformanceDaily
	for rows.Next() {
		var p domain.StaffPerformanceDaily
		var idStr, staffIDStr, createdStr, updatedStr string
		err := rows.Scan(&idStr, &staffIDStr, &p.BusinessDate, &p.InvoiceCount, &p.CustomerCount, &p.ServiceCount, &p.Revenue, &p.CommissionAmount, &createdStr, &updatedStr)
		if err != nil {
			return nil, apperror.Database("scan staff performance", err)
		}
		p.ID, _ = uuid.Parse(idStr)
		p.StaffID, _ = uuid.Parse(staffIDStr)
		p.CreatedAt, _ = time.Parse(time.RFC3339, createdStr)
		p.UpdatedAt, _ = time.Parse(time.RFC3339, updatedStr)
		results = append(results, p)
	}
	return results, nil
}

// GetTopPerformers returns top performing staff by revenue.
func (r *PerformanceRepository) GetTopPerformers(ctx context.Context, dateFrom, dateTo string, limit int) ([]domain.StaffPerformanceSummary, error) {
	if limit <= 0 {
		limit = 10
	}

	query := `
		SELECT s.id, s.full_name,
			COALESCE(SUM(p.revenue), 0),
			COALESCE(SUM(p.customer_count), 0),
			COALESCE(SUM(p.invoice_count), 0),
			COALESCE(SUM(p.service_count), 0),
			COALESCE(SUM(p.commission_amount), 0)
		FROM staff s
		INNER JOIN staff_performance_daily p ON s.id = p.staff_id
		WHERE p.business_date BETWEEN ? AND ?
		GROUP BY s.id, s.full_name
		HAVING SUM(p.revenue) > 0
		ORDER BY SUM(p.revenue) DESC
		LIMIT ?`

	rows, err := r.db.QueryContext(ctx, query, dateFrom, dateTo, limit)
	if err != nil {
		return nil, apperror.Database("get top performers", err)
	}
	defer rows.Close()

	return r.scanSummaries(rows)
}

// GetRevenueTrend returns daily revenue totals for the date range.
func (r *PerformanceRepository) GetRevenueTrend(ctx context.Context, dateFrom, dateTo string) ([]domain.RevenueTrendPoint, error) {
	query := `
		SELECT business_date, COALESCE(SUM(revenue), 0)
		FROM staff_performance_daily
		WHERE business_date BETWEEN ? AND ?
		GROUP BY business_date
		ORDER BY business_date ASC`

	rows, err := r.db.QueryContext(ctx, query, dateFrom, dateTo)
	if err != nil {
		return nil, apperror.Database("get revenue trend", err)
	}
	defer rows.Close()

	var results []domain.RevenueTrendPoint
	for rows.Next() {
		var pt domain.RevenueTrendPoint
		if err := rows.Scan(&pt.Date, &pt.Revenue); err != nil {
			return nil, apperror.Database("scan revenue trend", err)
		}
		results = append(results, pt)
	}
	return results, nil
}

// GetStaffRevenueTrend returns daily revenue for a specific staff member.
func (r *PerformanceRepository) GetStaffRevenueTrend(ctx context.Context, staffID uuid.UUID, dateFrom, dateTo string) ([]domain.RevenueTrendPoint, error) {
	query := `
		SELECT business_date, revenue
		FROM staff_performance_daily
		WHERE staff_id = ? AND business_date BETWEEN ? AND ?
		ORDER BY business_date ASC`

	rows, err := r.db.QueryContext(ctx, query, staffID.String(), dateFrom, dateTo)
	if err != nil {
		return nil, apperror.Database("get staff revenue trend", err)
	}
	defer rows.Close()

	var results []domain.RevenueTrendPoint
	for rows.Next() {
		var pt domain.RevenueTrendPoint
		if err := rows.Scan(&pt.Date, &pt.Revenue); err != nil {
			return nil, apperror.Database("scan staff revenue trend", err)
		}
		results = append(results, pt)
	}
	return results, nil
}

func (r *PerformanceRepository) scanSummaries(rows *sql.Rows) ([]domain.StaffPerformanceSummary, error) {
	var results []domain.StaffPerformanceSummary
	rank := 0
	for rows.Next() {
		rank++
		var s domain.StaffPerformanceSummary
		var staffIDStr string
		err := rows.Scan(&staffIDStr, &s.StaffName, &s.Revenue, &s.CustomerCount, &s.InvoiceCount, &s.ServiceCount, &s.Commission)
		if err != nil {
			return nil, apperror.Database("scan performance summary", err)
		}
		s.StaffID, _ = uuid.Parse(staffIDStr)
		if s.InvoiceCount > 0 {
			s.AvgBill = s.Revenue / float64(s.InvoiceCount)
		}
		s.Rank = rank
		results = append(results, s)
	}
	return results, nil
}
