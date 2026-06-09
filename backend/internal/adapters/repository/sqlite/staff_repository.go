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

// StaffRepository is the SQLite implementation of ports.StaffRepository.
type StaffRepository struct {
	db  *sql.DB
	log *slog.Logger
}

// NewStaffRepository creates a new StaffRepository.
func NewStaffRepository(db *sql.DB, log *slog.Logger) *StaffRepository {
	return &StaffRepository{db: db, log: log}
}

// Create inserts a new staff record.
func (r *StaffRepository) Create(ctx context.Context, staff *domain.Staff) error {
	query := `
		INSERT INTO staff (id, staff_code, full_name, phone, email, gender, designation, joining_date, base_salary, commission_percentage, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := r.db.ExecContext(ctx, query,
		staff.ID.String(),
		staff.StaffCode,
		staff.FullName,
		staff.Phone,
		staff.Email,
		staff.Gender,
		staff.Designation,
		staff.JoiningDate.Format(time.RFC3339),
		staff.BaseSalary,
		staff.CommissionPercentage,
		staff.Status,
		staff.CreatedAt.Format(time.RFC3339),
		staff.UpdatedAt.Format(time.RFC3339),
	)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed: staff.phone") {
			return apperror.Conflict("phone number already exists")
		}
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return apperror.Conflict("staff record already exists")
		}
		return apperror.Database("create staff", err)
	}
	return nil
}

// GetByID retrieves a staff member by ID.
func (r *StaffRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Staff, error) {
	query := `
		SELECT id, staff_code, full_name, phone, email, gender, designation, joining_date, base_salary, commission_percentage, status, created_at, updated_at
		FROM staff
		WHERE id = ?`

	row := r.db.QueryRowContext(ctx, query, id.String())
	staff, err := r.scanStaff(row)
	if err == sql.ErrNoRows {
		return nil, apperror.NotFound("staff", id.String())
	}
	if err != nil {
		return nil, apperror.Database("get staff by id", err)
	}
	return staff, nil
}

// GetByPhone retrieves a staff member by phone number.
func (r *StaffRepository) GetByPhone(ctx context.Context, phone string) (*domain.Staff, error) {
	query := `
		SELECT id, staff_code, full_name, phone, email, gender, designation, joining_date, base_salary, commission_percentage, status, created_at, updated_at
		FROM staff
		WHERE phone = ?`

	row := r.db.QueryRowContext(ctx, query, phone)
	staff, err := r.scanStaff(row)
	if err == sql.ErrNoRows {
		return nil, apperror.NotFound("staff", phone)
	}
	if err != nil {
		return nil, apperror.Database("get staff by phone", err)
	}
	return staff, nil
}

// List returns staff members matching the filter.
func (r *StaffRepository) List(ctx context.Context, filter ports.StaffFilter) ([]domain.Staff, int, error) {
	var conditions []string
	var args []interface{}

	if filter.Status != "" {
		conditions = append(conditions, "status = ?")
		args = append(args, filter.Status)
	}
	if filter.Designation != "" {
		conditions = append(conditions, "designation = ?")
		args = append(args, filter.Designation)
	}
	if filter.Search != "" {
		conditions = append(conditions, "(full_name LIKE ? OR phone LIKE ? OR staff_code LIKE ?)")
		search := "%" + filter.Search + "%"
		args = append(args, search, search, search)
	}

	where := ""
	if len(conditions) > 0 {
		where = "WHERE " + strings.Join(conditions, " AND ")
	}

	// Count total
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM staff %s", where)
	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, apperror.Database("count staff", err)
	}

	// Fetch rows
	limit := filter.Limit
	if limit <= 0 {
		limit = 50
	}
	offset := filter.Offset
	if offset < 0 {
		offset = 0
	}

	dataQuery := fmt.Sprintf(`
		SELECT id, staff_code, full_name, phone, email, gender, designation, joining_date, base_salary, commission_percentage, status, created_at, updated_at
		FROM staff %s
		ORDER BY full_name ASC
		LIMIT ? OFFSET ?`, where)

	dataArgs := append(args, limit, offset)
	rows, err := r.db.QueryContext(ctx, dataQuery, dataArgs...)
	if err != nil {
		return nil, 0, apperror.Database("list staff", err)
	}
	defer rows.Close()

	var results []domain.Staff
	for rows.Next() {
		staff, err := r.scanStaffRows(rows)
		if err != nil {
			return nil, 0, apperror.Database("scan staff row", err)
		}
		results = append(results, *staff)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, apperror.Database("iterate staff rows", err)
	}

	return results, total, nil
}

// Update modifies an existing staff record.
func (r *StaffRepository) Update(ctx context.Context, staff *domain.Staff) error {
	staff.UpdatedAt = time.Now().UTC()

	query := `
		UPDATE staff
		SET full_name = ?, phone = ?, email = ?, gender = ?, designation = ?, joining_date = ?, base_salary = ?, commission_percentage = ?, status = ?, updated_at = ?
		WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query,
		staff.FullName,
		staff.Phone,
		staff.Email,
		staff.Gender,
		staff.Designation,
		staff.JoiningDate.Format(time.RFC3339),
		staff.BaseSalary,
		staff.CommissionPercentage,
		staff.Status,
		staff.UpdatedAt.Format(time.RFC3339),
		staff.ID.String(),
	)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed: staff.phone") {
			return apperror.Conflict("phone number already exists")
		}
		return apperror.Database("update staff", err)
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		return apperror.NotFound("staff", staff.ID.String())
	}
	return nil
}

