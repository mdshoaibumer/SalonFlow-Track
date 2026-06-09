package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/salonflow/salonflow-track/internal/core/domain"
	"github.com/salonflow/salonflow-track/internal/core/ports"
	"github.com/salonflow/salonflow-track/pkg/apperror"
)

// SalaryRepository is the SQLite implementation of ports.SalaryRepository.
type SalaryRepository struct {
	db  *sql.DB
	log *slog.Logger
}

// NewSalaryRepository creates a new SalaryRepository.
func NewSalaryRepository(db *sql.DB, log *slog.Logger) *SalaryRepository {
	return &SalaryRepository{db: db, log: log}
}

// --- Salary Cycles ---

// CreateCycle inserts a new salary cycle.
func (r *SalaryRepository) CreateCycle(ctx context.Context, cycle *domain.SalaryCycle) error {
	query := `INSERT INTO salary_cycles (id, month, year, status, generated_at, generated_by, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := r.db.ExecContext(ctx, query,
		cycle.ID.String(), cycle.Month, cycle.Year, cycle.Status,
		cycle.GeneratedAt, cycle.GeneratedBy,
		cycle.CreatedAt.Format(time.RFC3339), cycle.UpdatedAt.Format(time.RFC3339),
	)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return apperror.Conflict(domain.ErrSalaryCycleExists.Error())
		}
		return apperror.Database("create salary cycle", err)
	}
	return nil
}

// GetCycleByID retrieves a salary cycle by ID.
func (r *SalaryRepository) GetCycleByID(ctx context.Context, id uuid.UUID) (*domain.SalaryCycle, error) {
	query := `SELECT id, month, year, status, COALESCE(generated_at, ''), COALESCE(generated_by, ''), created_at, updated_at
		FROM salary_cycles WHERE id = ?`

	var c domain.SalaryCycle
	var idStr, createdStr, updatedStr string
	err := r.db.QueryRowContext(ctx, query, id.String()).Scan(
		&idStr, &c.Month, &c.Year, &c.Status, &c.GeneratedAt, &c.GeneratedBy, &createdStr, &updatedStr,
	)
	if err == sql.ErrNoRows {
		return nil, apperror.NotFound("salary_cycle", id.String())
	}
	if err != nil {
		return nil, apperror.Database("get salary cycle", err)
	}
	c.ID, _ = uuid.Parse(idStr)
	c.CreatedAt, _ = time.Parse(time.RFC3339, createdStr)
	c.UpdatedAt, _ = time.Parse(time.RFC3339, updatedStr)
	return &c, nil
}

// GetCycleByMonthYear retrieves a salary cycle by month and year.
func (r *SalaryRepository) GetCycleByMonthYear(ctx context.Context, month, year int) (*domain.SalaryCycle, error) {
	query := `SELECT id, month, year, status, COALESCE(generated_at, ''), COALESCE(generated_by, ''), created_at, updated_at
		FROM salary_cycles WHERE month = ? AND year = ?`

	var c domain.SalaryCycle
	var idStr, createdStr, updatedStr string
	err := r.db.QueryRowContext(ctx, query, month, year).Scan(
		&idStr, &c.Month, &c.Year, &c.Status, &c.GeneratedAt, &c.GeneratedBy, &createdStr, &updatedStr,
	)
	if err == sql.ErrNoRows {
		return nil, apperror.NotFound("salary_cycle", fmt.Sprintf("%d-%d", year, month))
	}
	if err != nil {
		return nil, apperror.Database("get salary cycle by month", err)
	}
	c.ID, _ = uuid.Parse(idStr)
	c.CreatedAt, _ = time.Parse(time.RFC3339, createdStr)
	c.UpdatedAt, _ = time.Parse(time.RFC3339, updatedStr)
	return &c, nil
}

// ListCycles returns salary cycles with filtering.
func (r *SalaryRepository) ListCycles(ctx context.Context, filter ports.SalaryCycleFilter) ([]domain.SalaryCycle, int, error) {
	where := []string{"1=1"}
	args := []interface{}{}

	if filter.Year > 0 {
		where = append(where, "year = ?")
		args = append(args, filter.Year)
	}
	if filter.Status != "" {
		where = append(where, "status = ?")
		args = append(args, filter.Status)
	}

	whereClause := strings.Join(where, " AND ")

	var total int
	countQ := fmt.Sprintf("SELECT COUNT(*) FROM salary_cycles WHERE %s", whereClause)
	if err := r.db.QueryRowContext(ctx, countQ, args...).Scan(&total); err != nil {
		return nil, 0, apperror.Database("count salary cycles", err)
	}

	limit := filter.Limit
	if limit <= 0 {
		limit = 50
	}
	query := fmt.Sprintf(`SELECT id, month, year, status, COALESCE(generated_at, ''), COALESCE(generated_by, ''), created_at, updated_at
		FROM salary_cycles WHERE %s ORDER BY year DESC, month DESC LIMIT ? OFFSET ?`, whereClause)

	args = append(args, limit, filter.Offset)
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, apperror.Database("list salary cycles", err)
	}
	defer rows.Close()

	var cycles []domain.SalaryCycle
	for rows.Next() {
		var c domain.SalaryCycle
		var idStr, createdStr, updatedStr string
		err := rows.Scan(&idStr, &c.Month, &c.Year, &c.Status, &c.GeneratedAt, &c.GeneratedBy, &createdStr, &updatedStr)
		if err != nil {
			return nil, 0, apperror.Database("scan salary cycle", err)
		}
		c.ID, _ = uuid.Parse(idStr)
		c.CreatedAt, _ = time.Parse(time.RFC3339, createdStr)
		c.UpdatedAt, _ = time.Parse(time.RFC3339, updatedStr)
		cycles = append(cycles, c)
	}
	return cycles, total, nil
}

// UpdateCycleStatus updates a salary cycle's status.
func (r *SalaryRepository) UpdateCycleStatus(ctx context.Context, id uuid.UUID, status, generatedAt, generatedBy string) error {
	query := `UPDATE salary_cycles SET status = ?, generated_at = ?, generated_by = ?, updated_at = ? WHERE id = ?`
	result, err := r.db.ExecContext(ctx, query, status, generatedAt, generatedBy, time.Now().UTC().Format(time.RFC3339), id.String())
	if err != nil {
		return apperror.Database("update salary cycle status", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return apperror.NotFound("salary_cycle", id.String())
	}
	return nil
}

// --- Salary Records ---

// CreateRecord inserts a new salary record.
func (r *SalaryRepository) CreateRecord(ctx context.Context, record *domain.SalaryRecord) error {
	query := `INSERT INTO salary_records (id, salary_cycle_id, staff_id, base_salary, commission_amount, bonus_amount, advance_amount, deduction_amount, gross_salary, net_salary, payment_status, payment_date, notes, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := r.db.ExecContext(ctx, query,
		record.ID.String(), record.SalaryCycleID.String(), record.StaffID.String(),
		record.BaseSalary, record.CommissionAmount, record.BonusAmount,
		record.AdvanceAmount, record.DeductionAmount, record.GrossSalary, record.NetSalary,
		record.PaymentStatus, record.PaymentDate, record.Notes,
		record.CreatedAt.Format(time.RFC3339), record.UpdatedAt.Format(time.RFC3339),
	)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return apperror.Conflict("salary record already exists for this staff in this cycle")
		}
		return apperror.Database("create salary record", err)
	}
	return nil
}

