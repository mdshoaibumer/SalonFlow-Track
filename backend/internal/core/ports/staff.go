package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/salonflow/salonflow-track/internal/core/domain"
)

// StaffRepository defines persistence operations for Staff.
type StaffRepository interface {
	Create(ctx context.Context, staff *domain.Staff) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Staff, error)
	GetByPhone(ctx context.Context, phone string) (*domain.Staff, error)
	List(ctx context.Context, filter StaffFilter) ([]domain.Staff, int, error)
	Update(ctx context.Context, staff *domain.Staff) error
	Delete(ctx context.Context, id uuid.UUID) error
	CountByStatus(ctx context.Context) (total, active, inactive int, err error)
}

// StaffFilter holds query parameters for listing staff.
type StaffFilter struct {
	Status      string
	Designation string
	Search      string
	Limit       int
	Offset      int
}