// Delete removes a staff member (hard delete).
func (r *StaffRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM staff WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query, id.String())
	if err != nil {
		return apperror.Database("delete staff", err)
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		return apperror.NotFound("staff", id.String())
	}
	return nil
}

// CountByStatus returns count of total, active, and inactive staff.
func (r *StaffRepository) CountByStatus(ctx context.Context) (total, active, inactive int, err error) {
	query := `SELECT COUNT(*), COALESCE(SUM(CASE WHEN status = 'active' THEN 1 ELSE 0 END), 0), COALESCE(SUM(CASE WHEN status = 'inactive' THEN 1 ELSE 0 END), 0) FROM staff`
	err = r.db.QueryRowContext(ctx, query).Scan(&total, &active, &inactive)
	if err != nil {
		return 0, 0, 0, apperror.Database("count staff by status", err)
	}
	return total, active, inactive, nil
}

// --- scan helpers ---

func (r *StaffRepository) scanStaff(row *sql.Row) (*domain.Staff, error) {
	var s domain.Staff
	var id string
	var joiningDate, createdAt, updatedAt string

	err := row.Scan(&id, &s.StaffCode, &s.FullName, &s.Phone, &s.Email, &s.Gender, &s.Designation, &joiningDate, &s.BaseSalary, &s.CommissionPercentage, &s.Status, &createdAt, &updatedAt)
	if err != nil {
		return nil, err
	}

	return r.buildStaff(id, joiningDate, createdAt, updatedAt, &s)
}

func (r *StaffRepository) scanStaffRows(rows *sql.Rows) (*domain.Staff, error) {
	var s domain.Staff
	var id string
	var joiningDate, createdAt, updatedAt string

	err := rows.Scan(&id, &s.StaffCode, &s.FullName, &s.Phone, &s.Email, &s.Gender, &s.Designation, &joiningDate, &s.BaseSalary, &s.CommissionPercentage, &s.Status, &createdAt, &updatedAt)
	if err != nil {
		return nil, err
	}

	return r.buildStaff(id, joiningDate, createdAt, updatedAt, &s)
}

func (r *StaffRepository) buildStaff(id, joiningDate, createdAt, updatedAt string, s *domain.Staff) (*domain.Staff, error) {
	parsed, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("parse staff id: %w", err)
	}
	s.ID = parsed
	s.JoiningDate, _ = time.Parse(time.RFC3339, joiningDate)
	s.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	s.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)

	return s, nil
}
