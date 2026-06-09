package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/salonflow/salonflow-track/internal/core/domain"
)

// CustomerRepository defines persistence operations for Customer.
type CustomerRepository interface {
	Create(ctx context.Context, customer *domain.Customer) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Customer, error)
	GetByPhone(ctx context.Context, phone string) (*domain.Customer, error)
	List(ctx context.Context, filter CustomerFilter) ([]domain.Customer, int, error)
	Update(ctx context.Context, customer *domain.Customer) error
	Delete(ctx context.Context, id uuid.UUID) error
	CountByStatus(ctx context.Context) (total, active, inactive int, err error)
	CountNewThisMonth(ctx context.Context) (int, error)
	CountBirthdayToday(ctx context.Context) (int, error)
}

// CustomerFilter holds query parameters for listing customers.
type CustomerFilter struct {
	Status string
	Search string
	Limit  int
	Offset int
}