// GetRecordByID retrieves a salary record by ID.
func (r *SalaryRepository) GetRecordByID(ctx context.Context, id uuid.UUID) (*domain.SalaryRecord, error) {
	query := `SELECT sr.id, sr.salary_cycle_id, sr.staff_id, s.full_name, sr.base_salary, sr.commission_amount, sr.bonus_amount, sr.advance_amount, sr.deduction_amount, sr.gross_salary, sr.net_salary, sr.payment_status, COALESCE(sr.payment_date, ''), sr.notes, sr.created_at, sr.updated_at
		FROM salary_records sr JOIN staff s ON sr.staff_id = s.id WHERE sr.id = ?`

	var rec domain.SalaryRecord
	var idStr, cycleStr, staffStr, createdStr, updatedStr string
	err := r.db.QueryRowContext(ctx, query, id.String()).Scan(
		&idStr, &cycleStr, &staffStr, &rec.StaffName,
		&rec.BaseSalary, &rec.CommissionAmount, &rec.BonusAmount,
		&rec.AdvanceAmount, &rec.DeductionAmount, &rec.GrossSalary, &rec.NetSalary,
		&rec.PaymentStatus, &rec.PaymentDate, &rec.Notes,
		&createdStr, &updatedStr,
	)
	if err == sql.ErrNoRows {
		return nil, apperror.NotFound("salary_record", id.String())
	}
	if err != nil {
		return nil, apperror.Database("get salary record", err)
	}
	rec.ID, _ = uuid.Parse(idStr)
	rec.SalaryCycleID, _ = uuid.Parse(cycleStr)
	rec.StaffID, _ = uuid.Parse(staffStr)
	rec.CreatedAt, _ = time.Parse(time.RFC3339, createdStr)
	rec.UpdatedAt, _ = time.Parse(time.RFC3339, updatedStr)
	return &rec, nil
}

