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

// CustomerRepository is the SQLite implementation of ports.CustomerRepository.
type CustomerRepository struct {
	db  *sql.DB
	log *slog.Logger
}

// NewCustomerRepository creates a new CustomerRepository.
func NewCustomerRepository(db *sql.DB, log *slog.Logger) *CustomerRepository {
	return &CustomerRepository{db: db, log: log}
}

// Create inserts a new customer record.
func (r *CustomerRepository) Create(ctx context.Context, customer *domain.Customer) error {
	query := `
		INSERT INTO customers (id, customer_code, full_name, phone, email, gender, date_of_birth, anniversary_date, address, notes, total_visits, total_spent, last_visit_date, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := r.db.ExecContext(ctx, query,
		customer.ID.String(),
		customer.CustomerCode,
		customer.FullName,
		customer.Phone,
		customer.Email,
		customer.Gender,
		timeToNullString(customer.DateOfBirth),
		timeToNullString(customer.AnniversaryDate),
		customer.Address,
		customer.Notes,
		customer.TotalVisits,
		customer.TotalSpent,
		timeToNullString(customer.LastVisitDate),
		customer.Status,
		customer.CreatedAt.Format(time.RFC3339),
		customer.UpdatedAt.Format(time.RFC3339),
	)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed: customers.phone") {
			return apperror.Conflict("phone number already exists")
		}
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return apperror.Conflict("customer already exists")
		}
		return apperror.Database("create customer", err)
	}
	return nil
}

// GetByID retrieves a customer by ID.
func (r *CustomerRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Customer, error) {
	query := `
		SELECT id, customer_code, full_name, phone, email, gender, date_of_birth, anniversary_date, address, notes, total_visits, total_spent, last_visit_date, status, created_at, updated_at
		FROM customers WHERE id = ?`

	row := r.db.QueryRowContext(ctx, query, id.String())
	cust, err := r.scanCustomer(row)
	if err == sql.ErrNoRows {
		return nil, apperror.NotFound("customer", id.String())
	}
	if err != nil {
		return nil, apperror.Database("get customer by id", err)
	}
	return cust, nil
}

// GetByPhone retrieves a customer by phone number.
func (r *CustomerRepository) GetByPhone(ctx context.Context, phone string) (*domain.Customer, error) {
	query := `
		SELECT id, customer_code, full_name, phone, email, gender, date_of_birth, anniversary_date, address, notes, total_visits, total_spent, last_visit_date, status, created_at, updated_at
		FROM customers WHERE phone = ?`

	row := r.db.QueryRowContext(ctx, query, phone)
	cust, err := r.scanCustomer(row)
	if err == sql.ErrNoRows {
		return nil, apperror.NotFound("customer", phone)
	}
	if err != nil {
		return nil, apperror.Database("get customer by phone", err)
	}
	return cust, nil
}

// List returns customers matching the filter.
func (r *CustomerRepository) List(ctx context.Context, filter ports.CustomerFilter) ([]domain.Customer, int, error) {
	var conditions []string
	var args []interface{}

	if filter.Status != "" {
		conditions = append(conditions, "status = ?")
		args = append(args, filter.Status)
	}
	if filter.Search != "" {
		conditions = append(conditions, "(full_name LIKE ? OR phone LIKE ? OR customer_code LIKE ?)")
		search := "%" + filter.Search + "%"
		args = append(args, search, search, search)
	}

	where := ""
	if len(conditions) > 0 {
		where = "WHERE " + strings.Join(conditions, " AND ")
	}

	// Count total
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM customers %s", where)
	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, apperror.Database("count customers", err)
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
		SELECT id, customer_code, full_name, phone, email, gender, date_of_birth, anniversary_date, address, notes, total_visits, total_spent, last_visit_date, status, created_at, updated_at
		FROM customers %s
		ORDER BY full_name ASC
		LIMIT ? OFFSET ?`, where)

	dataArgs := append(args, limit, offset)
	rows, err := r.db.QueryContext(ctx, dataQuery, dataArgs...)
	if err != nil {
		return nil, 0, apperror.Database("list customers", err)
	}
	defer rows.Close()

	var results []domain.Customer
	for rows.Next() {
		cust, err := r.scanCustomerRows(rows)
		if err != nil {
			return nil, 0, apperror.Database("scan customer row", err)
		}
		results = append(results, *cust)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, apperror.Database("iterate customer rows", err)
	}

	return results, total, nil
}

