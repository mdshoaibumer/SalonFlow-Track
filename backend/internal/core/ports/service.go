package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/salonflow/salonflow-track/internal/core/domain"
)

// ServiceRepository defines persistence operations for Service.
type ServiceRepository interface {
	Create(ctx context.Context, service *domain.Service) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Service, error)
	GetByName(ctx context.Context, name string) (*domain.Service, error)
	List(ctx context.Context, filter ServiceFilter) ([]domain.Service, int, error)
	Update(ctx context.Context, service *domain.Service) error
	Delete(ctx context.Context, id uuid.UUID) error
	CountByStatus(ctx context.Context) (total, active, inactive int, err error)
}

// ServiceFilter holds query parameters for listing services.
type ServiceFilter struct {
	Status   string
	Category string
	Search   string
	Limit    int
	Offset   int
}
