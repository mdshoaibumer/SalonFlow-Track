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

// ServiceRepository is the SQLite implementation of ports.ServiceRepository.
type ServiceRepository struct {
	db  *sql.DB
	log *slog.Logger
}

// NewServiceRepository creates a new ServiceRepository.
func NewServiceRepository(db *sql.DB, log *slog.Logger) *ServiceRepository {
	return &ServiceRepository{db: db, log: log}
}

// Create inserts a new service record.
func (r *ServiceRepository) Create(ctx context.Context, service *domain.Service) error {
	query := `
		INSERT INTO services (id, service_code, name, category, description, duration_minutes, price, cost_price, commission_type, commission_value, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := r.db.ExecContext(ctx, query,
		service.ID.String(),
		service.ServiceCode,
		service.Name,
		service.Category,
		service.Description,
		service.DurationMinutes,
		service.Price,
		service.CostPrice,
		service.CommissionType,
		service.CommissionValue,
		service.Status,
		service.CreatedAt.Format(time.RFC3339),
		service.UpdatedAt.Format(time.RFC3339),
	)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return apperror.Conflict("service already exists")
		}
		return apperror.Database("create service", err)
	}
	return nil
}

// GetByID retrieves a service by ID.
func (r *ServiceRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Service, error) {
	query := `
		SELECT id, service_code, name, category, description, duration_minutes, price, cost_price, commission_type, commission_value, status, created_at, updated_at
		FROM services WHERE id = ?`

	row := r.db.QueryRowContext(ctx, query, id.String())
	svc, err := r.scanService(row)
	if err == sql.ErrNoRows {
		return nil, apperror.NotFound("service", id.String())
	}
	if err != nil {
		return nil, apperror.Database("get service by id", err)
	}
	return svc, nil
}

// GetByName retrieves a service by name.
func (r *ServiceRepository) GetByName(ctx context.Context, name string) (*domain.Service, error) {
	query := `
		SELECT id, service_code, name, category, description, duration_minutes, price, cost_price, commission_type, commission_value, status, created_at, updated_at
		FROM services WHERE LOWER(name) = LOWER(?)`

	row := r.db.QueryRowContext(ctx, query, name)
	svc, err := r.scanService(row)
	if err == sql.ErrNoRows {
		return nil, apperror.NotFound("service", name)
	}
	if err != nil {
		return nil, apperror.Database("get service by name", err)
	}
	return svc, nil
}

// List returns services matching the filter.
func (r *ServiceRepository) List(ctx context.Context, filter ports.ServiceFilter) ([]domain.Service, int, error) {
	var conditions []string
	var args []interface{}

	if filter.Status != "" {
		conditions = append(conditions, "status = ?")
		args = append(args, filter.Status)
	}
	if filter.Category != "" {
		conditions = append(conditions, "category = ?")
		args = append(args, filter.Category)
	}
	if filter.Search != "" {
		conditions = append(conditions, "(name LIKE ? OR service_code LIKE ? OR category LIKE ?)")
		search := "%" + filter.Search + "%"
		args = append(args, search, search, search)
	}

	where := ""
	if len(conditions) > 0 {
		where = "WHERE " + strings.Join(conditions, " AND ")
	}

	// Count total
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM services %s", where)
	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, apperror.Database("count services", err)
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
		SELECT id, service_code, name, category, description, duration_minutes, price, cost_price, commission_type, commission_value, status, created_at, updated_at
		FROM services %s
		ORDER BY name ASC
		LIMIT ? OFFSET ?`, where)

	dataArgs := append(args, limit, offset)
	rows, err := r.db.QueryContext(ctx, dataQuery, dataArgs...)
	if err != nil {
		return nil, 0, apperror.Database("list services", err)
	}
	defer rows.Close()

	var results []domain.Service
	for rows.Next() {
		svc, err := r.scanServiceRows(rows)
		if err != nil {
			return nil, 0, apperror.Database("scan service row", err)
		}
		results = append(results, *svc)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, apperror.Database("iterate service rows", err)
	}

	return results, total, nil
}

// Update modifies an existing service record.
func (r *ServiceRepository) Update(ctx context.Context, service *domain.Service) error {
	service.UpdatedAt = time.Now().UTC()

	query := `
		UPDATE services
		SET name = ?, category = ?, description = ?, duration_minutes = ?, price = ?, cost_price = ?, commission_type = ?, commission_value = ?, status = ?, updated_at = ?
		WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query,
		service.Name,
		service.Category,
		service.Description,
		service.DurationMinutes,
		service.Price,
		service.CostPrice,
		service.CommissionType,
		service.CommissionValue,
		service.Status,
		service.UpdatedAt.Format(time.RFC3339),
		service.ID.String(),
	)
	if err != nil {
		return apperror.Database("update service", err)
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		return apperror.NotFound("service", service.ID.String())
	}
	return nil
}

// Delete removes a service (hard delete).
func (r *ServiceRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM services WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query, id.String())
	if err != nil {
		return apperror.Database("delete service", err)
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		return apperror.NotFound("service", id.String())
	}
	return nil
}

// CountByStatus returns count of total, active, and inactive services.
func (r *ServiceRepository) CountByStatus(ctx context.Context) (total, active, inactive int, err error) {
	query := `SELECT COUNT(*), COALESCE(SUM(CASE WHEN status = 'active' THEN 1 ELSE 0 END), 0), COALESCE(SUM(CASE WHEN status = 'inactive' THEN 1 ELSE 0 END), 0) FROM services`
	err = r.db.QueryRowContext(ctx, query).Scan(&total, &active, &inactive)
	if err != nil {
		return 0, 0, 0, apperror.Database("count services by status", err)
	}
	return total, active, inactive, nil
}

// --- scan helpers ---

func (r *ServiceRepository) scanService(row *sql.Row) (*domain.Service, error) {
	var s domain.Service
	var id string
	var createdAt, updatedAt string

	err := row.Scan(&id, &s.ServiceCode, &s.Name, &s.Category, &s.Description, &s.DurationMinutes, &s.Price, &s.CostPrice, &s.CommissionType, &s.CommissionValue, &s.Status, &createdAt, &updatedAt)
	if err != nil {
		return nil, err
	}

	return r.buildService(id, createdAt, updatedAt, &s)
}

func (r *ServiceRepository) scanServiceRows(rows *sql.Rows) (*domain.Service, error) {
	var s domain.Service
	var id string
	var createdAt, updatedAt string

	err := rows.Scan(&id, &s.ServiceCode, &s.Name, &s.Category, &s.Description, &s.DurationMinutes, &s.Price, &s.CostPrice, &s.CommissionType, &s.CommissionValue, &s.Status, &createdAt, &updatedAt)
	if err != nil {
		return nil, err
	}

	return r.buildService(id, createdAt, updatedAt, &s)
}

func (r *ServiceRepository) buildService(id, createdAt, updatedAt string, s *domain.Service) (*domain.Service, error) {
	parsed, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	s.ID = parsed

	if t, err := time.Parse(time.RFC3339, createdAt); err == nil {
		s.CreatedAt = t
	}
	if t, err := time.Parse(time.RFC3339, updatedAt); err == nil {
		s.UpdatedAt = t
	}

	return s, nil
}