// Update modifies an existing customer record.
func (r *CustomerRepository) Update(ctx context.Context, customer *domain.Customer) error {
	customer.UpdatedAt = time.Now().UTC()

	query := `
		UPDATE customers
		SET full_name = ?, phone = ?, email = ?, gender = ?, date_of_birth = ?, anniversary_date = ?, address = ?, notes = ?, total_visits = ?, total_spent = ?, last_visit_date = ?, status = ?, updated_at = ?
		WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query,
		customer.FullName,
		customer.Phone,
		customer.Email,
		customer.Gender,
		timeToNullString(customer.DateOfBirth),
		timeToNullString(customer.AnniversaryDate),
		customer.Address,
		customer.Notes,
		customer.TotalVisits,
		customer.TotalSpent,
		timeToNullString(customer.LastVisitDate),
		customer.Status,
		customer.UpdatedAt.Format(time.RFC3339),
		customer.ID.String(),
	)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed: customers.phone") {
			return apperror.Conflict("phone number already exists")
		}
		return apperror.Database("update customer", err)
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		return apperror.NotFound("customer", customer.ID.String())
	}
	return nil
}

// Delete removes a customer (hard delete).
func (r *CustomerRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM customers WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query, id.String())
	if err != nil {
		return apperror.Database("delete customer", err)
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		return apperror.NotFound("customer", id.String())
	}
	return nil
}

// CountByStatus returns count of total, active, and inactive customers.
func (r *CustomerRepository) CountByStatus(ctx context.Context) (total, active, inactive int, err error) {
	query := `SELECT COUNT(*), COALESCE(SUM(CASE WHEN status = 'active' THEN 1 ELSE 0 END), 0), COALESCE(SUM(CASE WHEN status = 'inactive' THEN 1 ELSE 0 END), 0) FROM customers`
	err = r.db.QueryRowContext(ctx, query).Scan(&total, &active, &inactive)
	if err != nil {
		return 0, 0, 0, apperror.Database("count customers by status", err)
	}
	return total, active, inactive, nil
}

// CountNewThisMonth returns count of customers created this month.
func (r *CustomerRepository) CountNewThisMonth(ctx context.Context) (int, error) {
	now := time.Now().UTC()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC).Format(time.RFC3339)
	query := `SELECT COUNT(*) FROM customers WHERE created_at >= ?`
	var count int
	if err := r.db.QueryRowContext(ctx, query, startOfMonth).Scan(&count); err != nil {
		return 0, apperror.Database("count new customers", err)
	}
	return count, nil
}

// CountBirthdayToday returns count of customers with birthday today.
func (r *CustomerRepository) CountBirthdayToday(ctx context.Context) (int, error) {
	now := time.Now().UTC()
	monthDay := fmt.Sprintf("-%02d-%02d", now.Month(), now.Day())
	query := `SELECT COUNT(*) FROM customers WHERE date_of_birth LIKE ?`
	var count int
	if err := r.db.QueryRowContext(ctx, query, "%"+monthDay+"%").Scan(&count); err != nil {
		return 0, apperror.Database("count birthdays today", err)
	}
	return count, nil
}

// --- helpers ---

func timeToNullString(t *time.Time) interface{} {
	if t == nil {
		return nil
	}
	return t.Format(time.RFC3339)
}

func (r *CustomerRepository) scanCustomer(row *sql.Row) (*domain.Customer, error) {
	var c domain.Customer
	var id string
	var dob, anniversary, lastVisit sql.NullString
	var createdAt, updatedAt string

	err := row.Scan(&id, &c.CustomerCode, &c.FullName, &c.Phone, &c.Email, &c.Gender, &dob, &anniversary, &c.Address, &c.Notes, &c.TotalVisits, &c.TotalSpent, &lastVisit, &c.Status, &createdAt, &updatedAt)
	if err != nil {
		return nil, err
	}

	return r.buildCustomer(id, dob, anniversary, lastVisit, createdAt, updatedAt, &c)
}

func (r *CustomerRepository) scanCustomerRows(rows *sql.Rows) (*domain.Customer, error) {
	var c domain.Customer
	var id string
	var dob, anniversary, lastVisit sql.NullString
	var createdAt, updatedAt string

	err := rows.Scan(&id, &c.CustomerCode, &c.FullName, &c.Phone, &c.Email, &c.Gender, &dob, &anniversary, &c.Address, &c.Notes, &c.TotalVisits, &c.TotalSpent, &lastVisit, &c.Status, &createdAt, &updatedAt)
	if err != nil {
		return nil, err
	}

	return r.buildCustomer(id, dob, anniversary, lastVisit, createdAt, updatedAt, &c)
}

func (r *CustomerRepository) buildCustomer(id string, dob, anniversary, lastVisit sql.NullString, createdAt, updatedAt string, c *domain.Customer) (*domain.Customer, error) {
	parsed, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	c.ID = parsed

	if dob.Valid {
		if t, err := time.Parse(time.RFC3339, dob.String); err == nil {
			c.DateOfBirth = &t
		} else if t, err := time.Parse("2006-01-02", dob.String); err == nil {
			c.DateOfBirth = &t
		}
	}
	if anniversary.Valid {
		if t, err := time.Parse(time.RFC3339, anniversary.String); err == nil {
			c.AnniversaryDate = &t
		} else if t, err := time.Parse("2006-01-02", anniversary.String); err == nil {
			c.AnniversaryDate = &t
		}
	}
	if lastVisit.Valid {
		if t, err := time.Parse(time.RFC3339, lastVisit.String); err == nil {
			c.LastVisitDate = &t
		} else if t, err := time.Parse("2006-01-02", lastVisit.String); err == nil {
			c.LastVisitDate = &t
		}
	}

	if t, err := time.Parse(time.RFC3339, createdAt); err == nil {
		c.CreatedAt = t
	}
	if t, err := time.Parse(time.RFC3339, updatedAt); err == nil {
		c.UpdatedAt = t
	}

	return c, nil
}