// ListRecordsByCycle returns all salary records for a cycle.
func (r *SalaryRepository) ListRecordsByCycle(ctx context.Context, cycleID uuid.UUID) ([]domain.SalaryRecord, error) {
	query := `SELECT sr.id, sr.salary_cycle_id, sr.staff_id, s.full_name, sr.base_salary, sr.commission_amount, sr.bonus_amount, sr.advance_amount, sr.deduction_amount, sr.gross_salary, sr.net_salary, sr.payment_status, COALESCE(sr.payment_date, ''), sr.notes, sr.created_at, sr.updated_at
		FROM salary_records sr JOIN staff s ON sr.staff_id = s.id WHERE sr.salary_cycle_id = ? ORDER BY sr.net_salary DESC`

	rows, err := r.db.QueryContext(ctx, query, cycleID.String())
	if err != nil {
		return nil, apperror.Database("list salary records", err)
	}
	defer rows.Close()

	var records []domain.SalaryRecord
	for rows.Next() {
		var rec domain.SalaryRecord
		var idStr, cycleStr, staffStr, createdStr, updatedStr string
		err := rows.Scan(
			&idStr, &cycleStr, &staffStr, &rec.StaffName,
			&rec.BaseSalary, &rec.CommissionAmount, &rec.BonusAmount,
			&rec.AdvanceAmount, &rec.DeductionAmount, &rec.GrossSalary, &rec.NetSalary,
			&rec.PaymentStatus, &rec.PaymentDate, &rec.Notes,
			&createdStr, &updatedStr,
		)
		if err != nil {
			return nil, apperror.Database("scan salary record", err)
		}
		rec.ID, _ = uuid.Parse(idStr)
		rec.SalaryCycleID, _ = uuid.Parse(cycleStr)
		rec.StaffID, _ = uuid.Parse(staffStr)
		rec.CreatedAt, _ = time.Parse(time.RFC3339, createdStr)
		rec.UpdatedAt, _ = time.Parse(time.RFC3339, updatedStr)
		records = append(records, rec)
	}
	return records, nil
}

// UpdateRecordPayment updates a salary record's payment status.
func (r *SalaryRepository) UpdateRecordPayment(ctx context.Context, id uuid.UUID, status, paymentDate string) error {
	query := `UPDATE salary_records SET payment_status = ?, payment_date = ?, updated_at = ? WHERE id = ?`
	result, err := r.db.ExecContext(ctx, query, status, paymentDate, time.Now().UTC().Format(time.RFC3339), id.String())
	if err != nil {
		return apperror.Database("update salary payment", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return apperror.NotFound("salary_record", id.String())
	}
	return nil
}

// --- Advances ---

// CreateAdvance inserts a new advance.
func (r *SalaryRepository) CreateAdvance(ctx context.Context, advance *domain.Advance) error {
	query := `INSERT INTO advances (id, staff_id, amount, advance_date, reason, recovered_amount, remaining_amount, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := r.db.ExecContext(ctx, query,
		advance.ID.String(), advance.StaffID.String(),
		advance.Amount, advance.AdvanceDate, advance.Reason,
		advance.RecoveredAmount, advance.RemainingAmount, advance.Status,
		advance.CreatedAt.Format(time.RFC3339), advance.UpdatedAt.Format(time.RFC3339),
	)
	if err != nil {
		return apperror.Database("create advance", err)
	}
	return nil
}

// GetAdvanceByID retrieves an advance by ID.
func (r *SalaryRepository) GetAdvanceByID(ctx context.Context, id uuid.UUID) (*domain.Advance, error) {
	query := `SELECT a.id, a.staff_id, s.full_name, a.amount, a.advance_date, a.reason, a.recovered_amount, a.remaining_amount, a.status, a.created_at, a.updated_at
		FROM advances a JOIN staff s ON a.staff_id = s.id WHERE a.id = ?`

	var adv domain.Advance
	var idStr, staffStr, createdStr, updatedStr string
	err := r.db.QueryRowContext(ctx, query, id.String()).Scan(
		&idStr, &staffStr, &adv.StaffName,
		&adv.Amount, &adv.AdvanceDate, &adv.Reason,
		&adv.RecoveredAmount, &adv.RemainingAmount, &adv.Status,
		&createdStr, &updatedStr,
	)
	if err == sql.ErrNoRows {
		return nil, apperror.NotFound("advance", id.String())
	}
	if err != nil {
		return nil, apperror.Database("get advance", err)
	}
	adv.ID, _ = uuid.Parse(idStr)
	adv.StaffID, _ = uuid.Parse(staffStr)
	adv.CreatedAt, _ = time.Parse(time.RFC3339, createdStr)
	adv.UpdatedAt, _ = time.Parse(time.RFC3339, updatedStr)
	return &adv, nil
}

// ListAdvances returns advances with filtering.
func (r *SalaryRepository) ListAdvances(ctx context.Context, filter ports.AdvanceFilter) ([]domain.Advance, int, error) {
	where := []string{"1=1"}
	args := []interface{}{}

	if filter.StaffID != "" {
		where = append(where, "a.staff_id = ?")
		args = append(args, filter.StaffID)
	}
	if filter.Status != "" {
		where = append(where, "a.status = ?")
		args = append(args, filter.Status)
	}

	whereClause := strings.Join(where, " AND ")

	var total int
	countQ := fmt.Sprintf("SELECT COUNT(*) FROM advances a WHERE %s", whereClause)
	if err := r.db.QueryRowContext(ctx, countQ, args...).Scan(&total); err != nil {
		return nil, 0, apperror.Database("count advances", err)
	}

	limit := filter.Limit
	if limit <= 0 {
		limit = 50
	}
	query := fmt.Sprintf(`SELECT a.id, a.staff_id, s.full_name, a.amount, a.advance_date, a.reason, a.recovered_amount, a.remaining_amount, a.status, a.created_at, a.updated_at
		FROM advances a JOIN staff s ON a.staff_id = s.id WHERE %s ORDER BY a.advance_date DESC LIMIT ? OFFSET ?`, whereClause)

	args = append(args, limit, filter.Offset)
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, apperror.Database("list advances", err)
	}
	defer rows.Close()

	var advances []domain.Advance
	for rows.Next() {
		var adv domain.Advance
		var idStr, staffStr, createdStr, updatedStr string
		err := rows.Scan(
			&idStr, &staffStr, &adv.StaffName,
			&adv.Amount, &adv.AdvanceDate, &adv.Reason,
			&adv.RecoveredAmount, &adv.RemainingAmount, &adv.Status,
			&createdStr, &updatedStr,
		)
		if err != nil {
			return nil, 0, apperror.Database("scan advance", err)
		}
		adv.ID, _ = uuid.Parse(idStr)
		adv.StaffID, _ = uuid.Parse(staffStr)
		adv.CreatedAt, _ = time.Parse(time.RFC3339, createdStr)
		adv.UpdatedAt, _ = time.Parse(time.RFC3339, updatedStr)
		advances = append(advances, adv)
	}
	return advances, total, nil
}

// UpdateAdvance updates an advance record.
func (r *SalaryRepository) UpdateAdvance(ctx context.Context, advance *domain.Advance) error {
	query := `UPDATE advances SET recovered_amount = ?, remaining_amount = ?, status = ?, updated_at = ? WHERE id = ?`
	result, err := r.db.ExecContext(ctx, query,
		advance.RecoveredAmount, advance.RemainingAmount, advance.Status,
		time.Now().UTC().Format(time.RFC3339), advance.ID.String(),
	)
	if err != nil {
		return apperror.Database("update advance", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return apperror.NotFound("advance", advance.ID.String())
	}
	return nil
}

// GetPendingAdvances returns approved/recovering advances for a staff member.
func (r *SalaryRepository) GetPendingAdvances(ctx context.Context, staffID uuid.UUID) ([]domain.Advance, error) {
	query := `SELECT id, staff_id, amount, advance_date, reason, recovered_amount, remaining_amount, status, created_at, updated_at
		FROM advances WHERE staff_id = ? AND status IN ('approved', 'recovering') ORDER BY advance_date ASC`

	rows, err := r.db.QueryContext(ctx, query, staffID.String())
	if err != nil {
		return nil, apperror.Database("get pending advances", err)
	}
	defer rows.Close()

	var advances []domain.Advance
	for rows.Next() {
		var adv domain.Advance
		var idStr, staffStr, createdStr, updatedStr string
		err := rows.Scan(&idStr, &staffStr, &adv.Amount, &adv.AdvanceDate, &adv.Reason,
			&adv.RecoveredAmount, &adv.RemainingAmount, &adv.Status, &createdStr, &updatedStr)
		if err != nil {
			return nil, apperror.Database("scan pending advance", err)
		}
		adv.ID, _ = uuid.Parse(idStr)
		adv.StaffID, _ = uuid.Parse(staffStr)
		adv.CreatedAt, _ = time.Parse(time.RFC3339, createdStr)
		adv.UpdatedAt, _ = time.Parse(time.RFC3339, updatedStr)
		advances = append(advances, adv)
	}
	return advances, nil
}

// GetTotalOutstandingAdvances returns the total remaining advance amount for a staff member.
func (r *SalaryRepository) GetTotalOutstandingAdvances(ctx context.Context, staffID uuid.UUID) (float64, error) {
	query := `SELECT COALESCE(SUM(remaining_amount), 0) FROM advances WHERE staff_id = ? AND status IN ('approved', 'recovering')`
	var total float64
	err := r.db.QueryRowContext(ctx, query, staffID.String()).Scan(&total)
	if err != nil {
		return 0, apperror.Database("get outstanding advances", err)
	}
	return total, nil
}

// GetPayrollStats returns payroll dashboard statistics.
func (r *SalaryRepository) GetPayrollStats(ctx context.Context, month, year int) (*ports.PayrollStats, error) {
	stats := &ports.PayrollStats{}

	// Get totals from the specific cycle
	cycleQ := `SELECT COALESCE(SUM(sr.net_salary), 0),
		COALESCE(SUM(CASE WHEN sr.payment_status = 'pending' THEN 1 ELSE 0 END), 0),
		COALESCE(SUM(CASE WHEN sr.payment_status = 'paid' THEN 1 ELSE 0 END), 0)
		FROM salary_records sr
		JOIN salary_cycles sc ON sr.salary_cycle_id = sc.id
		WHERE sc.month = ? AND sc.year = ?`

	err := r.db.QueryRowContext(ctx, cycleQ, month, year).Scan(
		&stats.TotalPayroll, &stats.PendingPayments, &stats.PaidSalaries,
	)
	if err != nil && err != sql.ErrNoRows {
		return nil, apperror.Database("get payroll stats", err)
	}

	// Outstanding advances
	advQ := `SELECT COALESCE(SUM(remaining_amount), 0) FROM advances WHERE status IN ('approved', 'recovering')`
	_ = r.db.QueryRowContext(ctx, advQ).Scan(&stats.OutstandingAdvances)

	return stats, nil
}
